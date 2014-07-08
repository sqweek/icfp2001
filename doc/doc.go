package doc

import (
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

type DecoratedText struct {
	Decoration // struct name by itself includes fields in this struct

	/* A list of tokens affected by this decoration. Each token will
	 * contain no whitespace, UNLESS the Tt flag is set in which
	 * case there will only be a single token */
	Tokens []string
}

// at this point a document can be represented by a sequence of DecoratedText structs

func (colour *TextColour) String() string {
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
	panic("unknown colour ", colour)
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
	panic("invalid underline state ", underline)
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
		d.Colour)
}
