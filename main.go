package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/YakDriver/fixit/fixers/verbfix"
)

func main() {
	force := flag.Bool("f", false, "Overwrite input file")
	new := flag.Bool("n", false, "Write new file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("must include file to fix! got none")
	}

	filename := flag.Args()[0]
	content, err := verbfix.FileContent(filename)
	if err != nil {
		log.Fatal(err)
	}

	if *force {
		err = verbfix.OverwriteFile(filename, verbfix.FixIt(content))
		if err != nil {
			log.Fatal(err)
		}
	} else if *new {
		newFilename := strings.TrimSuffix(filename, path.Ext(filename)) + "_new" + path.Ext(filename)
		err = verbfix.WriteFile(newFilename, verbfix.FixIt(content))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%s", verbfix.FixIt(content))
	}
}
