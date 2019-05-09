package main

import (
	"flag"
	"fmt"
	"os"
)

func showHelp() {
	fmt.Println("vorserve: used to run a VOice Recording SERVEr")
	fmt.Println("")
	flag.PrintDefaults()
	os.Exit(1)
}

func showError(msg string) {
	fmt.Println("error: " + msg)
	fmt.Println("")
	fmt.Println("use voserve -h to show usage information")
	os.Exit(1)
}
