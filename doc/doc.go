package doc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TextColour int

/* a const block defines multiple constants at once. 'iota' is a special
 * value that starts at zero and increments by one each time it is
 * evaluated. Note that if the type or value is left out of a constant
 * definition, the previous type/value is assumed. So even though
 * I've only explicitly referenced iota once, it's evaluated 8 times
 * and we get constants W=0 R=1 G=2 etc. */
const (
	/* Note that variables/struct fields are zeroed by default;
	 * so I put W as the zero value rather. */
	W TextColour = iota
	R
	G
	B
	C
	M
	Y
	K
)

/* Go doesn't explicitly label symbols as public/private, but it does
 * have the concept. It follows a simple rule - symbols that start with
 * a capital letter can be accessed from outside the package, symbols
 * that don't are internal to the package. */
type Decoration struct {
	B     bool
	Em    bool
	I     bool
	S     bool
	Tt    bool
	U     int
	Size  int
	Color TextColour
}

func (d Decoration) copy() Decoration {
	return Decoration{d.B, d.Em, d.I, d.S, d.Tt, d.U, d.Size, d.Color}
}

type Document struct{
	Parts []*DecoratedText
}

//func (a Decoration) Compare(b Decoration) bool {
//	
//	
//}

func (a Decoration) Equals(b Decoration) bool {
	return a.B == b.B &&
		(a.S || a.Em == b.Em) &&
		a.I == b.I &&
		a.S == b.S &&
		a.Tt == b.Tt &&
		a.U == b.U &&
		a.Size == b.Size &&
		a.Color == b.Color
}

func DefaultDecoration() Decoration {
	return Decoration{false,false,false,false,false,0,0,W}
}

type StackStr []string
func (s StackStr) Empty() bool { return len(s) == 0 }
func (s StackStr) Peek() string   { return s[len(s)-1] }
func (s *StackStr) Push(i string)  { (*s) = append((*s), i) }
func (s *StackStr) Pop() string {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}


type Stack []Decoration
func (s Stack) Empty() bool { return len(s) == 0 }
func (s Stack) Peek() Decoration   { return s[len(s)-1] }
func (s *Stack) Push(i Decoration)  { (*s) = append((*s), i) }
func (s *Stack) Pop() Decoration {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}
func (d Decoration) TagsTo(d2 Decoration,  tagStack StackStr) ([]string, StackStr) {
	var tags []string
	//bold
	fmt.Println("d1:"+d.String()+", d2:"+d2.String())
	if d.B != d2.B{
		if d.B==true {
			tag := tagStack.Pop()
			for tag != "B" {
				fmt.Println("tag:"+tag)
				tags = append(tags, "/"+tag)
				if tagStack.Empty() {
					break
				}
				tag = tagStack.Pop()
			}
			tags = append(tags, "/B")
		} else {
			tagStack.Push("B")
			tags = append(tags, "B")
		}
	}
	
	//EM
	if d.Em!=d2.Em {
		tagStack.Push("Em")
		tags = append(tags,"Em")
	}
	
	
	//colours
	if d.Color != d2.Color {
		tagStack.Push(strings.ToLower(d2.Color.String()))
		tags = append(tags, strings.ToLower(d2.Color.String()))
	}
	return tags, tagStack
}
func (d Decoration) NTagsTo(d2 Decoration, tagStack StackStr) int {
	tags,_ := d.TagsTo(d2,tagStack)
	return len(tags)
}


func (document Document) GenerateSML() string {
	
	result := ""
	var tagStack StackStr
	
	current := DefaultDecoration()
	
	for _,part := range document.Parts {
		
		next := part.Decoration
		if next.Equals(current) == false {
			
			//min,_ := current.NTagsTo(next,tagStack)
			//should step back up the stack to check if there's a quicker way
			var tags StackStr
			fmt.Println("tagstack before:"+strings.Join(tagStack,","))
			tags,tagStack = current.TagsTo(next, tagStack)
			fmt.Println("tagstack after:"+strings.Join(tagStack,","))
			//need to search ahead to optimize the order of tags
			
			for _,tag := range tags {
				result += "<"+tag+">"
			}
			
			current = next
			
		}
		
		result += strings.Join(part.Tokens," ")
			
	}
	
	for tagStack.Empty()==false {
		result += "</"+tagStack.Pop()+">"
	}
	
	return result
}

