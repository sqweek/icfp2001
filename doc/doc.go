package doc

import (
	"strconv"
	"fmt"
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

/* Returns the new Decoration state achieved by applying the provided tags */
func (src Decoration) Apply(tags ...string) Decoration {
	out := src.copy()
	for _, tag := range(tags) {
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
