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
		level = strings.ToUpper(l)
		lvlColor = levelColor(l)
	}
	color.New(color.FgHiBlack).Fprint(w.Out, formatTime(event[zerolog.TimestampFieldName]))
	color.New(lvlColor, color.Bold).Fprintf(w.Out, " [%s] ", level)
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
			color.New(color.FgGreen).Fprint(w.Out, strconv.Quote(value))
		case json.Number:
			color.New(color.FgMagenta).Fprint(w.Out, value.String())
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
	case "trace":
		return color.FgBlue
	case "error", "fatal", "panic":
		return color.FgRed
	default:
		return color.Reset
	}
}