func (document Document) Compact() Document {
	
	var document2 Document
	var last *DecoratedText
	
	for  _,current := range document.Parts {
		if last == nil {
			last = current
		} else if last.Decoration.Equals(current.Decoration) {
			if last.Tt {
				last.Tokens[0] += current.Tokens[0]
			} else {
				last.Tokens = append(last.Tokens, current.Tokens...)
			}
		} else {
			document2.Parts = append(document2.Parts, last)
			last = current
		}
	}
	
	document2.Parts = append(document2.Parts, last)
	return document2
}

/* Returns the new Decoration state achieved by applying the provided tags */
func (src Decoration) Apply(tags ...string) Decoration {
	out := src.copy()
	for _, tag := range tags {
		switch tag {
		case "B":
			out.B = true
		case "EM":
			out.Em = !src.Em
		case "I":
			out.I = true
		case "PL":
			out.U = 0
			out.B = false
			out.Em = false
			out.S = false
			out.Tt = false
		case "S":
			out.S = true
		case "TT":
			out.Tt = true
		case "U":
			if out.U < 3 {
				out.U++
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			out.Size, _ = strconv.Atoi(tag)
		case "w":
			out.Color = W
		case "r":
			out.Color = R
		case "g":
			out.Color = G
		case "b":
			out.Color = B
		case "c":
			out.Color = C
		case "y":
			out.Color = Y
		case "m":
			out.Color = M
		case "k":
			out.Color = K
		default:
			panic(tag)
		}
	}
	return out
}

type DecoratedText struct {
	Decoration // struct name by itself includes fields in this struct

	/* A list of tokens affected by this decoration. Each token will
	 * contain no whitespace, UNLESS the Tt flag is set in which
	 * case there will only be a single token */
	Tokens []string
}

// at this point a document can be represented by a sequence of DecoratedText structs

var reSpace *regexp.Regexp = regexp.MustCompile("[ \t\n\r\f]+")

/* converts runs of whitespace into a single space */
func collapseWhitespace(text string) string {
	return reSpace.ReplaceAllLiteralString(text, " ")
}

func NewDecoratedText(d Decoration, text string) *DecoratedText {
	var tokens []string
	if d.Tt {
		tokens = []string{text}
	} else {
		tokens = strings.Split(collapseWhitespace(text), " ")
	}
	return &DecoratedText{d, tokens}
}

func (colour TextColour) String() string {
	/* no breaks needed in switch statement, use 'fallthrough' keyword
	 * when you really need that behaviour. also switch works on pretty
	 * much any type :) */
	switch colour {
	case W:
		return "W"
	case R:
		return "R"
	case G:
		return "G"
	case B:
		return "B"
	case C:
		return "C"
	case M:
		return "M"
	case Y:
		return "Y"
	case K:
		return "K"
	}
	panic(colour)
}

func choose(cond bool, t rune, f rune) rune {
	if cond {
		return t
	}
	return f
}

func underlineStr(underline int) string {
	switch underline {
	case 0:
		return "|"
	case 1:
		return "-"
	case 2:
		return "="
	case 3:
		return "≡"
	}
	panic(underline)
}

func (d *Decoration) String() string {
	return fmt.Sprintf("%c%c%c%c%c%s%d%v",
		choose(d.B, 'B', 'b'),
		choose(d.Em, 'E', 'e'),
		choose(d.I, 'I', 'i'),
		choose(d.S, 'S', 's'),
		choose(d.Tt, 'T', 't'),
		underlineStr(d.U),
		d.Size,
		d.Color)
}

func (d *DecoratedText) String() string {
	return fmt.Sprintf("{%v} %v", &d.Decoration, d.Tokens)
}
