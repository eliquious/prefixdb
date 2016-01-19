package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	skip bool
	s    string
	stmt Node
	err  string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

// ParserTestSuite executes all the parser tests
type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) SetupTest() {
}

func (suite *ParserTestSuite) TearDownTest() {
}

func (suite *ParserTestSuite) validate(tests []TestCase) {
	for i, tt := range tests {
		if tt.skip {
			suite.T().Logf("skipping test of '%s'", tt.s)
			continue
		}
		stmt, err := ParseString(tt.s)

		if !reflect.DeepEqual(tt.err, errstring(err)) {
			suite.T().Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			suite.T().Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

// Ensure the parser will return an error for unknown statements
func (suite *ParserTestSuite) TestInvalidStatement() {
	var tests = []TestCase{

		// Errors
		{s: `a bad statement.`, err: `found IDENTIFIER (a), expected CREATE, DROP, SELECT, DELETE, UPSERT at line 1, char 1`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into CREATE KEYSPACE statements
func (suite *ParserTestSuite) TestCreateKeyspace() {
	var tests = []TestCase{
		{
			s:    `CREATE KEYSPACE acme WITH KEY id`,
			stmt: &CreateStatement{Keyspace: "acme", Keys: []string{"id"}},
		},
		{
			s:    `CREATE KEYSPACE acme WITH KEYS id, category`,
			stmt: &CreateStatement{Keyspace: "acme", Keys: []string{"id", "category"}},
		},
		{
			s:    `CREATE KEYSPACE acme WITH KEYS id`,
			stmt: &CreateStatement{Keyspace: "acme", Keys: []string{"id"}},
		},

		// Errors
		{s: `CREATE `, err: `found EOF, expected KEYSPACE at line 1, char 9`},
		{s: `CREATE KEYSPACE `, err: `found EOF, expected keyspace at line 1, char 18`},
		{s: `CREATE KEYSPACE acme.example.`, err: `found EOF, expected identifier at line 1, char 30`},
		{s: `CREATE KEYSPACE acme.example. `, err: `found WS, expected identifier at line 1, char 30`},
		{s: `CREATE KEYSPACE .example`, err: `found ., expected keyspace at line 1, char 17`},
		{s: `CREATE KEYSPACE acme`, err: `found EOF, expected WITH at line 1, char 22`},
		{s: `CREATE KEYSPACE acme WITH`, err: `found EOF, expected KEY, KEYS at line 1, char 27`},
		{s: `CREATE KEYSPACE acme WITH KEY`, err: `found EOF, expected identifier at line 1, char 31`},
		{s: `CREATE KEYSPACE acme WITH KEYS`, err: `found EOF, expected identifier at line 1, char 32`},
		{s: `CREATE KEYSPACE acme WITH KEYS id,`, err: `found EOF, expected identifier at line 1, char 35`},
		{s: `CREATE KEYSPACE acme WITH KEYS id, ""`, err: `found TEXTUAL, expected identifier at line 1, char 35`},
		{s: `CREATE KEYSPACE acme WITH KEY id,`, err: `found ,, expected EOF, SEMICOLON at line 1, char 33`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into DROP KEYSPACE statements
func (suite *ParserTestSuite) TestDropKeyspace() {
	var tests = []TestCase{
		{
			s:    `DROP KEYSPACE acme`,
			stmt: &DropStatement{Keyspace: "acme"},
		},

		// Errors
		{s: `DROP `, err: `found EOF, expected KEYSPACE at line 1, char 7`},
		{s: `DROP KEYSPACE `, err: `found EOF, expected keyspace at line 1, char 16`},
		{s: `DROP KEYSPACE acme.example.`, err: `found EOF, expected identifier at line 1, char 28`},
		{s: `DROP KEYSPACE acme.example. `, err: `found WS, expected identifier at line 1, char 28`},
		{s: `DROP KEYSPACE .example`, err: `found ., expected keyspace at line 1, char 15`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into SELECT statements
func (suite *ParserTestSuite) TestSelectStatement() {
	var tests = []TestCase{
		{
			s: `SELECT FROM users WHERE username = "bugs.bunny"`,
			stmt: &SelectStatement{Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value:        StringLiteral{"bugs.bunny"},
					},
				},
			},
		},
		{
			s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck"`,
			stmt: &SelectStatement{Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value: StringLiteralGroup{
							Values:   []string{"bugs.bunny", "daffy.duck"},
							Operator: OrOperator},
					},
				},
			},
		},
		{
			s: `SELECT FROM users WHERE username = "bugs.bunny" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01"`,
			stmt: &SelectStatement{Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value:        StringLiteral{"bugs.bunny"},
					},
					BetweenExpression{
						KeyAttribute: "timestamp",
						Values: StringLiteralGroup{
							Values:   []string{"2015-01-01", "2016-01-01"},
							Operator: AndOperator},
					},
				},
			},
		},
		{
			s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01"`,
			stmt: &SelectStatement{Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value: StringLiteralGroup{
							Values:   []string{"bugs.bunny", "daffy.duck"},
							Operator: OrOperator,
						},
					},
					BetweenExpression{
						KeyAttribute: "timestamp",
						Values: StringLiteralGroup{
							Values:   []string{"2015-01-01", "2016-01-01"},
							Operator: AndOperator,
						},
					},
				},
			},
		},
		{
			s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01" AND topic = "hunting"`,
			stmt: &SelectStatement{Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value: StringLiteralGroup{
							Values:   []string{"bugs.bunny", "daffy.duck"},
							Operator: OrOperator,
						},
					},
					BetweenExpression{
						KeyAttribute: "timestamp",
						Values: StringLiteralGroup{
							Values:   []string{"2015-01-01", "2016-01-01"},
							Operator: AndOperator,
						},
					},
					EqualityExpression{
						KeyAttribute: "topic",
						Value:        StringLiteral{"hunting"},
					},
				},
			},
		},

		// Errors
		{s: `SELECT`, err: `found EOF, expected FROM at line 1, char 8`},
		{s: `SELECT FROM `, err: `found EOF, expected keyspace at line 1, char 14`},
		{s: `SELECT FROM users`, err: `found EOF, expected WHERE at line 1, char 19`},
		{s: `SELECT FROM users WHERE`, err: `found EOF, expected identifier at line 1, char 25`},
		{s: `SELECT FROM users WHERE username`, err: `found EOF, expected EQ, BETWEEN at line 1, char 34`},
		{s: `SELECT FROM users WHERE username =`, err: `found EOF, expected string at line 1, char 35`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" OR`, err: `found EOF, expected string at line 1, char 52`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" AND`, err: `found EOF, expected identifier at line 1, char 53`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" OR`, err: `found OR, expected EOF, SEMICOLON, AND at line 1, char 65`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp`, err: `found EOF, expected EQ, BETWEEN at line 1, char 79`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN`, err: `found EOF, expected string at line 1, char 87`},
		{s: `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01"`, err: `found EOF, expected AND at line 1, char 99`},
	}

	suite.validate(tests)
}

// Ensure the parser can parse strings into UPSERT statements
func (suite *ParserTestSuite) TestUpsertStatement() {
	var tests = []TestCase{
		{
			s: `UPSERT "{...}" INTO users WHERE username = "bugs.bunny"`,
			stmt: &UpsertStatement{
				Value:    "{...}",
				Keyspace: "users",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value:        StringLiteral{"bugs.bunny"},
					},
				},
			},
		},
		{
			s: `UPSERT "{...}" INTO users.convo.timestamp WHERE username = "bugs.bunny" AND convo_id = "5" AND timestamp = "2015-01-01T00:00:00.001Z"`,
			stmt: &UpsertStatement{
				Value:    "{...}",
				Keyspace: "users.convo.timestamp",
				Where: []Expression{
					EqualityExpression{
						KeyAttribute: "username",
						Value:        StringLiteral{"bugs.bunny"},
					},
					EqualityExpression{
						KeyAttribute: "convo_id",
						Value:        StringLiteral{"5"},
					},
					EqualityExpression{
						KeyAttribute: "timestamp",
						Value:        StringLiteral{"2015-01-01T00:00:00.001Z"},
					},
				},
			},
		},

		// Errors
		{s: `UPSERT`, err: `found EOF, expected string at line 1, char 8`},
		{s: `UPSERT user`, err: `found IDENTIFIER (user), expected string at line 1, char 8`},
		{s: `UPSERT "..." `, err: `found EOF, expected INTO at line 1, char 15`},
		{s: `UPSERT "..." INTO`, err: `found EOF, expected keyspace at line 1, char 19`},
		{s: `UPSERT "..." INTO users`, err: `found EOF, expected WHERE at line 1, char 25`},
		{s: `UPSERT "..." INTO users WHERE`, err: `found EOF, expected identifier at line 1, char 31`},
		{s: `UPSERT "..." INTO users WHERE username`, err: `found EOF, expected EQ at line 1, char 40`},
		{s: `UPSERT "..." INTO users WHERE username =`, err: `found EOF, expected string at line 1, char 41`},
		{s: `UPSERT "..." INTO users WHERE username = "bugs.bunny" OR`, err: `OR not allowed at line 1, char 55`},
		{s: `UPSERT "..." INTO users WHERE username = "bugs.bunny" AND`, err: `found EOF, expected identifier at line 1, char 59`},
		{s: `UPSERT "..." INTO users WHERE username = "bugs.bunny" AND timestamp`, err: `found EOF, expected EQ at line 1, char 69`},
		{s: `UPSERT "..." INTO users WHERE username = "bugs.bunny" AND timestamp BETWEEN`, err: `BETWEEN not allowed at line 1, char 69`},
	}

	suite.validate(tests)
}

