# Go PINE

Implementation of [PINE](https://github.com/GovanifY/pine) client in Go language.

This is currently a WIP (work-in-progress).

## Goals

 * Plain Go
 * No dependencies to modules ouside of Go core

## How to use

```
go get github.com/Made-in-Slovakia/go-pine@latest
```

```
package main

import (
	pine "github.com/Made-in-Slovakia/go-pine"
	"log"
)

func main() {
	pine.DebugLogEnabled = true
	client := pine.NewClient(28011)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	answers, err := client.SendCommands([]pine.Command{
		pine.StatusCommand(),
		pine.VersionCommand(),
		pine.TitleCommand(),
	})

	if err != nil {
		log.Fatal(err)
	}

	for i, a := range answers {
		log.Printf("answer %v is %v", i, a)
	}

}
```

## TODOs

- [ ] TODOs in code
- [ ] Review from experienced Go developer
