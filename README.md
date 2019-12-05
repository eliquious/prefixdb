# prefixdb
PrefixDB was going to be a toy database that specialized in range scans of key-value pairs. The database was never completed. It does have a functional parser and lexer for the query language though.


## Example queries

There is a command-line tool that will read from stdin and write out the tokens under `tools/lexer`.

```
CREATE KEYSPACE acme.example.dynamite
DROP KEYSPACE acme
SELECT FROM users WHERE username = "bugs.bunny"
SELECT FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01" AND topic = "hunting"
DELETE FROM users WHERE username = "bugs.bunny"
DELETE FROM users WHERE username = "bugs.bunny" OR "daffy.duck" AND timestamp BETWEEN "2015-01-01" AND "2016-01-01" AND topic = "hunting"
UPSERT "{...}" INTO users WHERE username = "bugs.bunny"`
UPSERT "{...}" INTO users.convo.timestamp WHERE username = "bugs.bunny" AND convo_id = "5" AND timestamp = "2015-01-01T00:00:00.001Z"
```

## Parser Benchmark

The parser is quite fast. Most queries are in the several microsecond range for parsing.

```
pkg: github.com/eliquious/prefixdb/parser
BenchmarkBadStatement-12              	 1000000	      1115 ns/op
BenchmarkNestedCreateStatement-12     	  500000	      3683 ns/op
BenchmarkCreateStatement-12           	  500000	      2518 ns/op
BenchmarkDropStatement-12             	 1000000	      2373 ns/op
BenchmarkSelectStatement-12           	  300000	      4247 ns/op
BenchmarkComplexSelectStatement-12    	  200000	     10344 ns/op
BenchmarkDeleteStatement-12           	  300000	      4236 ns/op
BenchmarkComplexDeleteStatement-12    	  200000	     10335 ns/op
BenchmarkUpsertStatement-12           	  300000	      4737 ns/op
BenchmarkComplexUpsertStatement-12    	  200000	      9772 ns/op
BenchmarkLargeUpsertStatement-12      	  100000	     19269 ns/op
PASS
ok  	github.com/eliquious/prefixdb/parser	20.503s
```
