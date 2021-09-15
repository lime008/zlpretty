package zlpretty_test

import (
	"os"

	"github.com/lime008/zlpretty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExampleLogger() {
	log.Logger = log.Output(zlpretty.ConsoleWriter{Out: os.Stderr})
	log.Logger.Level(zerolog.DebugLevel)

	type testStruct struct {
		Name string
		Age  int
		Foo  bool
	}

	testStructValue := &testStruct{
		Name: "John",
		Age:  32,
		Foo:  true,
	}

	log.Info().Str("foo", "bar").Msg("hello")
	log.Trace().Caller().Bool("boolean", false).Str("foo", "bar").Int("number", 14).Msg("hello")
	log.Debug().Msg("hi")
	log.Warn().Msg("WARNING")
	log.Info().Interface("testStruct", testStructValue).Msg("hello")

	testSlice := []struct {
		Name   string
		Age    int
		Foo    bool
		Nested *testStruct
	}{
		{
			Name:   "John",
			Age:    32,
			Foo:    true,
			Nested: testStructValue,
		},
		{
			Name: "Jane",
			Age:  21,
			Foo:  false,
		},
	}

	log.Error().Interface("testSlice", testSlice).Msg("This is a slice")

	// Output:
}