// errstring converts an error to its string representation.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func BenchmarkBadStatement(b *testing.B) {
	stmt := "a bad statement"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkNestedCreateStatement(b *testing.B) {
	stmt := "CREATE KEYSPACE acme.example.dynamite"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkCreateStatement(b *testing.B) {
	stmt := "CREATE KEYSPACE acme"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkDropStatement(b *testing.B) {
	stmt := "DROP KEYSPACE acme"
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkSelectStatement(b *testing.B) {
	stmt := `SELECT FROM users WHERE username = "bugs.bunny"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkComplexSelectStatement(b *testing.B) {
	stmt := `SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01" AND topic = "hunting"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkDeleteStatement(b *testing.B) {
	stmt := `DELETE FROM users WHERE username = "bugs.bunny"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkComplexDeleteStatement(b *testing.B) {
	stmt := `DELETE FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01" AND topic = "hunting"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkUpsertStatement(b *testing.B) {
	stmt := `UPSERT "{...}" INTO users WHERE username = "bugs.bunny"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

func BenchmarkComplexUpsertStatement(b *testing.B) {
	stmt := `UPSERT "{...}" INTO users.convo.timestamp WHERE username = "bugs.bunny" AND convo_id = "5" AND timestamp = "2015-01-01T00:00:00.001Z"`
	for i := 0; i < b.N; i++ {
		NewParser(strings.NewReader(stmt)).ParseStatement()
	}
}

// func BenchmarkShowNamespacesStatement(b *testing.B) {
// 	stmt := "SHOW KEYSPACES"
// 	for i := 0; i < b.N; i++ {
// 		NewParser(strings.NewReader(stmt)).ParseStatement()
// 	}
// }
