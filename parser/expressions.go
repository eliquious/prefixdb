package parser

import "bytes"

type Operator int

const (
	EqualityOperator Operator = iota
	AndOperator
	OrOperator
	BetweenOperator
)

func (o Operator) String() string {
	switch o {
	case EqualityOperator:
		return " = "
	case AndOperator:
		return " AND "
	case OrOperator:
		return " OR "
	case BetweenOperator:
		return " BETWEEN "
	}
	return ""
}

type Expression interface {
	Operator() Operator
	String() string
}

type EqualityExpression struct {
	KeyAttribute string
	Value        Node
}

func (e EqualityExpression) NodeType() NodeType {
	return ExpressionType
}

func (e EqualityExpression) Operator() Operator {
	return EqualityOperator
}

func (e EqualityExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString(" ")
	buf.WriteString(e.KeyAttribute)
	buf.WriteString(e.Operator().String())
	buf.WriteString(e.Value.String())
	return buf.String()
}

type BetweenExpression struct {
	KeyAttribute string
	Values       StringLiteralGroup
}

func (b BetweenExpression) NodeType() NodeType {
	return BetweenType
}

func (b BetweenExpression) Operator() Operator {
	return BetweenOperator
}

func (b BetweenExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString(" ")
	buf.WriteString(b.KeyAttribute)
	buf.WriteString(BetweenOperator.String())
	buf.WriteString(b.Values.String())
	return buf.String()
}
