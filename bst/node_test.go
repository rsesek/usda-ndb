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
	"testing"
)

func TestMakeNode(t *testing.T) {
	n := MakeNode(Pair{Value: "hello", Token: "world"})
	expected := "hello"
	if n.value != expected {
		t.Errorf("Expected Value to be %q, got %q", expected, n.value)
	}
	if len(n.tokens) != 1 {
		t.Errorf("Expected 1 Tokens, got %d", len(n.tokens))
	} else {
		expected = "world"
		if n.tokens[0] != expected {
			t.Errorf("Expected Token to be %q, got %q", expected, n.tokens[0])
		}
	}
}

func TestNodeInsertPair(t *testing.T) {
	const kValue = "hello"
	node := &Node{value: kValue}
	node.InsertPair(Pair{Value: kValue, Token: "doc1"})
	node.InsertPair(Pair{Value: kValue, Token: "doc2"})
	node.InsertPair(Pair{Value: kValue, Token: "doc1"})
	if len(node.tokens) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(node.tokens))
	}

	expectations := []struct {
		index    int
		expected string
	}{
		{0, "doc1"},
		{1, "doc2"},
		{2, "doc1"},
	}
	for i := 0; i < len(expectations); i++ {
		expected := expectations[i]
		actual := node.tokens[expected.index]
		if actual != expected.expected {
			t.Errorf("At token index %d, expected %q, got %q", expected.index, expected.expected, actual)
		}
	}
}

func TestPairCompare(t *testing.T) {
	p1 := Pair{Value: "abc"}
	p2 := Pair{Value: "bcd"}

	if c := p1.Less(p2); !c {
		t.Errorf("p1 comes before p2, got %t", c)
	}
	if c := p1.Equal(p2); c {
		t.Errorf("p1 is not equal to p2")
	}

	p2.Value = "abb"
	if c := p1.Less(p2); c {
		t.Errorf("p1 comes after p2")
	}

	p1.Value = "abd"
	if c := p1.Less(p2); c {
		t.Errorf("p1 comes after p2, got %t", c)
	}

	p2.Value = "abd"
	if c := p1.Less(p2); c {
		t.Errorf("p1 is the same as p2, not less, got %d", c)
	}
	if c := p1.Equal(p2); !c {
		t.Errorf("p1 is the same as p2")
	}
}
