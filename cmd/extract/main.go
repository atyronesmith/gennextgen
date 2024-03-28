package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/generate"
)

func main() {
	isVerbose := flag.Bool("verbose", false, "Print extra runtime information.")
	isHelp := flag.Bool("help", false, "Print usage information.")
	outDir := flag.String("outdir", "/tmp/", "Where to place output.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] stack_dir\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       stack_dir  -- Path to directory containing config-download stack.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isVerbose {
		fmt.Println("Verbose...")
	}

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if dirname := flag.Arg(0); dirname == "" {
		flag.Usage()
		os.Exit(1)
	} else {
		if err := generate.GenSecrets(dirname, *outDir); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		if err := generate.GenOpenStackControlPlane(dirname, *outDir); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

	}

}
