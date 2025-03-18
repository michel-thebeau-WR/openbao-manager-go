package main

import (
	"fmt"
	"os"

	baoConfig "github.com/michel-thebeau-WR/openbao-manager-go/baomon/config"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Print("No other functionality is supported at the moment.\n")
		return
	}

	args := os.Args[1:]
	if args[0] == "--read" || args[0] == "-r" {
		if len(args) < 2 {
			fmt.Fprint(os.Stderr, "Invalid number of arguments for --read/-r\n")
			return
		}

		var S baoConfig.MonitorConfig
		openedFile, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file %v. Error message: %v\n", args[1], err)
			return
		}
		defer openedFile.Close()

		err = S.ReadYAMLMonitorConfig(openedFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with parsing read yaml data: %v\n", err)
			return
		}
		fmt.Printf("Result: \n%#v\n", S)
	} else if args[0] == "--write" || args[0] == "-w" {
		if len(args) < 3 {
			fmt.Fprint(os.Stderr, "Invalid number of arguments for --read/-r\n")
			return
		}

		var S baoConfig.MonitorConfig
		openedFile, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file %v. Error message: %v\n", args[1], err)
			return
		}
		defer openedFile.Close()

		err = S.ReadYAMLMonitorConfig(openedFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with parsing read yaml data: %v\n", err)
			return
		}
		fmt.Printf("Read result: \n%#v\n", S)

		writeFile, err := os.Create(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %v for writing\n", args[2])
			return
		}
		defer writeFile.Close()
		err = S.WriteYAMLMonitorConfig(writeFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with writing parsed yaml data to file: %v\n", err)
			return
		}
		fmt.Print("Write Complete\n")
	} else {
		fmt.Print("No other functionality is supported at the moment.\n")
		return
	}
}
