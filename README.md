# zlpretty

### Forked from [github.com/UnnoTed/horizontal](https://github.com/UnnoTed/horizontal)

zlpretty is a pretty logging with focus on readability, based on the `zerolog.ConsoleWriter` but with some added features like json pretty printing and log line separator.

![zlpretty](https://i.imgur.com/RvuuYSj.png)

`go get -u github.com/lime008/zlpretty`

```go
package main

import (
	"os"

	"github.com/lime008/zlpretty"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(horizontal.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("hi")
	log.Debug().Msg("hello")
}

```
