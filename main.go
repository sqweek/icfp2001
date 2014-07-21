// icfp2001 project icfp2001.go
package main

import (
	"fmt"
	"strings"
	"github.com/sqweek/icfp2001/doc"
)

func tokens2chan(s string, tokens chan string) {
	/* "defer" is kind of like "finally". When the expression drops out
	 * of scope (ie when tokens2chan returns), close(tokens) will run */
	defer close(tokens)

	i := 0 // token start index
	for i < len(s) {
		var j int // token length
		if s[i] == '<' {
			/* note j is the position relative to i, because a slice of the
			 * original string starting at i is passed in */
			j = strings.IndexRune(s[i:], '>')
			if j == -1 || j <= 1 {
				return //malformed doc, missing close tag or empty tag
			}
			j++ /* include closing > */
		} else {
			j = strings.IndexRune(s[i:], '<')
			if j == -1 {
				tokens <- s[i:]
				return
			}
		}
		tokens <- s[i:i + j]
		i += j
	}
}

func parse(s string) doc.Document {
	tokens := make(chan string)
	/* 'go' spawns a seperate "goroutine" to run a function */
	go tokens2chan(s, tokens)

	var document doc.Document
	
	var context doc.Stack
	var current doc.Decoration = doc.DefaultDecoration()
	
	/* 'range' over a channel reads from the chan until it is closed */
	for token := range tokens {
		fmt.Println("token:", token)
		if token[0] == '<' {
			tag := token[1:len(token) - 1]
			if tag[0] == '/' {
				if !context.Empty() {
					current = context.Pop()
				}
			} else {
				context.Push(current)
				current = current.Apply(tag)
			}
		} else {
			document.Parts = append(document.Parts, doc.NewDecoratedText(current, token))
		}
	}
	return document
}

//const s1 = "<r>  xxx </r><r>asdf</r><b> yyy  </b>"
const s1 = "<B>bold<r>red and bold</r>just bold</B>"
func main() {
	document := parse(s1).Compact()
	for _, block := range(document.Parts) {
		fmt.Println(block)
	}
	
	
	sml := document.GenerateSML()
	fmt.Println(sml)
}
