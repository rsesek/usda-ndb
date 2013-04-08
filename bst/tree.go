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

// Package bst implements a Binary Search Tree. The tree node values are actually
// 2-tuples (search-key, app-token), with the search-key being the value that all
// values in a node have in common. The app-token is used to refer back to some
// object that is being searched for using the tree.
package bst

type Tree struct {
	root *Node
}

func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) Find(value string) []string {
	node := t.findNode(value, t.root)
	return nil
}

func (t *Tree) findNode(value string, node *Node) *Node {
	if node == nil || node.value == value {
		return node
	}

	if value < node.value {
		return t.findNode(value, node.left)
	}
	return t.findNode(value, node.right)
}

func (t *Tree) Insert(p Pair) {
	if t.root == nil {
		t.root = MakeNode(p)
	} else {
		t.insertOn(p, t.root)
	}
}

func (t *Tree) insertOn(p Pair, n *Node) {
	if p.Value == n.value {
		n.InsertPair(p)
	} else if p.Value < n.value {
		if n.left == nil {
			n.left = MakeNode(p)
		} else {
			t.insertOn(p, n.left)
		}
	} else {
		if n.right == nil {
			n.right = MakeNode(p)
		} else {
			t.insertOn(p, n.right)
		}
	}
}
