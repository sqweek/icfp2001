// icfp2001 project icfp2001.go
package main

import (
	"fmt"
//	"strings"
	"github.com/sqweek/icfp2001/doc"
)



func ChooseStr(cond bool, t string, f string) string {
	if cond {
		return t
	}
	return f
}

func parse(s string) doc.Document {
	//var document []*doc.DecoratedText
	var document doc.Document
	
	var context doc.Stack
	var current doc.Decoration
	
	text := ""
	
	for i:=0; i<len(s); i++ {
		
		if s[i]!='<' {
			text += string(s[i])
		} else {
			fmt.Println("text: '"+text+"'")
			if len(text)>0 {
				document.Parts = append(document.Parts, doc.NewDecoratedText(current, text))
			}
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
				current = current.Apply(token)
			}
			fmt.Println("token: '"+ChooseStr(isEndToken==true,"/","")+token+"'")
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
