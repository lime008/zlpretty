package zlpretty

import (
	"io"
	"testing"

	"github.com/UnnoTed/horizontal"
	"github.com/rs/zerolog"
)

var testStruct = struct {
	Name string
	Age  int
}{
	Name: "John",
	Age:  32,
}

func testLog(logger zerolog.Logger) {
	logger.Error().Str("foo", "bar").Interface("testStruct", testStruct).Msg("hello")
}

func BenchmarkZlpretty(b *testing.B) {
	logger := zerolog.New(ConsoleWriter{Out: io.Discard})
	logger.Level(zerolog.DebugLevel)

	for i := 0; i < b.N; i++ {
		testLog(logger)
	}
}

func BenchmarkHorizontal(b *testing.B) {
	logger := zerolog.New(horizontal.ConsoleWriter{Out: io.Discard})
	logger.Level(zerolog.DebugLevel)

	for i := 0; i < b.N; i++ {
		testLog(logger)
	}
}

func BenchmarkZerolog(b *testing.B) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: io.Discard})
	logger.Level(zerolog.DebugLevel)

	for i := 0; i < b.N; i++ {
		testLog(logger)
	}
}

func BenchmarkNoPretty(b *testing.B) {
	logger := zerolog.New(io.Discard)
	logger.Level(zerolog.DebugLevel)

	for i := 0; i < b.N; i++ {
		testLog(logger)
	}
}
