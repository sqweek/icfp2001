// icfp2001 project icfp2001.go
package main

import (
	"fmt"
//	"strings"
	"github.com/sqweek/icfp2001/doc"
)

type stack []doc.Decoration
func (s stack) Empty() bool { return len(s) == 0 }
func (s stack) Peek() doc.Decoration   { return s[len(s)-1] }
func (s *stack) Push(i doc.Decoration)  { (*s) = append((*s), i) }
func (s *stack) Pop() doc.Decoration {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}

func DefaultDecoration() doc.Decoration {
	return doc.Decoration{false,false,false,false,false,0,0,doc.W}
}

func ChooseStr(cond bool, t string, f string) string {
	if cond {
		return t
	}
	return f
}

func describe(s string){
	
	texts := make([]string, 100,100)
	decs := make([]doc.Decoration, 100,100)
	
	var context stack
	var current doc.Decoration
	
	text := ""
	
	for i:=0; i<len(s); i++ {
		
		if s[i]!='<' {
			text += string(s[i])
		} else {
			fmt.Println("text: '"+text+"'")
			texts = append(texts,text)
			decs = append(decs,current)
			text =""
			
			token := ""
			isEndToken := false
			
			i = i+1
			if s[i]=='/' {
				isEndToken = true
				i=i+1
			}
			
			for ; i<len(s) && s[i]!='>'; i++ {
				token += string(s[i])
			}
			
			if isEndToken==true {
				if context.Empty()==false {
					current = context.Pop()
				}
			} else {
				context.Push(current)
				/* TODO easy optimisation here by checking if the new state
				 * is the same as the existing state in which case the tag is
				 * redundant and can be ignored. */
				current = current.Apply(token)
			}
			fmt.Println("token: '"+ChooseStr(isEndToken==true,"/","")+token+"'")
		}
	}
}

const s1 = "<r>  xxx </r><b> yyy  </b>"
func main() {
	describe(s1)
}
