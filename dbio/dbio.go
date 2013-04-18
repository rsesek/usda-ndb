package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rsesek/usda-ndb/ndb"
)

var (
	asciidb = flag.String("asciidb", "", "The path to the ASCII database dumps.")
	output  = flag.String("output", "asciidb.gob.gz", "The path to the output file.")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *asciidb == "" {
		fmt.Fprintln(os.Stderr, "No -asciidb specified")
		flag.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "No -output specified")
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Reading database from %s", *asciidb)
	db, err := ndb.ReadDatabase(*asciidb)
	if err != nil {
		log.Fatalf("ndb.ReadDatabase: %v", err)
	}

	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer f.Close()

	log.Print("Writing compressed stream to %s", *output)

	w := gzip.NewWriter(f)
	defer w.Close()

	enc := gob.NewEncoder(w)
	if err := enc.Encode(db); err != nil {
		log.Fatalf("gob.Encode: %v", err)
	}

	log.Print("***** Done *****")
}
