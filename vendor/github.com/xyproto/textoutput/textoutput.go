//go:build !windows
// +build !windows

// Package textoutput offers a simple way to use vt100 and output colored text
package textoutput

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/vt100"
)

// CharAttribute is a rune and a color attribute
type CharAttribute struct {
	A vt100.AttributeColor
	R rune
}

// TextOutput keeps state about verbosity and if colors are enabled
type TextOutput struct {
	lightReplacer *strings.Replacer
	darkReplacer  *strings.Replacer
	color         bool
	enabled       bool
}

// Respect the NO_COLOR environment variable
var EnvNoColor = env.Bool("NO_COLOR")

// New creates a new TextOutput struct, which is
// enabled by default and with colors turned on.
// If the NO_COLOR environment variable is set, colors are disabled.
func New() *TextOutput {
	o := &TextOutput{nil, nil, !EnvNoColor, true}
	o.initializeTagReplacers()
	return o
}

// NewTextOutput can initialize a new TextOutput struct,
// which can have colors turned on or off and where the
// output can be enabled (verbose) or disabled (silent).
// If NO_COLOR is set, colors are disabled, regardless.
func NewTextOutput(color, enabled bool) *TextOutput {
	if EnvNoColor {
		color = false
	}
	o := &TextOutput{nil, nil, color, enabled}
	o.initializeTagReplacers()
	return o
}

// OutputTags will output text that may have tags like "<blue>", "</blue>" or "<off>" for
// enabling or disabling color attributes. Respects the color/enabled settings
// of this TextOutput.
func (o *TextOutput) OutputTags(colors ...string) {
	if o.enabled {
		fmt.Println(o.Tags(colors...))
	}
}

// Given a line with words and several color strings, color the words
// in the order of the colors. The last color will color the rest of the
// words.
func (o *TextOutput) OutputWords(line string, colors ...string) {
	if o.enabled {
		fmt.Println(o.Words(line, colors...))
	}
}

// Write a message to stdout if output is enabled
func (o *TextOutput) Println(msg ...interface{}) {
	if o.enabled {
		fmt.Println(o.InterfaceTags(msg...))
	}
}

// Write a message to the given io.Writer if output is enabled
func (o *TextOutput) Fprintln(w io.Writer, msg ...interface{}) {
	if o.enabled {
		fmt.Fprintln(w, o.InterfaceTags(msg...))
	}
}

// Write a message to stdout if output is enabled
func (o *TextOutput) Printf(msg ...interface{}) {
	if !o.enabled {
		return
	}
	count := len(msg)
	if count == 0 {
		return
	} else if count == 1 {
		if fmtString, ok := msg[0].(string); ok {
			fmt.Print(fmtString)
		}
	} else { // > 1
		if fmtString, ok := msg[0].(string); ok {
			fmt.Printf(o.InterfaceTags(fmtString), msg[1:]...)
		} else {
			// fail
			fmt.Printf("%v", msg...)
		}
	}
}

// Write a message to the given io.Writer if output is enabled
func (o *TextOutput) Fprintf(w io.Writer, msg ...interface{}) {
	if !o.enabled {
		return
	}
	count := len(msg)
	if count == 0 {
		return
	} else if count == 1 {
		if fmtString, ok := msg[0].(string); ok {
			fmt.Fprint(w, fmtString)
		}
	} else { // > 1
		if fmtString, ok := msg[0].(string); ok {
			fmt.Fprintf(w, o.InterfaceTags(fmtString), msg[1:]...)
		} else {
			// fail
			fmt.Fprintf(w, "%v", msg...)
		}
	}
}

// Write a message to stdout if output is enabled
func (o *TextOutput) Print(msg ...interface{}) {
	if o.enabled {
		fmt.Print(o.InterfaceTags(msg...))
	}
}

// Write a message to the given io.Writer if output is enabled
func (o *TextOutput) Fprint(w io.Writer, msg ...interface{}) {
	if o.enabled {
		fmt.Fprint(w, o.InterfaceTags(msg...))
	}
}

// Write an error message in red to stderr if output is enabled
func (o *TextOutput) Err(msg string) {
	if o.enabled {
		if o.color {
			vt100.Red.Error(msg)
		} else {
			vt100.Default.Error(msg)
		}
	}
}

