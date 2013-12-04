// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// This file tests the tree package.

package tree

import (
	"testing"
)

// test is a struct that describes a test case with its properties.
// In general most properties require that the tree was not initialized without
// randomization.
type test struct {
	vs []Value
	fn func(a, b Value) bool

	// The following fields do not depend on the insertion order
	siz     int
	max     Value
	min     Value
	str     string
	sorted  []Value
	missing []Value
	remOrd  []Value

	// height depends on the insertion order.
	root   Value
	height int
}

var (
	t1 = test{
		vs:      []Value{5, 2, 7, 3, 1, 6, 9, 4, 8, 2, 1, 8},
		siz:     9,
		max:     1,
		min:     9,
		str:     "[9 8 7 6 5 4 3 2 1]",
		sorted:  []Value{9, 8, 7, 6, 5, 4, 3, 2, 1},
		missing: []Value{0, 10, -5, 45, 347, -1},
		remOrd:  []Value{6, 8, 3, 1, 5},
		root:    5,
		height:  4,
		fn: func(a, b Value) bool {
			ai := a.(int)
			bi := b.(int)
			return ai > bi
		},
	}

	t2 = test{
		vs:      []Value{1, 2, 3, 4, 5, 6, 7, 8, 9},
		siz:     9,
		max:     9,
		min:     1,
		str:     "[1 2 3 4 5 6 7 8 9]",
		sorted:  []Value{1, 2, 3, 4, 5, 6, 7, 8, 9},
		missing: []Value{0, 10, -5, 45, 347, -1},
		remOrd:  []Value{2, 4, 1, 3},
		root:    1,
		height:  9,
		fn: func(a, b Value) bool {
			ai := a.(int)
			bi := b.(int)
			return ai < bi
		},
	}

	t3 = test{
		vs:      []Value{67.9, -1.5, 4e8, 567.34, -567.89, 0.0, 0.0},
		siz:     6,
		max:     4e8,
		min:     -567.89,
		str:     "[-567.89 -1.5 0 67.9 567.34 4e+08]",
		sorted:  []Value{-567.89, -1.5, 0.0, 67.9, 567.34, 4.0e8},
		missing: []Value{67.99, 0.00001, -1.6, 1.0, 400.0},
		remOrd:  []Value{67.9, 0.0},
		root:    67.9,
		height:  3,
		fn: func(a, b Value) bool {
			af := a.(float64)
			bf := b.(float64)
			return af < bf
		},
	}

	t4 = test{
		vs:      []Value{"Lisa", "Lukas", "Ben", "Chris", "Chris", "Benni", "Sara", "Patrick"},
		siz:     7,
		max:     "Sara",
		min:     "Ben",
		str:     "[Ben Benni Chris Lisa Lukas Patrick Sara]",
		sorted:  []Value{"Ben", "Benni", "Chris", "Lisa", "Lukas", "Patrick", "Sara"},
		missing: []Value{"Dan", "Benjamin", "Christopher", "Marietta", "Wolfgang", "Ruth"},
		remOrd:  []Value{"Patrick", "Sara", "Lisa", "Lukas"},
		root:    "Lisa",
		height:  4,
		fn: func(a, b Value) bool {
			as := a.(string)
			bs := b.(string)
			return as < bs
		},
	}

	tests = []test{t1, t2, t3, t4}
)

func TestTree(o *testing.T) {
	for _, want := range tests {
		equals := func(a, b Value) bool {
			return !want.fn(a, b) && !want.fn(b, a)
		}

		tree := New(want.fn)
		tree.Init(want.vs)

		if siz := tree.Len(); siz != want.siz {
			o.Errorf("tree.Len() = %v; want %v", siz, want.siz)
		}
		if max := tree.Max().Val(); !equals(max, want.max) {
			o.Errorf("tree.Max().Val() = %v; want %v", max, want.max)
		}
		if min := tree.Min().Val(); !equals(min, want.min) {
			o.Errorf("tree.Min().Val() = %v; want %v", min, want.min)
		}
		if out := tree.Slice(); !sliceEquals(out, want.sorted) {
			o.Errorf("tree.Slice() = %v; want %v", out, want.sorted)
		}
		if str := tree.String(); str != want.str {
			o.Errorf("tree.String() = %v; want %v", str, want.str)
		}
		if root := tree.Root().Val(); !equals(root, want.root) {
			o.Errorf("tree.Root().Val() = %v; want %v", root, want.root)
		}
		if height := tree.Height(); height != want.height {
			o.Errorf("tree.Height() = %v; want %v", height, want.height)
		}
		for _, v := range want.sorted {
			if !tree.Contains(v) {
				o.Errorf("tree.Contains(%v) = false; want true", v)
			}
		}
		for _, v := range want.missing {
			if tree.Contains(v) {
				o.Errorf("tree.Contains(%v) = true; want false", v)
			}
		}
	}
}

