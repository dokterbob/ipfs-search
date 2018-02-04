/*

Search engine for IPFS using Elasticsearch, RabbitMQ and Tika.
*/
package main

import (
	"context"
	"fmt"
	"github.com/ipfs-search/ipfs-search/commands"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Prefix logging with filename and line number: "d.go:23"
	// log.SetFlags(log.Lshortfile)

	// Logging w/o prefix
	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "ipfs-search"
	app.Usage = "IPFS search engine."

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add `HASH` to crawler queue",
			Action:  add,
		},
		{
			Name:    "crawl",
			Aliases: []string{"c"},
			Usage:   "start crawler",
			Action:  crawl,
		},
	}

	app.Run(os.Args)
}

func add(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Please supply one hash as argument.", 1)
	}

	hash := c.Args().Get(0)

	fmt.Printf("Adding hash '%s' to queue\n", hash)

	err := commands.AddHash(hash)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

// onSigTerm calls f() when SIGTERM (control-C) is received
func onSigTerm(f func()) {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var fail = func() {
		<-sigChan
		os.Exit(1)
	}

	var quit = func() {
		<-sigChan

		go fail()

		fmt.Println("Received SIGTERM, quitting... One more SIGTERM and we'll abort!")
		f()
	}

	go quit()
}

func crawl(c *cli.Context) error {
	fmt.Printf("Starting worker\n")

	ctx, cancel := context.WithCancel(context.Background())

	// Allow SIGTERM / Control-C quit through context
	onSigTerm(cancel)

	err := commands.Crawl(ctx)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
