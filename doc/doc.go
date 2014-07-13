package doc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TextColour rune

const (
	Def TextColour = '?' /* default colour */
	W TextColour = 'w'
	R TextColour = 'r'
	G TextColour = 'g'
	B TextColour = 'b'
	C TextColour = 'c'
	M TextColour = 'm'
	Y TextColour = 'y'
	K TextColour = 'k'
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
	return Decoration{false,false,false,false,false,0,0,Def}
}

type TagStack []string
func (s TagStack) Empty() bool { return len(s) == 0 }
func (s TagStack) Peek() string   { return s[len(s)-1] }
func (s *TagStack) Push(i string)  { (*s) = append((*s), i) }
func (s *TagStack) Pop() string {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}
func (s *TagStack) HasColour() bool {
	for _, tag := range *s {
		switch tag {
		case "w", "r", "g", "b", "c", "y", "m", "k":
			return true
		}
	}
	return false
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
func (d Decoration) TagsTo(d2 Decoration,  tagStack TagStack) ([]string, TagStack) {
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
		if d2.Color == Def {
			// keep popping until we get back to default colour
			for tagStack.HasColour() {
				tag := tagStack.Pop()
				tags = append(tags, "/" + tag)
			}
		} else {
			tagStack.Push(string(d2.Color))
			tags = append(tags, string(d2.Color))
		}
	}
	return tags, tagStack
}
func (d Decoration) NTagsTo(d2 Decoration, tagStack TagStack) int {
	tags,_ := d.TagsTo(d2,tagStack)
	return len(tags)
}


func (document Document) GenerateSML() string {
	
	result := ""
	var tagStack TagStack
	
	current := DefaultDecoration()
	
	for _,part := range document.Parts {
		
		next := part.Decoration
		if next.Equals(current) == false {
			
			//min,_ := current.NTagsTo(next,tagStack)
			//should step back up the stack to check if there's a quicker way
			var tags TagStack
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
		if !current.HasContent() {
			continue
		}
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

	if last.HasContent() {
		document2.Parts = append(document2.Parts, last)
	}
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
		case "w", "r", "g", "b", "c", "y", "m", "k":
			out.Color = TextColour(tag[0])
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

func (dt *DecoratedText) HasContent() bool {
	return !(len(dt.Tokens) == 0 || (dt.Tt && len(dt.Tokens[0]) == 0))
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
	return fmt.Sprintf("%c%c%c%c%c%s%d%c",
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