func TestRandTree(o *testing.T) {
	for _, want := range tests {
		equals := func(a, b Value) bool {
			return !want.fn(a, b) && !want.fn(b, a)
		}

		tree := New(want.fn)
		tree.RandInit(want.vs)

		if siz := tree.Len(); siz != want.siz {
			o.Errorf("tree.Len() = %v; want %v", siz, want.siz)
		}
		if max := tree.Max().Val(); !equals(max, want.max) {
			o.Errorf("tree.Max().Val() = %v; want %v", max, want.max)
		}
		if min := tree.Min().Val(); !equals(min, want.min) {
			o.Errorf("tree.Min().Val() = %v; want %v", min, want.min)
		}
		if out := tree.Slice(); !sliceEquals(out, want.sorted) {
			o.Errorf("tree.Slice() = %v; want %v", out, want.sorted)
		}
		if str := tree.String(); str != want.str {
			o.Errorf("tree.String() = %v; want %v", str, want.str)
		}
		for _, v := range want.sorted {
			if !tree.Contains(v) {
				o.Errorf("tree.Contains(%v) = false; want true", v)
			}
		}
		for _, v := range want.missing {
			if tree.Contains(v) {
				o.Errorf("tree.Contains(%v) = true; want false", v)
			}
		}
	}
}

func TestNilMaxMin(o *testing.T) {
	tree := New(t1.fn)
	if min := tree.Min(); min != nil {
		o.Errorf("empty tree.Min() = %v; want nil", min)
	}
	if max := tree.Max(); max != nil {
		o.Errorf("empty tree.Max() = %v; want nil", max)
	}
}

func TestFind(o *testing.T) {
	for _, want := range tests {
		tree := New(want.fn)
		tree.RandInit(want.vs)

		for _, v := range want.sorted {
			if tree.Find(v) == nil {
				o.Errorf("tree.Find(%v) = nil; want *Node", v)
			}
			if !tree.root.Contains(v) {
				o.Errorf("tree.root.Contains(%v) = false; want true", v)
			}
		}
		for _, v := range want.missing {
			if tree.Find(v) != nil {
				o.Errorf("tree.Find(%v) != nil; want nil", v)
			}
		}
	}
}

func TestDelete(o *testing.T) {
	for _, want := range tests {
		tree := New(want.fn)
		tree.RandInit(want.vs)

		for _, v := range want.missing {
			if tree.Delete(v) {
				o.Errorf("tree.Delete(%v) = true; want false", v)
			}
		}
		for i, v := range want.remOrd {
			if !tree.Delete(v) {
				o.Errorf("tree.Delete(%v) = false; want true", v)
			}
			if tree.Delete(v) {
				o.Errorf("tree.Delete(%v) = true; want false", v)
			}
			if siz := tree.Len(); siz != want.siz-(i+1) {
				o.Errorf("tree.Len() = %v; want %v", siz, want.siz-(i+1))
			}
			if tree.Contains(v) {
				o.Errorf("tree.Contains(%v) = true; want false", v)
			}
		}
	}
}

func TestNextPrev(o *testing.T) {
	for _, want := range tests {
		equals := func(a, b Value) bool {
			return !want.fn(a, b) && !want.fn(b, a)
		}

		tree := New(want.fn)
		tree.RandInit(want.vs)

		n := tree.Min()
		for _, v := range want.sorted {
			if val := n.Val(); !equals(val, v) {
				o.Errorf("Node.Val() = %v; want %v", val, v)
			}
			n = n.Next()
		}
		if n != nil {
			o.Errorf("Node = %v; want nil", n)
		}

		n = tree.Max()
		num := len(want.sorted) - 1
		for i := range want.sorted {
			if val := n.Val(); !equals(val, want.sorted[num-i]) {
				o.Errorf("Node.Val() = %v; want %v", val, want.sorted[num-i])
			}
			n = n.Prev()
		}
		if n != nil {
			o.Errorf("Node = %v; want nil", n)
		}
	}
}

// sliceEquals returns true if two slices are equal.
func sliceEquals(a, b []Value) bool {
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
