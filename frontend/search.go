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
	"net/http"
	"sort"
	"strings"
)

type resultList []searchResult

type searchResult struct {
	NDBID        string
	FoodGroup    int
	Description  string
	Manufacturer string `json:",omitempty"`
	Score        int
}

func (l resultList) Len() int {
	return len(l)
}

func (l resultList) Less(i, j int) bool {
	return l[i].Score > l[j].Score
}

func (l resultList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (s *server) search(rw http.ResponseWriter, req *http.Request) {
	q := req.FormValue("q")
	terms := strings.Split(strings.ToLower(q), " ")

	// For each search term, start a new goroutine to search the BST. It is
	// threadsafe for reads/access.
	queries := make(chan []string)
	for _, term := range terms {
		go func(term string) {
			queries <- s.db.FindFood(term)
		}(term)
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
	results := make(resultList, len(scores))
	i := 0
	for id, score := range scores {
		food := s.db.Foods[id]
		results[i] = searchResult{
			NDBID:        food.NDBID,
			FoodGroup:    food.FoodGroup,
			Description:  food.LongDescription,
			Manufacturer: food.Manufacturer,
			Score:        score,
		}
		i++
	}
	sort.Sort(results)
	jsonResponse(rw, results)
}
