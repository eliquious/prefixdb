package parser

type Operator int

const (
	EqualityOperator Operator = iota
	AndOperator
	OrOperator
	BetweenOperator
)

type Expression interface {
	Operator() Operator
	Left() Node
	Right() Node
}

type EqualityExpression struct {
	KeyAttribute string
	Value        StringLiteral
}

func (e EqualityExpression) NodeType() NodeType {
	return ExpressionType
}

func (e EqualityExpression) Operator() Operator {
	return EqualityOperator
}

func (e EqualityExpression) Left() KeyAttribute {
	return KeyAttribute{e.KeyAttribute}
}

func (e EqualityExpression) Right() StringLiteral {
	return e.Value
}

type AndExpression struct {
	L Node
	R Node
}

func (a AndExpression) NodeType() NodeType {
	return ExpressionType
}

func (a AndExpression) Operator() Operator {
	return AndOperator
}

func (a AndExpression) Left() Node {
	return a.L
}

func (a AndExpression) Right() Node {
	return a.R
}

type OrExpression struct {
	L Node
	R Node
}

func (o OrExpression) NodeType() NodeType {
	return ExpressionType
}

func (o OrExpression) Operator() Operator {
	return OrOperator
}

func (o OrExpression) Left() Node {
	return o.L
}

func (o OrExpression) Right() Node {
	return o.R
}

type BetweenExpression struct {
	Attribute string
	Values    AndExpression
}

func (b BetweenExpression) NodeType() NodeType {
	return BetweenType
}

func (b BetweenExpression) Operator() Operator {
	return BetweenOperator
}

func (b BetweenExpression) Left() KeyAttribute {
	return KeyAttribute{b.Attribute}
}

func (b BetweenExpression) Right() Node {
	return b.Values
}
