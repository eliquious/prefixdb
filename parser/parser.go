package parser

import (
	"io"
	"strings"

	"github.com/eliquious/lexer"
	tokens "github.com/eliquious/prefixdb/lexer"
)

// Parser represents an PrefixDB parser.
type Parser struct {
	s *lexer.TokenBuffer
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: lexer.NewTokenBuffer(r)}
}

// ParseString parses a statement string and returns its AST representation.
func ParseString(s string) (Node, error) {
	return NewParser(strings.NewReader(s)).ParseStatement()
}

// ParseStatement parses a string and returns a Statement AST object.
func (p *Parser) ParseStatement() (Node, error) {

	// Inspect the first token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case tokens.CREATE:
		return p.parseCreateStatement()
	case tokens.DROP:
		return p.parseDropStatement()
	case tokens.SELECT:
		return p.parseSelectStatement()
	case tokens.DELETE:
		return p.parseDeleteStatement()
	case tokens.UPSERT:
		return p.parseUpsertStatement()
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"CREATE", "DROP", "SELECT", "DELETE", "UPSERT"}, pos)
	}
}

// parseCreateStatement parses a string and returns an AST object.
// This function assumes the "CREATE" token has already been consumed.
func (p *Parser) parseCreateStatement() (Node, error) {

	// Inspect the first token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case tokens.KEYSPACE:
		return p.parseCreateKeyspaceStatement()
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"KEYSPACE"}, pos)
	}
}

// parseCreateKeyspaceStatement parses a string and returns a CreateKeyspaceStatement.
// This function assumes the "CREATE" token has already been consumed.
func (p *Parser) parseCreateKeyspaceStatement() (*CreateStatement, error) {
	stmt := &CreateStatement{}

	// Parse the name of the keyspace to be used
	lit, err := p.parseKeyspace()
	if err != nil {
		return nil, err
	}
	stmt.Keyspace = lit

	// Inspect the WITH token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case tokens.WITH:

		// Inspect the KEY or KEYS token.
		tok, pos, lit := p.scanIgnoreWhitespace()
		switch tok {
		case tokens.KEY:
			k, err := p.parseIdent()
			if err != nil {
				return nil, err
			} else {
				stmt.Keys = append(stmt.Keys, k)
			}
		case tokens.KEYS:
			k, err := p.parseIdentList()
			if err != nil {
				return nil, err
			} else {
				stmt.Keys = k
			}
		default:
			return nil, NewParseError(tokstr(tok, lit), []string{"KEY", "KEYS"}, pos)
		}
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"WITH"}, pos)
	}

	// Verify end of query
	tok, pos, lit = p.scanIgnoreWhitespace()
	switch tok {
	case lexer.EOF:
	case lexer.SEMICOLON:
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"EOF", "SEMICOLON"}, pos)
	}

	return stmt, nil
}

// parseDropStatement parses a string and returns a Statement AST object.
// This function assumes the "DROP" token has already been consumed.
func (p *Parser) parseDropStatement() (Node, error) {
	// Inspect the first token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case tokens.KEYSPACE:
		return p.parseDropKeyspaceStatement()
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"KEYSPACE"}, pos)
	}
}

// parseDropKeyspaceStatement parses a string and returns a DropStatement.
// This function assumes the "DROP" token has already been consumed.
func (p *Parser) parseDropKeyspaceStatement() (*DropStatement, error) {
	stmt := &DropStatement{}

	// Parse the name of the keyspace to be used
	lit, err := p.parseKeyspace()
	if err != nil {
		return nil, err
	}
	stmt.Keyspace = lit

	return stmt, nil
}

// parseSelectStatement parses a string and returns an AST object.
// This function assumes the "SELECT" token has already been consumed.
func (p *Parser) parseSelectStatement() (Node, error) {

	ks, where, err := p.parseFromWhere()
	if err != nil {
		return nil, err
	}

	return &SelectStatement{
		Keyspace: ks,
		Where:    where,
	}, nil
}

// parseDeleteStatement parses a string and returns an AST object.
// This function assumes the "DELETE" token has already been consumed.
func (p *Parser) parseDeleteStatement() (Node, error) {
	ks, where, err := p.parseFromWhere()
	if err != nil {
		return nil, err
	}

	return &DeleteStatement{
		Keyspace: ks,
		Where:    where,
	}, nil
}

