package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	process "github.com/atyronesmith/gennextgen/pkg/process"
	genUtils "github.com/atyronesmith/gennextgen/pkg/utils"
)

func main() {
	isHelp := flag.Bool("help", false, "Print usage information.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] dirPath", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       path/tripleo...  -- Path to tripleo-overcloud-envrionment.yaml file.\n")

		flag.PrintDefaults()
	}

	flag.Parse()

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if flag.NArg() != 1 && flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	toeFile := ""
	tordFile := ""
	tobdFile := ""

	if flag.NArg() == 1 {
		fileInfo, err := os.Stat(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !fileInfo.IsDir() {
			fmt.Println("Error: Path is not a directory.")
			os.Exit(1)
		}
		fmt.Printf("Directory: %s\n", flag.Arg(0))
		genUtils.SetRootDir(flag.Arg(0))
	} else {
		toeFile = flag.Arg(0)
		tordFile = flag.Arg(1)
		tobdFile = flag.Arg(2)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Current working directory: %s\n", cwd)

	tripleoData, err := process.GetTripleoData(toeFile, tordFile, tobdFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s, err := genUtils.StructToYaml(tripleoData.Environment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(s))

	re := regexp.MustCompile(`(?m).*(Parameters|HostnameFormat|NetworkConfigTemplate|ExtraGroupVars|ExtraConfig|Count|ControlPlaneSubnet|Services)`)

	for key, value := range tripleoData.Environment.ParmaterDefaults.Params {
		if re.MatchString(key) {
			fmt.Printf("Key: %s, Value: %v\n", key, value)
		}
	}

	s, err = genUtils.StructToYaml(tripleoData.Roles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(s))

	s, err = genUtils.StructToYaml(tripleoData.Deployment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(s))
}
