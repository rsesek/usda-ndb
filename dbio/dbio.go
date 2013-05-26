//
// USDA-NDB Viewer
// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

/*
Command dbio reads an ASCII database into a github.com/rsesek/usda-ndb/ndb.ASCII object
and writes it back out as compressed GOB file.
*/

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
