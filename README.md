# Go PINE

Implementation of [PINE](https://github.com/GovanifY/pine) client in Go language.

This is currently a WIP (work-in-progress).

## Goals

 * Plain Go
 * No dependencies to other modules except Go standard libraries

## How to use

Install `go-pine` module.

```
go get github.com/Made-in-Slovakia/go-pine@latest
```

Reading basic information as `Status`, `Version` and `Title` from emulator.

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

Same code but using `Exchange`.

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

	e := pine.NewExchange(client, 2)
	e.AddCommand(pine.StatusCommand())
	e.AddCommand(pine.VersionCommand())
	e.AddCommand(pine.TitleCommand())
	err := e.Execute()
	if err != nil {
		log.Fatal(err)
	}

	status, err := e.ReadUint32(0)
	if err != nil {
		log.Printf("status is %d", status)
	}

	version, err := e.ReadString(1)
	if err != nil {
		log.Printf("version is %s", version)
	}

	title, err := e.ReadString(2)
	if err != nil {
		log.Printf("title is %s", title)
	}

}
```

## TODOs

- [ ] Resolve TODOs in code
- [ ] Test implementation of Unix sockets for Linux and MacOS version
- [ ] Sockets for Linux and MacOS version should use slot variable
- [ ] Review from experienced Go developer
