package terminal

// code form: https://github.com/daviddengcn/go-colortext

import (
	"io"
	"os"
	"strings"
)

type Color int

const (
	None = Color(iota)
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

var color_name = []string{
	"None",
	"Black", "Red", "Green", "Yellow",
	"Blue", "Magenta", "Cyan", "White",
	"BrightBlack", "BrightRed", "BrightGreen", "BrightYellow",
	"BrightBlue", "BrightMagenta", "BrightCyan", "BrightWhite"}

func (self Color) String() string {
	return color_name[self]
}

func Parse(str string) (Color, bool) {
	for i, name := range color_name {
		if strings.EqualFold(name, str) {
			return Color(i), true
		}
	}

	return None, false
}

var Writer io.Writer = os.Stdout

// reset color by default
func Reset() {
	resetColor()
}

// Change sets the foreground and background colors. If the value of the color is None,
// the corresponding color keeps unchanged.
func Change(fg Color, bg Color) {
	fgBright := false
	bgBright := false
	if fg >= BrightBlack {
		fg = fg - BrightBlack + 1
		fgBright = true
	}

	if bg >= BrightBlack {
		bg = bg - BrightBlack + 1
		bgBright = true
	}
	changeColor(fg, fgBright, bg, bgBright)
}

// Foreground changes the foreground color.
func Foreground(color Color) {
	Change(color, None)
}

// Background changes the background color.
func Background(color Color) {
	Change(None, color)
}
