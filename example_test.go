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

	testStruct := struct {
		Name string
		Age  int
		Foo  bool
	}{
		Name: "John",
		Age:  32,
		Foo:  true,
	}

	log.Debug().Msg("hi")
	log.Debug().Msg("")
	log.Debug().Str("foo", "bar").Int("number", 14).Msg("hello")
	log.Info().Str("foo", "bar").Msg("hello")
	log.Error().Interface("testStruct", testStruct).Msg("hello")

	// Output:
}
