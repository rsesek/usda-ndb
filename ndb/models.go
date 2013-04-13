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

// A Nutrient represents either a macro or micronutrient that is measured
// for a food item in the database.
type Nutrient struct {
	// 3-digit code that identifiers the nutrient. Key.
	NutrientID int
	// The units of measure.
	Units string
	// Description.
	Description string
	// The order used in official reports.
	SortOrder int
}

type FoodGroup struct {
	// 4-digit code identifying the food group.
	GroupCode int
	// Name of the food group.
	Description string
}

// A Food represents a foodstuff whose nutritional content has been measured in the NDB.
type Food struct {
	// 5-digit identification number for the food. Key.
	// String to preserve leading zeros (apparently).
	NDBID string
	// 4-digit code indicating food group to which this food belongs.
	FoodGroup int
	// Long (200 chars) description.
	LongDescription string
	// Short (60 chars) description.
	ShortDescription string
	// Other names used to describe the food.
	CommonNames string
	// Scientific name of the food.
	ScientificName string
	// If applicable, the manufacturer of the food.
	Manufacturer string
	// Description of the inedible parts of the food.
	RefuseDescription string
	// The percentage of the food that is refuse.
	Refuse int
	// Nutrients of the food.
	Nutrients []FoodNutrient
	// The common household weights/units.
	Weights []Weight
}

// A FoodNutrieint is a measured nutrient value for a food item.
type FoodNutrient struct {
	// 3-digit code that identifiers the nutrient. Key.
	NutrientID int
	// Edible portion (amount in 100 grams).
	Value float32
	// Number of data points used to calculate the value.
	DataPoints int
}

// A Weight is a common measure of a food item that contains a factor for
// multiplying a FoodNutrient.Value to get the Value in common units.
//
// The formula is N = (V*W) / 100
//   N = nutrient value for household measure.
//   V = nutrient value per 100g.
//   W = g weight of portion.
type Weight struct {
	// Sequence number.
	Sequence int
	// Amount of units (e.g. 1 in 1 cup).
	Amount float32
	// The description of the weight.
	Description string
	// The weight in grams for this unit.
	WeightG float32
}
