package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"

	"github.com/newtechlab/vor/vorserve/data"
)

var (
	fHelp bool
	fHTTP string
	fData string
	fSalt string
)

var (
	globalStorage data.Storage
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&fHelp, "h", false, "show this information")
	flag.StringVar(&fHTTP, "http", ":5000", "interface and port to bind to")
	flag.StringVar(&fData, "data", "", "data storage path, supports local folder or S3 bucket, formatted as s3:bucketname or file:path")
	flag.StringVar(&fSalt, "salt", "", "salt to use, if not specified a random one is used")
}

func main() {
	flag.Parse()

	if fHelp {
		showHelp()
	}
	if fData == "" {
		showError("you must provide a value for the data flag")
	}
	if fSalt == "" {
		// initiate a random salt that will be used
		buf := make([]byte, 512/8)
		n, err := rand.Read(buf)
		if n != 512/8 || err != nil {
			log.Fatalln("error reading random data: ", err, n)
		}
		fSalt = base64.StdEncoding.EncodeToString(buf)
	}
	if len(fSalt) < 32 {
		log.Fatalln("to short a salt, must be at least 32 characters long")
	}

	setupStorage()
	registerHandlers()
	runServer()
}

func setupStorage() {
	var err error
	globalStorage, err = data.NewStorage(fData)
	if err != nil {
		log.Fatalln("error creating data storage: ", err)
	}
}
