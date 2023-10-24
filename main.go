package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	a "github.com/dissipative/opabinia/internal/app"
	"github.com/jessevdk/go-flags"
)

//go:embed VERSION
var version string

type Options struct {
	Init    bool `short:"i" long:"init" description:"Init markdown project"`
	Serve   bool `short:"s" long:"serve" description:"Serve serving project"`
	Compile bool `short:"c" long:"compile" description:"Compile project to static html site"`
	Version bool `short:"v" long:"version" description:"Show a version"`
}

func (o *Options) empty() bool {
	return !(o.Init || o.Serve || o.Compile || o.Version)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from", r)
			os.Exit(1)
		}
	}()

	app, err := a.NewApp()
	if err != nil {
		fatal("Fatal error", err)
	}

	opts := new(Options)

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = a.Name
	parser.Usage = "[OPTIONS]"

	_, err = parser.ParseArgs(os.Args[1:])
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		fmt.Println("Error occurred: ", err)
		parser.WriteHelp(bufio.NewWriter(os.Stdout))
		os.Exit(1)
	}

	if opts.empty() {
		parser.WriteHelp(bufio.NewWriter(os.Stdout))
		os.Exit(0)
	}

	switch {
	case opts.Version:
		fmt.Print(version)
	case opts.Init:
		if err = app.DoInit(); err != nil {
			fatal("Init fatal error", err)
		}
	case opts.Compile:
		if err = app.Compile(); err != nil {
			fatal("Compile fatal error", err)
		}
	case opts.Serve:
		if err = app.DoServe(); err != nil {
			fatal("Serve fatal error", err)
		}
	}
}

func fatal(msg string, problem any) {
	fmt.Printf("%s: %s\n", msg, problem)
	os.Exit(1)
}