// Write an error message to stderr and quit with exit code 1
func (o *TextOutput) ErrExit(msg string) {
	o.Err(msg)
	os.Exit(1)
}

// Deprectated
func (o *TextOutput) IsEnabled() bool {
	return o.enabled
}

// Enabled returns true if any output is enabled
func (o *TextOutput) Enabled() bool {
	return o.enabled
}

// Disabled returns true if all output is disabled
func (o *TextOutput) Disabled() bool {
	return !o.enabled
}

func (o *TextOutput) DarkRed(s string) string {
	if o.color {
		return vt100.Red.Get(s)
	}
	return s
}

func (o *TextOutput) DarkGreen(s string) string {
	if o.color {
		return vt100.Green.Get(s)
	}
	return s
}

func (o *TextOutput) DarkYellow(s string) string {
	if o.color {
		return vt100.Yellow.Get(s)
	}
	return s
}

func (o *TextOutput) DarkBlue(s string) string {
	if o.color {
		return vt100.Blue.Get(s)
	}
	return s
}

func (o *TextOutput) DarkPurple(s string) string {
	if o.color {
		return vt100.Magenta.Get(s)
	}
	return s
}

func (o *TextOutput) DarkCyan(s string) string {
	if o.color {
		return vt100.Cyan.Get(s)
	}
	return s
}

func (o *TextOutput) DarkGray(s string) string {
	if o.color {
		return vt100.DarkGray.Get(s)
	}
	return s
}

func (o *TextOutput) LightRed(s string) string {
	if o.color {
		return vt100.LightRed.Get(s)
	}
	return s
}

func (o *TextOutput) LightGreen(s string) string {
	if o.color {
		return vt100.LightGreen.Get(s)
	}
	return s
}

func (o *TextOutput) LightYellow(s string) string {
	if o.color {
		return vt100.LightYellow.Get(s)
	}
	return s
}

func (o *TextOutput) LightBlue(s string) string {
	if o.color {
		return vt100.LightBlue.Get(s)
	}
	return s
}

func (o *TextOutput) LightPurple(s string) string {
	if o.color {
		return vt100.LightMagenta.Get(s)
	}
	return s
}

func (o *TextOutput) LightCyan(s string) string {
	if o.color {
		return vt100.LightCyan.Get(s)
	}
	return s
}

func (o *TextOutput) White(s string) string {
	if o.color {
		return vt100.White.Get(s)
	}
	return s
}

// Given a line with words and several color strings, color the words
// in the order of the colors. The last color will color the rest of the
// words.
func (o *TextOutput) Words(line string, colors ...string) string {
	if o.color {
		return vt100.Words(line, colors...)
	}
	return line
}

// Change the color state in the terminal emulator
func (o *TextOutput) ColorOn(attribute1, attribute2 int) string {
	if !o.color {
		return ""
	}
	return fmt.Sprintf("\033[%d;%dm", attribute1, attribute2)
}

// Change the color state in the terminal emulator
func (o *TextOutput) ColorOff() string {
	if !o.color {
		return ""
	}
	return "\033[0m"
}

// Replace <blue> with starting a light blue color attribute and <off> with using the default attributes.
// </blue> can also be used for using the default attributes.
func (o *TextOutput) LightTags(colors ...string) string {
	return o.lightReplacer.Replace(strings.Join(colors, ""))
}

// Same as LightTags
func (o *TextOutput) Tags(colors ...string) string {
	return o.LightTags(colors...)
}

// InterfaceTags is the same as LightTags, but with interfaces
func (o *TextOutput) InterfaceTags(colors ...interface{}) string {
	var sb strings.Builder
	for _, color := range colors {
		if colorString, ok := color.(string); ok {
			sb.WriteString(colorString)
		} else {
			sb.WriteString(fmt.Sprintf("%s", color))
		}
	}
	return o.LightTags(sb.String())
}

// Replace <blue> with starting a light blue color attribute and <off> with using the default attributes.
// </blue> can also be used for using the default attributes.
func (o *TextOutput) DarkTags(colors ...string) string {
	return o.darkReplacer.Replace(strings.Join(colors, ""))
}

