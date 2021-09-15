package zlpretty

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

var (
	f         = initFormatter()
	size      int
	separator []byte
	indent    = 4
)

func initFormatter() *colorjson.Formatter {
	f := colorjson.NewFormatter()
	f.Indent = indent
	f.NumberColor = color.New(color.FgHiMagenta)
	f.KeyColor = color.New(color.FgCyan)
	// f.Indent = "    "
	// f.Prefix = "    "
	// // json colors
	// f.SpaceColor = color.New(color.FgRed, color.Bold)
	// f.CommaColor = color.New(color.FgWhite, color.Bold)
	// f.ColonColor = color.New(color.FgYellow, color.Bold)
	// f.ObjectColor = color.New(color.FgGreen, color.Bold)
	// f.ArrayColor = color.New(color.FgHiRed)
	// f.FieldColor = color.New(color.FgCyan)
	// f.StringColor = color.New(color.FgHiYellow)
	// f.TrueColor = color.New(color.FgCyan, color.Bold)
	// f.FalseColor = color.New(color.FgHiRed)
	// f.NumberColor = color.New(color.FgHiMagenta)
	// f.NullColor = color.New(color.FgWhite, color.Bold)
	// f.StringQuoteColor = color.New(color.FgBlue, color.Bold)
	return f
}

type ConsoleWriter struct {
	Out     io.Writer
	NoColor bool
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	d := jsoniter.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&event)
	if err != nil {
		return
	}
	level := "????"
	lvlColor := color.Reset
	if l, ok := event[zerolog.LevelFieldName].(string); ok {
		level = strings.ToUpper(l)[0:4]
		lvlColor = levelColor(l)
	}
	color.New(color.FgHiBlack).Fprint(w.Out, formatTime(event[zerolog.TimestampFieldName]))
	color.New(lvlColor).Fprintf(w.Out, " [%s] ", level)
	if message, ok := event[zerolog.MessageFieldName].(string); ok {
		w.Out.Write([]byte(message))
	}
	for field := range event {
		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName:
			continue
		}

		color.New(color.FgCyan).Fprint(w.Out, "\n", strings.Repeat(" ", indent), field, "=")
		switch value := event[field].(type) {
		case string:
			w.Out.Write([]byte(strconv.Quote(value)))
		case json.Number:
			w.Out.Write([]byte(value.String()))
		default:
			v, err := f.Marshal(value)
			if err != nil {
				return 0, err
			}

			v = bytes.ReplaceAll(v, []byte("\n"), []byte("\n"+strings.Repeat(" ", indent)))

			w.Out.Write(v)
		}
	}
	w.Out.Write([]byte("\n"))
	return
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

func levelColor(level string) color.Attribute {
	switch level {
	case "debug":
		return color.FgMagenta
	case "info":
		return color.FgGreen
	case "warn":
		return color.FgYellow
	case "error", "fatal", "panic":
		return color.FgRed
	default:
		return color.Reset
	}
}
