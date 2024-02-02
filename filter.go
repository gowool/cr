package cr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
)

const (
	tokenizePattern  = `=|<>|>|<|>=|<=|!=|AND|OR|IS\s+NULL|IS\s+NOT\s+NULL|LIKE|NOT\s+LIKE|ILIKE|NOT\s+ILIKE|IN|NOT\s+IN|\(|\)|"([^"\\]*(\\.[^"\\]*)*)"|\'([^\'\\]*(\\.[^\'\\]*)*)\'|\S+`
	normalizePatters = `(?i)\s+AND\s+|\s+OR\s+|\s+NULL|\s+IS\s+NULL|\s+IS\s+NOT\s+NULL|\s+LIKE\s+|\s+NOT\s+LIKE\s+|\s+ILIKE\s+|\s+NOT\s+ILIKE\s+|\s+IN\s+|\s+NOT\s+IN\s+`
)

var (
	tokenizeRegex  *regexp.Regexp
	normalizeRegex *regexp.Regexp

	normalizeFilter atomic.Bool
)

func init() {
	EnableNormalize()

	tokenizeRegex = regexp.MustCompile(tokenizePattern)
	normalizeRegex = regexp.MustCompile(normalizePatters)
}

func DisableNormalize() {
	normalizeFilter.Store(false)
}

func EnableNormalize() {
	normalizeFilter.Store(true)
}

type Filter struct {
	Operator   Operator
	Conditions []interface{}
}

func (f Filter) IsEmpty() bool {
	return len(f.Conditions) == 0
}

func (f Filter) ToSQL() (string, []interface{}) {
	var (
		conditions []string
		args       []interface{}
	)

	for _, cond := range f.Conditions {
		switch c := cond.(type) {
		case string:
			conditions = append(conditions, c)
		case Condition:
			if c.Operator.IsEmpty() {
				c.Operator = OpEqual
			}

			if c.Operator.Has(OpIS) && c.Value == nil {
				conditions = append(conditions, fmt.Sprintf("%s %s NULL", c.Column, c.Operator))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s (?)", c.Column, c.Operator))
				args = append(args, c.Value)
			}
		case Filter:
			if newStr, newArgs := c.ToSQL(); newStr != "" {
				conditions = append(conditions, fmt.Sprintf("(%s)", newStr))
				args = append(args, newArgs...)
			}
		}
	}

	operator := f.Operator
	if operator.IsEmpty() {
		operator = OpAND
	}

	return strings.Join(conditions, fmt.Sprintf(" %s ", operator)), args
}

func ParseFilter(filter string) (f Filter) {
	if filter = strings.TrimSpace(filter); filter != "" {
		if normalizeFilter.Load() {
			filter = normalize(filter)
		}
		f = toFilter(toTree(tokenize(filter)))
	}
	return
}

func normalize(filter string) string {
	return normalizeRegex.ReplaceAllStringFunc(filter, func(w string) string { return strings.ToUpper(w) })
}

func tokenize(filter string) []string {
	return tokenizeRegex.FindAllString(filter, -1)
}

func toTree(tokens []string) []string {
	operators := map[string]int{
		"OR":          1,
		"AND":         2,
		"=":           3,
		"<>":          3,
		">":           3,
		"<":           3,
		">=":          3,
		"<=":          3,
		"!=":          3,
		"IS NULL":     3,
		"IS NOT NULL": 3,
		"LIKE":        3,
		"NOT LIKE":    3,
		"ILIKE":       3,
		"NOT ILIKE":   3,
		"IN":          3,
		"NOT IN":      3,
	}

	var (
		stack  []string
		output []string
	)

	for _, token := range tokens {
		if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				output = append(output, op)
			}
			// Pop the "("
			stack = stack[:len(stack)-1]
		} else if precedence, isOperator := operators[token]; isOperator {
			for len(stack) > 0 {
				if stack[len(stack)-1] == "(" {
					break
				}
				topOp := stack[len(stack)-1]
				if topPrecedence, isTopOperator := operators[topOp]; isTopOperator && topPrecedence >= precedence {
					stack = stack[:len(stack)-1]
					output = append(output, topOp)
				} else {
					break
				}
			}
			stack = append(stack, token)
		} else {
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		op := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		output = append(output, op)
	}

	return output
}

func toFilter(tokens []string) Filter {
	var stack []interface{}

	for _, token := range tokens {
		t := strings.ToUpper(token)
		switch t {
		case "AND", "OR":
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, Filter{
				Operator:   Operator(t),
				Conditions: []interface{}{left, right},
			})
		case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE", "=", "<>", ">", "<", ">=", "<=", "!=":
			right := stack[len(stack)-1].(string)
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1].(string)
			stack = stack[:len(stack)-1]
			stack = append(stack, Condition{
				Column:   left,
				Operator: Operator(t),
				Value:    cast(right),
			})
		case "IN", "NOT IN":
			var (
				left  string
				right string
			)
			for len(stack) > 0 {
				r := stack[len(stack)-1]

				if s, ok := r.(string); ok {
					right = left + right
					left = s
					stack = stack[:len(stack)-1]
					continue
				}
				break
			}

			stack = append(stack, Condition{
				Column:   left,
				Operator: Operator(t),
				Value: apply(strings.Split(right, ","), func(item string) interface{} {
					return cast(item)
				}),
			})
		case "IS NULL", "IS NOT NULL":
			left := stack[len(stack)-1].(string)
			stack = stack[:len(stack)-1]
			stack = append(stack, Condition{
				Column:   left,
				Operator: Operator(strings.TrimSuffix(t, " NULL")),
			})
		default:
			stack = append(stack, token)
		}
	}

	if len(stack) == 1 {
		switch f := stack[0].(type) {
		case Filter:
			return f
		case Condition:
			return Filter{Conditions: []interface{}{f}}
		}
	}

	return Filter{}
}

func cast(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return value
	}

	last := len(value) - 1
	for _, c := range []byte{'\'', '"'} {
		if value[0] == c && value[last] == c {
			value = value[1:last]

			if v, err := uuid.Parse(value); err == nil {
				return v
			}

			if v, err := toTime(value); err == nil {
				return v
			}

			return strings.ReplaceAll(value, string([]byte{'\\', c}), string(c))
		}
	}

	if strings.ToLower(value) == "true" {
		return true
	}

	if strings.ToLower(value) == "false" {
		return false
	}

	if v, err := strconv.ParseInt(value, 0, 0); err == nil {
		return v
	}

	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v
	}

	return value
}
