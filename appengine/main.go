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

package appengine

import (
	"compress/gzip"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/rsesek/usda-ndb/frontend"
	"github.com/rsesek/usda-ndb/ndb"
)

func init() {
	f, err := os.Open("./asciidb.gob.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	var db *ndb.ASCIIDB
	dec := gob.NewDecoder(r)
	if err = dec.Decode(&db); err != nil {
		panic(err)
	}

	server := frontend.NewServer(db, "./static")
	http.Handle("/", server)
}
