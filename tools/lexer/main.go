package main

import (
        "flag"
        "fmt"
        "os"
        "strconv"

        "github.com/eliquious/lexer"
        _ "github.com/eliquious/prefixdb/lexer"
)

var skipWhitespace = flag.Bool("w", false, "skip whitespace tokens")

func main() {
        flag.Parse()
        l := lexer.NewScanner(os.Stdin)
        for {
                tok, pos, lit := l.Scan()

                // exit if EOF
                if tok == lexer.EOF {
                        break
                }

                // skip whitespace tokens
                if tok == lexer.WS && *skipWhitespace {
                        continue
                }

                // Print token
                if len(lit) > 0 {
                        fmt.Printf("[%4d:%-3d] %10s - %s\n", pos.Line, pos.Char, tok, strconv.QuoteToASCII(lit))
                } else {
                        fmt.Printf("[%4d:%-3d] %10s\n", pos.Line, pos.Char, tok)
                }
        }
}
