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
