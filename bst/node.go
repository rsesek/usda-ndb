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

type Node struct {
	left   *Node
	right  *Node
	value  string
	tokens []string
}

func MakeNode(p Pair) *Node {
	n := &Node{value: p.Value}
	n.InsertPair(p)
	return n
}

func (n *Node) Value() string {
	return n.value
}

func (n *Node) InsertPair(p Pair) {
	if n.value != p.Value {
		panic("Cannot insert a Pair into a Node that does not match Value")
	}
	n.tokens = append(n.tokens, p.Token)
}

type Pair struct {
	Value string
	Token string
}

func (p *Pair) Less(other Pair) bool {
	return p.Value < other.Value
}

func (p *Pair) Equal(other Pair) bool {
	return p.Value == other.Value
}
