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
)

type ASCIIDB struct {
	basePath   string
	FoodGroups []FoodGroup
	Nutrients  []Nutrient
	Foods      map[string]Food
}

func ReadDatabase(base string) (*ASCIIDB, error) {
	db := &ASCIIDB{
		basePath: base,
		Foods:    make(map[string]Food, 8000),
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

	log.Print("Database loaded")

	fmt.Printf("%#v", *db)
	return db, nil
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

		refuse, err := intyString(parts[8])
		if err != nil {
			return fmt.Errorf("readFoods: Refuse: %v", err)
		}

		id := trimString(parts[0])

		db.Foods[id] = Food{
			NDBID:             id,
			FoodGroup:         foodGroup,
			LongDescription:   trimString(parts[2]),
			ShortDescription:  trimString(parts[3]),
			CommonNames:       trimString(parts[4]),
			Manufacturer:      trimString(parts[5]),
			RefuseDescription: trimString(parts[7]),
			Refuse:            refuse,
		}

		return nil
	})
}

func (db *ASCIIDB) readFoodNutrients() error {
	return ReadFile(path.Join(db.basePath, "NUT_DATA.txt"), func(line string) error {
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
