package parser

import (
	"bytes"
	"strings"
)

type NodeType int

const (
	CreateKeyspaceType NodeType = iota
	DropKeyspaceType
	SelectType
	UpsertType
	DeleteType
	StringLiteralType
	StringLiteralGroupType
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

func (CreateStatement) NodeType() NodeType {
	return CreateKeyspaceType
}

// String returns a string representation
func (c CreateStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("CREATE KEYSPACE ")
	buf.WriteString(c.Keyspace)
	buf.WriteString(" WITH ")
	if len(c.Keys) > 1 {
		buf.WriteString("KEYS ")
	} else {
		buf.WriteString("Key ")
	}

	buf.WriteString(strings.Join(c.Keys, ", ") + ";")
	return buf.String()
}

type DropStatement struct {
	Keyspace string
}

func (DropStatement) NodeType() NodeType {
	return DropKeyspaceType
}

// String returns a string representation
func (d DropStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("DROP KEYSPACE ")
	buf.WriteString(d.Keyspace)
	buf.WriteString(";")
	return buf.String()
}

type SelectStatement struct {
	Keyspace string
	Where    []Expression
}

func (SelectStatement) NodeType() NodeType {
	return SelectType
}

// String returns a string representation
func (s SelectStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("SELECT FROM ")
	buf.WriteString(s.Keyspace)
	buf.WriteString(" WHERE ")

	var filters []string
	for _, exp := range s.Where {
		filters = append(filters, exp.String())
	}
	buf.WriteString(strings.Join(filters, " AND "))
	buf.WriteString(";")
	return buf.String()
}

type UpsertStatement struct {
	Value    string
	Keyspace string
	Where    []Expression
}

func (UpsertStatement) NodeType() NodeType {
	return UpsertType
}

// String returns a string representation
func (u UpsertStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("UPSERT ")
	buf.WriteString(u.Value)
	buf.WriteString(" INTO ")
	buf.WriteString(u.Keyspace)
	buf.WriteString(" WHERE ")

	var filters []string
	for _, exp := range u.Where {
		filters = append(filters, exp.String())
	}
	buf.WriteString(strings.Join(filters, " AND "))
	buf.WriteString(";")
	return buf.String()
}

type DeleteStatement struct {
	Keyspace string
	Where    []Expression
}

func (DeleteStatement) NodeType() NodeType {
	return DeleteType
}

// String returns a string representation
func (d DeleteStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("DELETE FROM ")
	buf.WriteString(d.Keyspace)
	buf.WriteString(" WHERE ")

	var filters []string
	for _, exp := range d.Where {
		filters = append(filters, exp.String())
	}
	buf.WriteString(strings.Join(filters, " AND "))
	buf.WriteString(";")
	return buf.String()
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

type StringLiteralGroup struct {
	Values   []string
	Operator Operator
}

func (s StringLiteralGroup) NodeType() NodeType {
	return StringLiteralGroupType
}

func (s StringLiteralGroup) String() string {
	return strings.Join(s.Values, s.Operator.String())
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
