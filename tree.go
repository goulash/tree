// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package tree provides a binary search tree implementation.
//
// There are a few more functions in the coming, when I have time.
package tree

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Value interface{}

// Node represents the internal nodes of a binary search tree.
//
// If a node is not nil, then it must store a value, and may contain
// links to one or two subtrees (left and right). A node always has
// a pointer to the parent node, unless it is the root node.
type Node struct {
	val Value

	parent *Node
	left   *Node
	right  *Node

	// lessFn is in each Node to allow us to do things with Nodes without
	// knowing what Tree they are in.
	lessFn func(a, b Value) bool
}

// Val returns the value stored by the node.
func (n *Node) Val() Value {
	return n.val
}

// Next returns the next node after n, or nil if n is the last node.
func (n *Node) Next() *Node {
	if n.right != nil {
		return n.right.Min()
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	return p
}

// Prev returns the previous node before n, or nil if n is the first node.
func (n *Node) Prev() *Node {
	if n.left != nil {
		return n.left.Max()
	}
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}
	return p
}

// String provides a string representation of the elements in the (sub)tree,
// in ascending order.
func (t *Node) String() string {
	var buf bytes.Buffer

	// walk prints out all the elements in the (sub)tree.
	var walk func(n *Node)
	walk = func(n *Node) {
		if n != nil {
			walk(n.left)
			fmt.Fprintf(&buf, "%v ", n.val)
			walk(n.right)
		}
	}

	fmt.Fprint(&buf, "[")
	if t != nil {
		walk(t)

		// remove the ", " at the end
		buf.Truncate(buf.Len() - 1)
	}
	fmt.Fprintf(&buf, "]")
	return buf.String()
}

// Find searches the (sub)tree for the value v and returns the node if it is
// found.
func (t *Node) Find(v Value) *Node {
	for t != nil {
		if t.lessFn(v, t.val) {
			t = t.left
		} else if t.lessFn(t.val, v) {
			t = t.right
		} else {
			return t
		}
	}
	return nil
}

// Contains searches the (sub)tree for the value v and returns true if it is
// found.
func (t *Node) Contains(v Value) bool {
	return t.Find(v) != nil
}

// Max returns the maximum value found in the (sub)tree,
// or nil if the (sub)tree is empty.
func (t *Node) Max() *Node {
	if t == nil {
		return nil
	}
	for t.right != nil {
		t = t.right
	}
	return t
}

// Min returns the minimum value found in the (sub)tree,
// or nil if the (sub)tree is empty.
func (t *Node) Min() *Node {
	if t == nil {
		return nil
	}
	for t.left != nil {
		t = t.left
	}
	return t
}

// Height calculates the maximum height of the (sub)tree.
//
// Note: a tree with 0 elements has a height of 0; a tree with 1 element
// must have a height of 1; a tree with 2 elements a height of 2; and a
// tree with n elements must have a height >= lg n.
func (t *Node) Height() int {
	if t == nil {
		return 0
	}

	l, r := t.left.Height(), t.right.Height()
	return 1 + maxInt(l, r)
}

// maxInt returns the greater of two numbers a and b.
func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Tree represents a binary search tree.
//
// The zero value of a Tree is a ready to use tree. Do note however, that
// you will need a pointer to the tree to use any of the methods.
type Tree struct {
	root *Node
	size int

	lessFn func(a, b Value) bool
}

// New returns a new Tree to use, with lessFn as the function that gives a < b.
//
// A Tree must be created with this function, otherwise trying to insert into
// it will cause a panic.
func New(lessFn func(a, b Value) bool) *Tree {
	return &Tree{lessFn: lessFn}
}

// Init initializes the tree with elements from vs in the given order.
//
// In the general case, RandInit should be used, as Init cannot hope to have
// a logarithmic height if vs is sorted in any way.
func (t *Tree) Init(vs []Value) {
	for _, v := range vs {
		t.Insert(v)
	}
}

// RandInit initializes the tree with elements from vs in random order.
//
// If you know that vs is already random, then you should use Init instead,
// as it will be slightly more efficient (about twice as fast).
//
// Note: you are responsible for seeding math/rand.
func (t *Tree) RandInit(vs []Value) {
	n := len(vs)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	for i := n - 1; i >= 0; i-- {
		r := rand.Intn(i + 1)
		t.Insert(vs[idx[r]])
		idx[r], idx[i] = idx[i], idx[r]
	}
}

// Len returns the size of the tree.
func (t *Tree) Len() int {
	return t.size
}

// Root returns the root node of the tree, which is nil if the tree is empty.
func (t *Tree) Root() *Node {
	return t.root
}

func (t *Tree) String() string { return t.root.String() }

// Slice returns the tree as a slice.
//
// The slice is completely disconnected from the tree, you can do whatever you
// want with it.
func (t *Tree) Slice() []Value {
	array := make([]Value, t.Len())

	var slice func(t *Node, i int) int
	slice = func(n *Node, i int) int {
		if n != nil {
			if n.left != nil {
				i = slice(n.left, i)
			}
			array[i] = n.val
			i++
			if n.right != nil {
				i = slice(n.right, i)
			}
		}
		return i
	}

	slice(t.root, 0)
	return array
}

func (t *Tree) Find(v Value) *Node { return t.root.Find(v) }

func (t *Tree) Contains(v Value) bool { return t.root.Find(v) != nil }

func (t *Tree) Max() *Node { return t.root.Max() }

func (t *Tree) Min() *Node { return t.root.Min() }

func (t *Tree) Height() int { return t.root.Height() }

// Range returns the search range [from, to] as a slice.
//func (t *Tree) Range(from, to Value) []Value {
//	return nil
//	// TODO
//}

// Insert inserts a value v into the tree if it does not exist and returns the
// node containing it.
//
// Note: if the value v already is in the tree, nothing happens.
// If you really want to know if it succeeded, check Len() before and after.
func (t *Tree) Insert(v Value) *Node {
	var n *Node
	x := t.root
	for x != nil {
		n = x
		if t.lessFn(v, x.val) {
			x = x.left
		} else if t.lessFn(x.val, v) {
			x = x.right
		} else {
			return x
		}
	}

	z := &Node{v, n, nil, nil, t.lessFn}
	if n == nil {
		// then the tree was empty => n = nil
		t.root = z
	} else if t.lessFn(v, n.val) {
		n.left = z
	} else { // n.val < v
		n.right = z
	}
	t.size++
	return z
}

// Delete removes the value v from the tree, returning true if successful.
func (t *Tree) Delete(v Value) bool {
	if n := t.Find(v); n != nil {
		t.removeNode(n)
		return true
	}
	return false
}

// removeNode removes a node from the tree.
// Note: we assume that n != nil!
func (t *Tree) removeNode(n *Node) {
	if n.left == nil {
		t.transplant(n, n.right)
	} else if n.right == nil {
		t.transplant(n, n.left)
	} else {
		// successor of n
		s := n.right
		for s.left != nil {
			s = s.left
		}
		if s.parent != n {
			t.transplant(s, s.right)
			s.right = n.right
			s.right.parent = s
		}
		t.transplant(n, s)
		s.left = n.left
		s.left.parent = s
	}
	t.size--
}

// transplant replaces n with m in the tree.
// Note: we assume that n != nil!
func (t *Tree) transplant(u, v *Node) {
	if u.parent == nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}
