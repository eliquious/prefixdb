package lexer

import (
	"github.com/eliquious/lexer"
)

func init() {
	// Loads keyword tokens into lexer
	lexer.LoadTokenMap(tokenMap)
}

const (
	// Starts the keywords with an offset from the built in tokens
	startKeywords lexer.Token = iota + 1000

	// CREATE starts a CREATE KEYSPACE query.
	CREATE

	// SELECT starts a SELECT FROM query.
	SELECT

	// UPSERT inserts or replaces a key-value pair.
	UPSERT

	// INTO sets which keyspace the new key-value pair will be inserted into.
	INTO

	// DELETE deletes keys from a keyspace.
	DELETE

	// DROP deletes an entire keyspace.
	DROP

	// FROM specifies which keyspace to select and delete from.
	FROM

	// WHERE allows for key filtering when selecting or deleting.
	WHERE

	// KEYSPACE signifies what is being created.
	KEYSPACE

	// WITH is used to prefix query options.
	WITH

	// KEY signifies a key attribute follows.
	KEY

	// KEYS signifies several key attributes follow.
	KEYS
	endKeywords

	// Separates the keywords from the conditionals
	startConditionals

	// BETWEEN filters an attribute by two values.
	BETWEEN
	endConditionals
)

var tokenMap = map[lexer.Token]string{

	CREATE:   "CREATE",
	SELECT:   "SELECT",
	DELETE:   "DELETE",
	UPSERT:   "UPSERT",
	INTO:     "INTO",
	DROP:     "DROP",
	FROM:     "FROM",
	WHERE:    "WHERE",
	KEYSPACE: "KEYSPACE",
	WITH:     "WITH",
	KEY:      "KEY",
	KEYS:     "KEYS",
	BETWEEN:  "BETWEEN",
}

// IsKeyword returns true if the token is a keyword.
func IsKeyword(tok lexer.Token) bool {
	return tok > startKeywords && tok < endKeywords
}

// IsConditional returns true if the token is a conditional clause.
func IsConditional(tok lexer.Token) bool {
	return tok > startConditionals && tok < endConditionals
}