func (o *TextOutput) DisableColors() {
	o.color = false
	o.initializeTagReplacers()
}

func (o *TextOutput) EnableColors() {
	o.color = true
	o.initializeTagReplacers()
}

func (o *TextOutput) Disable() {
	o.enabled = false
}

func (o *TextOutput) Enable() {
	o.enabled = true
}

func (o *TextOutput) initializeTagReplacers() {
	// Initialize tag replacement tables, with as few memory allocations as possible (no append)
	off := vt100.NoColor()
	rs := make([]string, len(vt100.LightColorMap)*4+2)
	i := 0
	if o.color {
		for key, value := range vt100.LightColorMap {
			rs[i] = "<" + key + ">"
			i++
			rs[i] = value.String()
			i++
			rs[i] = "</" + key + ">"
			i++
			rs[i] = off
			i++
		}
		rs[i] = "<off>"
		i++
		rs[i] = off
	} else {
		for key := range vt100.LightColorMap {
			rs[i] = "<" + key + ">"
			i++
			rs[i] = ""
			i++
			rs[i] = "</" + key + ">"
			i++
			rs[i] = ""
			i++
		}
		rs[i] = "<off>"
		i++
		rs[i] = ""
	}
	o.lightReplacer = strings.NewReplacer(rs...)
	// Initialize the replacer for the dark color scheme, while reusing the rs slice
	i = 0
	if o.color {
		for key, value := range vt100.DarkColorMap {
			rs[i] = "<" + key + ">"
			i++
			rs[i] = value.String()
			i++
			rs[i] = "</" + key + ">"
			i++
			rs[i] = off
			i++
		}
		rs[i] = "<off>"
		i++
		rs[i] = off
	} else {
		for key := range vt100.DarkColorMap {
			rs[i] = "<" + key + ">"
			i++
			rs[i] = ""
			i++
			rs[i] = "</" + key + ">"
			i++
			rs[i] = ""
			i++
		}
		rs[i] = "<off>"
		i++
		rs[i] = ""
	}
	o.darkReplacer = strings.NewReplacer(rs...)
}

// ExtractToSlice iterates over an ANSI encoded string, parsing out color codes and places it in
// a slice of CharAttribute. Each CharAttribute in the slice represents a character in the
// input string and its corresponding color attributes. This function handles escaping sequences
// and converts ANSI color codes to vt100.AttributeColor structs.
// The returned uint is the number of stored elements.
func (o *TextOutput) ExtractToSlice(s string, pcc *[]CharAttribute) uint {
	var (
		escaped      bool
		colorcode    strings.Builder
		currentColor vt100.AttributeColor
	)
	counter := uint(0)
	for _, r := range s {
		switch {
		case escaped && r == 'm':
			colorAttributes := strings.Split(strings.TrimPrefix(colorcode.String(), "["), ";")
			if len(colorAttributes) != 1 || colorAttributes[0] != "0" {
				for _, attribute := range colorAttributes {
					if attributeNumber, err := strconv.Atoi(attribute); err == nil { // success
						currentColor = append(currentColor, byte(attributeNumber))
					} else {
						continue
					}
				}
				// Strip away leading 0 color attribute, if there are more than 1
				if len(currentColor) > 1 && currentColor[0] == 0 {
					currentColor = currentColor[1:]
				}
			} else {
				currentColor = vt100.NewAttributeColor()
			}
			colorcode.Reset()
			escaped = false
		case r == '\033':
			escaped = true
		case escaped && r != 'm':
			colorcode.WriteRune(r)
		default:
			(*pcc)[counter] = CharAttribute{currentColor, r}
			counter++
		}
	}
	return counter
}

// Extract iterates over an ANSI encoded string, parsing out color codes and creating a slice
// of CharAttribute structures. Each CharAttribute in the slice represents a character in the
// input string and its corresponding color attributes. This function handles escaping sequences
// and converts ANSI color codes to vt100.AttributeColor structs.
func (o *TextOutput) Extract(s string) []CharAttribute {
	cc := make([]CharAttribute, len(s))
	n := o.ExtractToSlice(s, &cc)
	return cc[:n]
}
