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

package bst

import (
	"reflect"
	"testing"
)

// InOrderTokens is used to test everything else, so make sure it works manually.
func TestInOrderTokens(t *testing.T) {
	/*
						8
				3				10
		1			6					14
				4		7			13
	*/
	tree := NewTree()
	tree.root = &Node{
		value:  "8",
		tokens: []string{"g"},
		left: &Node{
			value:  "3",
			tokens: []string{"b"},
			left: &Node{
				value:  "1",
				tokens: []string{"a"},
			},
			right: &Node{
				value:  "6",
				tokens: []string{"d", "e"},
				left: &Node{
					value:  "4",
					tokens: []string{"c"},
				},
				right: &Node{
					value:  "7",
					tokens: []string{"f"},
				},
			},
		},
		right: &Node{
			value:  "10",
			tokens: []string{"h"},
			right: &Node{
				value:  "14",
				tokens: []string{"k", "l"},
				left: &Node{
					value:  "13",
					tokens: []string{"i", "j"},
				},
			},
		},
	}

	expected := "abcdefghijkl"
	var actual string
	for v := range tree.InOrderTokens() {
		actual += v
	}
	if expected != actual {
		t.Errorf("Expected traversal to get %q, got %q", expected, actual)
	}
}

func TestInsertRoot(t *testing.T) {
	tree := NewTree()
	tree.Insert(Pair{Value: "hello", Token: "world"})
	actual := <-tree.InOrderTokens()
	expected := "world"
	if actual != expected {
		t.Errorf("Root should have token value %q, got %q", expected, actual)
	}
}

func TestInsert(t *testing.T) {
	pairs := []Pair{
		{"moo", "cow"},
		{"baaa", "sheep"},
		{"crow", "crow"},
		{"bark", "dog"},
		{"bark", "chipmunk"},
		{"hoot", "owl"},
	}
	tree := NewTree()
	for _, p := range pairs {
		tree.Insert(p)
	}

	expected := []string{
		"sheep",
		"dog",
		"chipmunk",
		"crow",
		"owl",
		"cow",
	}
	i := 0

	for actual := range tree.InOrderTokens() {
		if i >= len(expected) {
			t.Errorf("Got unexpected %q", actual)
			continue
		}
		if actual != expected[i] {
			t.Errorf("Token %d should be %q, got %q", expected[i], actual)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Missing %d tokens from stream", len(expected)-i)
	}
}

func TestFind(t *testing.T) {
	pairs := []Pair{
		{"moo", "cow"},
		{"baaa", "sheep"},
		{"crow", "crow"},
		{"bark", "dog"},
		{"bark", "chipmunk"},
		{"hoot", "owl"},
	}
	tree := NewTree()
	for _, p := range pairs {
		tree.Insert(p)
	}

	expectations := []struct {
		value  string
		tokens []string
	}{
		{"moo", []string{"cow"}},
		{"mooo", nil},
		{"baaa", []string{"sheep"}},
		{"crow", []string{"crow"}},
		{"CROW", nil},
		{"bark", []string{"dog", "chipmunk"}},
		{"hoot", []string{"owl"}},
		{"hot", nil},
	}
	for _, expected := range expectations {
		actual := tree.Find(expected.value)
		if !reflect.DeepEqual(expected.tokens, actual) {
			t.Errorf("When finding %q, expected %v, got %v", expected.tokens, actual)
		}
	}
}
