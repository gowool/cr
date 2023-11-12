package cr

import (
	"fmt"
	"strings"
)

type Operator string

const (
	OpEmpty Operator = ""

	OpAdd Operator = "+"
	OpSub Operator = "-"
	OpMul Operator = "*"
	OpDiv Operator = "/"
	OpMod Operator = "%"

	OpBitAnd  Operator = "&"
	OpBitOr   Operator = "|"
	OpBitExOr Operator = "^"

	OpEqual    Operator = "="
	OpNotEqual Operator = "<>"
	OpGt       Operator = ">"
	OpGte      Operator = ">="
	OpLt       Operator = "<"
	OpLte      Operator = "<="

	OpAddEquals Operator = "+="
	OpSubEquals Operator = "-="
	OpMulEquals Operator = "*="
	OpDivEquals Operator = "/="
	OpModEquals Operator = "%="

	OpBitAndEquals Operator = "&="
	OpBitExEquals  Operator = "^-="
	OpBitOrEquals  Operator = "|*="

	OpIS      Operator = "IS"
	OpALL     Operator = "ALL"
	OpAND     Operator = "AND"
	OpANY     Operator = "ANY"
	OpBETWEEN Operator = "BETWEEN"
	OpEXISTS  Operator = "EXISTS"
	OpIN      Operator = "IN"
	OpLIKE    Operator = "LIKE"
	OpILIKE   Operator = "ILIKE"
	OpSIMILAR Operator = "SIMILAR TO"
	OpNOT     Operator = "NOT"
	OpOR      Operator = "OR"
	OpSOME    Operator = "SOME"
)

func (o Operator) String() string {
	return string(o)
}

func (o Operator) Append(a Operator) Operator {
	return Operator(fmt.Sprintf("%s %s", o, a))
}

func (o Operator) Prepend(a Operator) Operator {
	return Operator(fmt.Sprintf("%s %s", a, o))
}

func (o Operator) Has(a Operator) bool {
	return o == a || strings.Contains(o.String(), a.String())
}

func (o Operator) IsEmpty() bool {
	return o == OpEmpty
}
