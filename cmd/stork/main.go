package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/huboleo/storkey/internal"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "save":
		runSave(os.Args[2:])
	case "pull":
		runPull(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func runSave(args []string) {
	fs := flag.NewFlagSet("save", flag.ContinueOnError)
	deleteFiles := fs.Bool("d", false, "delete env files after saving")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "save does not take positional arguments")
		os.Exit(2)
	}

	if err := internal.Save(*deleteFiles); err != nil {
		fmt.Fprintf(os.Stderr, "save failed: %v\n", err)
		os.Exit(1)
	}
}

func runPull(args []string) {
	fs := flag.NewFlagSet("pull", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "pull does not take positional arguments")
		os.Exit(2)
	}

	if err := internal.Pull(); err != nil {
		fmt.Fprintf(os.Stderr, "pull failed: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  stork save [-d]")
	fmt.Fprintln(os.Stderr, "  stork pull")
}
