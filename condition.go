package cr

type Condition struct {
	Column   string
	Operator Operator
	Value    interface{}
}
