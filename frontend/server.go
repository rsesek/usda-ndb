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

package frontend

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rsesek/usda-ndb/ndb"
)

var (
	debug = flag.Bool("debug", false, "Debug mode: log all requests")
)

// NewServer creates a HTTP Handler that will serve static files from staticDir and
// various API endpoints using the ASCIIDB db.
func NewServer(db *ndb.ASCIIDB, staticDir string) http.Handler {
	s := &server{
		db:        db,
		staticDir: staticDir,
		mux:       http.NewServeMux(),
	}
	s.init()
	return s
}

type server struct {
	db        *ndb.ASCIIDB
	staticDir string
	mux       *http.ServeMux
}

func (s *server) init() {
	s.mux.Handle("/", http.FileServer(http.Dir(s.staticDir)))
	s.handleMethod("/_/search", (*server).search)
	s.handleMethod("/_/foodGroups", (*server).foodGroups)
	s.handleMethod("/_/nutrients", (*server).nutrients)
	s.handleMethod("/_/food/", (*server).getFood)
}

// Convience method to work around https://code.google.com/p/go/issues/detail?id=2280.
// This can be removed after Go1.1.
func (s *server) handleMethod(url string, meth func(*server, http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(url, func(rw http.ResponseWriter, req *http.Request) {
		meth(s, rw, req)
	})
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if *debug {
		log.Printf("%s %s %s", req.Proto, req.Method, req.URL)
	}

	s.mux.ServeHTTP(rw, req)
}

func (s *server) search(rw http.ResponseWriter, req *http.Request) {
	q := req.FormValue("q")
	terms := strings.Split(q, " ")

	// For each search term, start a new goroutine to search the BST. It is
	// threadsafe for reads/access.
	queries := make(chan []string)
	for _, term := range terms {
		go func() {
			queries <- s.db.FindFood(term)
		}()
	}

	// A realllllly stupid scoring algorithm just counts the number of times
	// the NDBID comes up.
	scores := make(map[string]int)
	for i := 0; i < len(terms); i++ {
		for _, id := range <-queries {
			scores[id]++
		}
	}

	// Collect the results into a response list.
	// TODO(rsesek): Trim out nutrient data? Paginate?
	foods := make([]*ndb.Food, len(scores))
	i := 0
	for id, _ := range scores {
		foods[i] = s.db.Foods[id]
		i++
	}
	jsonResponse(rw, scores)
}

func (s *server) foodGroups(rw http.ResponseWriter, req *http.Request) {
	jsonResponse(rw, s.db.FoodGroups)
}

func (s *server) nutrients(rw http.ResponseWriter, req *http.Request) {
	jsonResponse(rw, s.db.Nutrients)
}

func (s *server) getFood(rw http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")
	id := parts[len(parts)-1]
	if food, ok := s.db.Foods[id]; ok {
		jsonResponse(rw, food)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Error: Could not find food with id %s", id)
	}
}

func jsonResponse(rw http.ResponseWriter, resp interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		panic(err.Error())
	}
}
