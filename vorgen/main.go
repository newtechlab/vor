package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/newtechlab/vor/vorgen/config"
)

var (
	fConfig     string
	fDumpConfig bool
	fHelp       bool
)

func init() {
	flag.StringVar(&fConfig, "config", "./config.json", "path to config file")
	flag.BoolVar(&fHelp, "h", false, "show this info")
	flag.BoolVar(&fDumpConfig, "dump", false, "dump a sample config file")
}

func main() {
	flag.Parse()

	if fHelp {
		showHelp()
	}
	if fDumpConfig {
		showConfig()
	}

	config := loadConfig()
	project := generateTwillioProject(config)
	dumpTwillio(project)
}

func showHelp() {
	fmt.Println("vorgen: generate Twiliio project for voice recording")
	fmt.Println("")
	flag.PrintDefaults()
	os.Exit(0)
}

func showConfig() {
	fmt.Println(config.Default())
	os.Exit(0)
}

func loadConfig() config.Config {
	f, err := os.Open(fConfig)
	if err != nil {
		log.Fatalln("error loading config: ", err)
	}
	defer f.Close()
	c, err := config.LoadConfig(bufio.NewReader(f))
	if err != nil {
		log.Fatalln("error loading config from "+fConfig+": ", err)
	}
	return c
}

func dumpTwillio(a interface{}) {
	bstd := bufio.NewWriter(os.Stdout)
	defer bstd.Flush()
	enc := json.NewEncoder(bstd)
	enc.SetIndent("", "  ")
	err := enc.Encode(a)
	if err != nil {
		log.Fatalln("error encoding twillio JSON: ", err)
	}
}