func (p *Parser) parseFromWhere() (string, []Expression, error) {
	var keyspace string
	var exprs []Expression

	// Inspect the FROM token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case tokens.FROM:

		// Parse the name of the keyspace to be used
		lit, err := p.parseKeyspace()
		if err != nil {
			return "", nil, err
		}
		keyspace = lit

		// Inspect the WHERE token.
		tok, pos, lit := p.scanIgnoreWhitespace()
		switch tok {
		case tokens.WHERE:

			// Parse the WHERE clause
			exp, err := p.parseWhereClause(true, true)
			if err != nil {
				return "", nil, err
			}
			exprs = exp

		default:
			return "", nil, NewParseError(tokstr(tok, lit), []string{"WHERE"}, pos)
		}

	default:
		return "", nil, NewParseError(tokstr(tok, lit), []string{"FROM"}, pos)
	}
	return keyspace, exprs, nil
}

// parseUpsertStatement parses a string and returns an AST object.
// This function assumes the "UPSERT" token has already been consumed.
func (p *Parser) parseUpsertStatement() (Node, error) {

	// Parse value
	value, err := p.parseString()
	if err != nil {
		return nil, err
	}

	// Read INTO token
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != tokens.INTO {
		return nil, NewParseError(tokstr(tok, lit), []string{"INTO"}, pos)
	}

	// Parse keyspace name
	ks, err := p.parseKeyspace()
	if err != nil {
		return nil, err
	}

	// Read WHERE token
	tok, pos, lit = p.scanIgnoreWhitespace()
	if tok != tokens.WHERE {
		return nil, NewParseError(tokstr(tok, lit), []string{"WHERE"}, pos)
	}

	// Parse WHERE clause
	where, err := p.parseWhereClause(false, false)
	if err != nil {
		return nil, err
	}

	return &UpsertStatement{
		Value:    value,
		Keyspace: ks,
		Where:    where,
	}, nil
}

// parseWhereClause parses a string and returns an AST object.
// This function assumes the "WHERE" token has already been consumed.
func (p *Parser) parseWhereClause(allowBetween bool, allowLogicalOR bool) ([]Expression, error) {
	var expr []Expression

	// Read expression
	exp, err := p.parseExpression(allowBetween, allowLogicalOR)
	if err != nil {
		return expr, err
	}
	expr = append(expr, exp)

OUTER:
	for {

		// Test if there is another expression
		tok, pos, lit := p.scanIgnoreWhitespace()
		switch tok {
		case lexer.EOF, lexer.SEMICOLON:
			break OUTER
		case lexer.AND:

			exp, err := p.parseExpression(allowBetween, allowLogicalOR)
			if err != nil {
				return expr, err
			}
			expr = append(expr, exp)

		default:
			return nil, NewParseError(tokstr(tok, lit), []string{"EOF", "SEMICOLON", "AND"}, pos)
		}
	}
	return expr, nil
}

// parseExpression parses a string and returns an AST object.
func (p *Parser) parseExpression(allowBetween, allowLogicalOR bool) (Expression, error) {

	// Inspect the FROM token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT {
		return nil, NewParseError(tokstr(tok, lit), []string{"identifier"}, pos)
	}
	ident := lit

	// Inspect the operator token.
	tok, pos, lit = p.scanIgnoreWhitespace()
	switch tok {
	case lexer.EQ:
		expr, err := p.parseEqualityExpression(ident, allowLogicalOR)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case tokens.BETWEEN:
		if !allowBetween {
			return nil, &ParseError{Message: "BETWEEN not allowed", Pos: pos}
		}

		expr, err := p.parseBetweenExpression(ident)
		if err != nil {
			return nil, err
		}
		return expr, nil
	default:
		if allowBetween {
			return nil, NewParseError(tokstr(tok, lit), []string{"EQ", "BETWEEN"}, pos)
		}
		return nil, NewParseError(tokstr(tok, lit), []string{"EQ"}, pos)
	}
}

