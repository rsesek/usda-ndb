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

package ndb

import (
	"fmt"
	"strings"

	"github.com/rsesek/usda-ndb/bst"
)

func (db *ASCIIDB) addTermsForFood(food *Food) {
	// Join all the descriptions together to create search terms.
	search := strings.ToLower(fmt.Sprintf("%s %s %s %s",
		food.LongDescription, food.ShortDescription, food.CommonNames, food.Manufacturer))
	var last int
	for i := 0; i < len(search); i++ {
		c := search[i]
		if c == ',' || c == ' ' || c == '&' || c == '/' || c == '!' || c == '-' || c == '.' {
			part := search[last:i]
			if len(part) > 2 {
				db.searchTree.Insert(bst.Pair{part, food.NDBID})
			}
			last = i + 1
		}
	}
}

func (db *ASCIIDB) RebuildSearchIndex() {
	db.searchTree = bst.NewTree()
	for _, food := range db.Foods {
		db.addTermsForFood(food)
	}
}
