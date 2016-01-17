package parser

type NodeType int

const (
	CreateKeyspaceType NodeType = iota
	DropKeyspaceType
	SelectType
	UpsertType
	DeleteType
	StringLiteralType
	ExpressionType
	KeyAttributeType
	BetweenType
)

type Node interface {
	NodeType() NodeType
	String() string
}

type CreateStatement struct {
	Keyspace string
	Keys     []string
}

type DropStatement struct {
	Keyspace string
}

type SelectStatement struct {
	Keyspace string
	Where    Expression
}

type UpsertStatement struct {
	Keyspace string
	Where    Expression
}

type DeleteStatement struct {
	Keyspace string
	Where    Expression
}

type StringLiteral struct {
	Value string
}

func (s StringLiteral) NodeType() NodeType {
	return StringLiteralType
}

func (s StringLiteral) String() string {
	return s.Value
}

type KeyAttribute struct {
	Attribute string
}

func (s KeyAttribute) NodeType() NodeType {
	return KeyAttributeType
}

func (s KeyAttribute) String() string {
	return s.Attribute
}
