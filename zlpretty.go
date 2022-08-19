package zlpretty

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

var (
	indentSize = 2
	indentChar = " "
	indent     = strings.Repeat(indentChar, indentSize)
)

var (
	colorDefault = color.New(color.Reset)

	colorTime   = color.New(color.FgHiBlack)
	colorField  = color.New(color.FgCyan)
	colorString = color.New(color.FgGreen)
	colorNumber = color.New(color.FgMagenta)

	colorDebug = color.New(color.FgMagenta, color.Bold)
	colorInfo  = color.New(color.FgGreen, color.Bold)
	colorWarn  = color.New(color.FgYellow, color.Bold)
	colorTrace = color.New(color.FgBlue, color.Bold)
	colorError = color.New(color.FgRed, color.Bold)
)

var jsonScheme = &json.ColorScheme{
	Int:       createFormat(color.FgMagenta),
	Uint:      createFormat(color.FgMagenta),
	Float:     createFormat(color.FgMagenta),
	Bool:      createFormat(color.FgYellow),
	String:    createFormat(color.FgGreen),
	Binary:    createFormat(color.FgRed),
	ObjectKey: createFormat(color.FgCyan),
	Null:      createFormat(color.FgBlue),
}

type ConsoleWriter struct {
	Out     io.Writer
	NoColor bool
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	err = json.UnmarshalWithOption(p, &event, json.DecodeFieldPriorityFirstWin())
	if err != nil {
		return
	}

	level := "????"
	lvlColor := colorDefault
	if l, ok := event[zerolog.LevelFieldName].(string); ok {
		level = strings.ToUpper(l)
		lvlColor = levelColor(l)
		delete(event, zerolog.LevelFieldName)
	}
	w.Print(colorTime, formatTime(event[zerolog.TimestampFieldName]))
	delete(event, zerolog.TimestampFieldName)
	w.Print(lvlColor, " ["+level+"] ")
	if message, ok := event[zerolog.MessageFieldName].(string); ok {
		w.Out.Write([]byte(message))
		delete(event, zerolog.MessageFieldName)
	}

	for field := range event {
		w.Print(colorField, "\n"+indent+field)
		w.Out.Write([]byte(": "))
		switch value := event[field].(type) {
		case string:
			w.Print(colorString, strconv.Quote(value))
		case json.Number:
			w.Print(colorNumber, value.String())
		default:
			v, err := json.MarshalIndentWithOption(value, indent, indent, w.jsonOptions)
			if err != nil {
				return 0, err
			}

			w.Out.Write(v)
		}
	}
	w.Out.Write([]byte("\n"))
	return
}

func (w ConsoleWriter) jsonOptions(opts *json.EncodeOption) {
	json.DisableHTMLEscape()(opts)
	json.DisableNormalizeUTF8()(opts)

	if w.NoColor {
		return
	}

	json.Colorize(jsonScheme)(opts)
}

func formatTime(t interface{}) string {
	switch t := t.(type) {
	case string:
		return t
	case json.Number:
		u, _ := t.Int64()
		return time.Unix(u, 0).Format(time.RFC3339)
	}
	return "<nil>"
}

func (cw ConsoleWriter) Print(c *color.Color, s string) {
	if cw.NoColor {
		cw.Out.Write([]byte(s))
		return
	}
	c.Fprint(cw.Out, s)
}

func levelColor(level string) *color.Color {
	switch level {
	case "debug":
		return colorDebug
	case "info":
		return colorInfo
	case "warn":
		return colorWarn
	case "trace":
		return colorTrace
	case "error", "fatal", "panic":
		return colorError
	default:
		return colorDefault
	}
}

func createFormat(c color.Attribute) json.ColorFormat {
	return json.ColorFormat{
		Header: fmt.Sprintf("%s[%dm", "\x1b", c),
		Footer: fmt.Sprintf("%s[%dm", "\x1b", 0),
	}
}
