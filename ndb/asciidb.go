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
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/rsesek/usda-ndb/bst"
)

type ASCIIDB struct {
	basePath   string
	FoodGroups []FoodGroup
	Nutrients  []Nutrient
	Foods      map[string]*Food
	searchTree *bst.Tree
}

func ReadDatabase(base string) (*ASCIIDB, error) {
	db := &ASCIIDB{
		basePath:   base,
		Foods:      make(map[string]*Food, 8000),
		searchTree: bst.NewTree(),
	}

	log.Print("Loading food groups")
	if err := db.readFoodGroups(); err != nil {
		return nil, err
	}

	log.Print("Loading nutrient definitions")
	if err := db.readNutrientDefinitions(); err != nil {
		return nil, err
	}

	log.Print("Loading food database")
	if err := db.readFoods(); err != nil {
		return nil, err
	}

	log.Print("Loading food nutrients information")
	if err := db.readFoodNutrients(); err != nil {
		return nil, err
	}

	log.Print("Loading weight information")
	if err := db.readWeights(); err != nil {
		return nil, err
	}

	log.Print("Database loaded")
	log.Printf("... %d foods", len(db.Foods))

	return db, nil
}

// FindFood performs a text search for Foods named |name| and returns a slice of
// NDBIDs for matches, or nil on none.
func (db *ASCIIDB) FindFood(name string) []string {
	return db.searchTree.Find(name)
}

func (db *ASCIIDB) readFoodGroups() error {
	return ReadFile(path.Join(db.basePath, "FD_GROUP.txt"), func(line string) error {
		parts := strings.Split(line, "^")
		if len(parts) != 2 {
			return fmt.Errorf("Expected 2 parts, got %d from a FD_GROUP", len(parts))
		}
		code, err := intyString(parts[0])
		if err != nil {
			return fmt.Errorf("readFoodGroups: %v", err)
		}
		db.FoodGroups = append(db.FoodGroups, FoodGroup{
			GroupCode:   code,
			Description: trimString(parts[1]),
		})
		return nil
	})
}

func (db *ASCIIDB) readNutrientDefinitions() error {
	return ReadFile(path.Join(db.basePath, "NUTR_DEF.txt"), func(line string) error {
		parts := strings.Split(line, "^")
		if len(parts) != 6 {
			return fmt.Errorf("Expected 6 parts, got %d from a NUTR_DEF", len(parts))
		}

		id, err := intyString(parts[0])
		if err != nil {
			return fmt.Errorf("readNutrientDefinitions: %v", err)
		}

		db.Nutrients = append(db.Nutrients, Nutrient{
			NutrientID:  id,
			Units:       trimString(parts[1]),
			Description: trimString(parts[3]),
		})

		return nil
	})
}

func (db *ASCIIDB) readFoods() error {
	return ReadFile(path.Join(db.basePath, "FOOD_DES.txt"), func(line string) error {
		parts := strings.Split(line, "^")
		if len(parts) != 14 {
			return fmt.Errorf("Expected 14 parts, got %d from a FOOD_DES", len(parts))
		}

		foodGroup, err := intyString(parts[1])
		if err != nil {
			return fmt.Errorf("readFoods: FoodGroup: %v", err)
		}

		var refuse int
		if s := trimString(parts[8]); s != "" {
			refuse, err = intyString(s)
			if err != nil {
				return fmt.Errorf("readFoods: Refuse: %v", err)
			}
		}

		id := trimString(parts[0])

		food := &Food{
			NDBID:             id,
			FoodGroup:         foodGroup,
			LongDescription:   trimString(parts[2]),
			ShortDescription:  trimString(parts[3]),
			CommonNames:       trimString(parts[4]),
			Manufacturer:      trimString(parts[5]),
			RefuseDescription: trimString(parts[7]),
			Refuse:            refuse,
		}
		db.Foods[id] = food

		// Join all the descriptions together to create search terms.
		search := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
			food.LongDescription, food.ShortDescription, food.CommonNames, food.Manufacturer))
		var last int
		for i := 0; i < len(search); i++ {
			c := search[i]
			if c == ',' || c == ' ' || c == '&' || c == '/' || c == '!' || c == '-' || c == '.' {
				part := search[last:i]
				if len(part) > 2 {
					db.searchTree.Insert(bst.Pair{part, id})
				}
				last = i + 1
			}
		}

		return nil
	})
}

func (db *ASCIIDB) readFoodNutrients() error {
	return ReadFile(path.Join(db.basePath, "NUT_DATA.txt"), func(line string) error {
		parts := strings.Split(line, "^")
		if len(parts) != 18 {
			return fmt.Errorf("Expected 18 parts, got %d from a NUT_DATA", len(parts))
		}

		id := trimString(parts[0])

		food, ok := db.Foods[id]
		if !ok {
			return fmt.Errorf("readFoodNutrients: Could not find food %s", id)
		}

		nutrientID, err := intyString(parts[1])
		if err != nil {
			return fmt.Errorf("readFoodNutrients: NutrientID: %v", err)
		}

		value, err := strconv.ParseFloat(trimString(parts[2]), 32)
		if err != nil {
			return fmt.Errorf("readFoodNutrients: Value: %v", err)
		}

		dataPoints, err := intyString(parts[3])
		if err != nil {
			return fmt.Errorf("readFoodNutrients: NutrientID: %v", err)
		}

		food.Nutrients = append(food.Nutrients, FoodNutrient{
			NutrientID: nutrientID,
			Value:      float32(value),
			DataPoints: dataPoints,
		})
		return nil
	})
}

func (db *ASCIIDB) readWeights() error {
	return ReadFile(path.Join(db.basePath, "WEIGHT.txt"), func(line string) error {
		parts := strings.Split(line, "^")
		if len(parts) != 7 {
			return fmt.Errorf("Expected 7 parts, got %d from a WEIGHT", len(parts))
		}

		id := trimString(parts[0])
		food, ok := db.Foods[id]
		if !ok {
			return fmt.Errorf("readWeights: Could not find food %s", id)
		}

		sequence, err := intyString(parts[1])
		if err != nil {
			return fmt.Errorf("readWeights: Sequence: %v", err)
		}

		amount, err := strconv.ParseFloat(trimString(parts[2]), 32)
		if err != nil {
			return fmt.Errorf("readWeights: Amount: %v", err)
		}

		weight, err := strconv.ParseFloat(trimString(parts[4]), 32)
		if err != nil {
			return fmt.Errorf("readWeights: WeightG: %v", err)
		}

		food.Weights = append(food.Weights, Weight{
			Sequence: sequence,
			Amount: float32(amount),
			Description: trimString(parts[3]),
			WeightG: float32(weight),
		})
		return nil
	})
}

// intyString turns a stringified number in the ASCII database dump format into an actual int.
func intyString(a string) (int, error) {
	return strconv.Atoi(trimString(a))
}

func trimString(s string) string {
	if s == "~~" {
		return ""
	}
	return strings.Trim(s, "~")
}