// parseEqualityExpression parses a string and returns an AST object.
func (p *Parser) parseEqualityExpression(ident string, allowLogicalOR bool) (Expression, error) {
	// expr := &EqualityExpression{KeyAttribute: ident}
	// var expr Expression

	// Parse the string value
	value, err := p.parseString()
	if err != nil {
		return nil, err
	}

	// Inspect the AND / OR token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case lexer.AND, lexer.EOF:
		p.unscan()
		return EqualityExpression{KeyAttribute: ident, Value: StringLiteral{value}}, nil
	case lexer.OR:
		// Upserts cannot contain OR equality clauses.
		if !allowLogicalOR {
			return nil, &ParseError{Message: "OR not allowed", Pos: pos}
		}

		values := []string{value}

		// Parse the string value
		value, err := p.parseString()
		if err != nil {
			return nil, err
		}
		values = append(values, value)
		return EqualityExpression{KeyAttribute: ident, Value: StringLiteralGroup{Operator: OrOperator, Values: values}}, nil
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"AND", "OR"}, pos)
	}
}

// parseBetweenExpression parses a string and returns an AST object.
func (p *Parser) parseBetweenExpression(ident string) (Expression, error) {
	expr := BetweenExpression{KeyAttribute: ident}

	// Parse the string value
	value, err := p.parseString()
	if err != nil {
		return nil, err
	}

	// Inspect the AND / OR token.
	tok, pos, lit := p.scanIgnoreWhitespace()
	switch tok {
	case lexer.AND:
		values := []string{value}

		// Parse the string value
		value, err := p.parseString()
		if err != nil {
			return nil, err
		}
		values = append(values, value)
		expr.Values = StringLiteralGroup{Operator: AndOperator, Values: values}
	default:
		return nil, NewParseError(tokstr(tok, lit), []string{"AND"}, pos)
	}
	return expr, nil
}

// parseKeyspace returns a keyspace title or an error
func (p *Parser) parseKeyspace() (string, error) {
	var keyspace string
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT {
		return "", NewParseError(tokstr(tok, lit), []string{"keyspace"}, pos)
	}
	keyspace = lit

	// Scan entire keyspace
	// Keyspaces are a period delimited list of identifiers
	var endPeriod bool
	for {
		tok, pos, lit = p.scan()
		if tok == lexer.DOT {
			keyspace += "."
			endPeriod = true
		} else if tok == lexer.IDENT {
			keyspace += lit
			endPeriod = false
		} else {
			break
		}
	}

	// remove last token
	p.unscan()

	// Keyspaces can't end on a period
	if endPeriod {
		return "", NewParseError(tokstr(tok, lit), []string{"identifier"}, pos)
	}
	return keyspace, nil
}

// parserString parses a string.
func (p *Parser) parseString() (string, error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != lexer.STRING {
		return "", NewParseError(tokstr(tok, lit), []string{"string"}, pos)
	}
	return lit, nil
}

// parseIdent parses an identifier.
func (p *Parser) parseIdent() (string, error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT {
		p.unscan()
		return "", NewParseError(tokstr(tok, lit), []string{"identifier"}, pos)
	}
	return lit, nil
}

// parseIdentList returns a list of attributes or an error
func (p *Parser) parseIdentList() ([]string, error) {
	var keys []string
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT {
		return keys, NewParseError(tokstr(tok, lit), []string{"identifier"}, pos)
	}
	keys = append(keys, lit)

	// Scan entire list
	// Key lists are comma delimited
	for {
		tok, pos, lit = p.scanIgnoreWhitespace()
		if tok == lexer.EOF {
			break
		} else if tok != lexer.COMMA {
			return keys, NewParseError(tokstr(tok, lit), []string{"COMMA", "EOF"}, pos)
		}

		k, err := p.parseIdent()
		if err != nil {
			return keys, err
		}
		keys = append(keys, k)
	}

	// remove last token
	p.unscan()

	return keys, nil
}

// scan returns the next token from the underlying scanner.
func (p *Parser) scan() (tok lexer.Token, pos lexer.Pos, lit string) { return p.s.Scan() }

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.s.Unscan() }

// peekRune returns the next rune that would be read by the scanner.
func (p *Parser) peekRune() rune { return p.s.Peek() }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok lexer.Token, pos lexer.Pos, lit string) {
	tok, pos, lit = p.scan()
	if tok == lexer.WS {
		tok, pos, lit = p.scan()
	}
	return
}

// tokstr returns a literal if provided, otherwise returns the token string.
func tokstr(tok lexer.Token, lit string) string {
	if tok == lexer.IDENT {
		return "IDENTIFIER (" + lit + ")"
	}
	return tok.String()
}
