// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binpack implements Jake Gordon's 2D binpacking algorithm.
//
// The algorithm used is described on Jake's blog here:
//
//   http://codeincomplete.com/posts/2011/5/7/bin_packing/
//
// And is also implemented by him in JavaScript here:
//
//   https://github.com/jakesgordon/bin-packing
//
package binpack

type Packable interface {
	// Len should return the number of blocks in total.
	Len() int

	// Size should return the width and height of the block n
	Size(n int) (width, height int)

	// Place should place the block n, at the position [x, y].
	Place(n, x, y int)
}

// Pack uses the packable interface, p, to pack two dimensional blocks onto a
// larger two dimensional grid.
//
// The algorithm does not start with an fixed width and height, instead it
// starts with the width and height of the first block and then grows as
// neccessary to fit each block into the overall grid. As the grid is grown the
// algorithm attempts to maintain an roughly square ratio by making 'smart'
// choices about whether to grow right or down.
//
// The returned width and height reflect how large the overall grid must be to
// contain at least each packed block.
//
// When growing, the algorithm can only grow to the right or down. If the new
// block is both wider and taller than the first [p.Size(0)] block, then the
// algorithm will be unable to pack the blocks, and [-1, -1] will be returned.
//
// To avoid the above problem, sort blocks by max(width, height).
//
// If the number of blocks is zero (p.Len() == 0) then this function is no-op
// and [0, 0] is returned.
func Pack(p Packable) (width, height int) {
	numBlocks := p.Len()

	if numBlocks == 0 {
		return 0, 0
	}

	w, h := p.Size(0)
	root := &node{
		x:      0,
		y:      0,
		width:  w,
		height: h,
	}

	p.Place(0, 0, 0)

	for i := 0; i < numBlocks; i++ {
		w, h = p.Size(i)

		node := root.find(w, h)
		if node != nil {
			node = node.split(w, h)

			// Update block in-place
			p.Place(i, node.x, node.y)

		} else {
			newRoot, grown := root.grow(w, h)
			if newRoot == nil {
				return -1, -1
			}

			// Update block in-place
			p.Place(i, grown.x, grown.y)

			root = newRoot
		}
	}

	return root.width, root.height
}

type node struct {
	x, y, width, height int
	right, down         *node
}

func (n *node) find(width, height int) *node {
	if n.right != nil || n.down != nil {
		right := n.right.find(width, height)
		if right != nil {
			return right
		}
		return n.down.find(width, height)
	} else if width <= n.width && height <= n.height {
		return n
	}
	return nil
}

func (n *node) split(width, height int) *node {
	n.down = &node{
		x:      n.x,
		y:      n.y + height,
		width:  n.width,
		height: n.height - height,
	}

	n.right = &node{
		x:      n.x + width,
		y:      n.y,
		width:  n.width - width,
		height: height,
	}

	return n
}

func (n *node) grow(width, height int) (root, grown *node) {
	canGrowDown := width <= n.width
	canGrowRight := height <= n.height

	// attempt to keep square-ish by growing right when height is much greater than width
	shouldGrowRight := canGrowRight && (n.height >= (n.width + width))

	// attempt to keep square-ish by growing down when width is much greater than height
	shouldGrowDown := canGrowDown && (n.width >= (n.height + height))

	if shouldGrowRight {
		return n.growRight(width, height)
	} else if shouldGrowDown {
		return n.growDown(width, height)
	} else if canGrowRight {
		return n.growRight(width, height)
	} else if canGrowDown {
		return n.growDown(width, height)
	}

	// need to ensure sensible root starting size to avoid this happening
	return nil, nil
}

func (n *node) growRight(width, height int) (root, grown *node) {
	newRoot := &node{
		x:      0,
		y:      0,
		width:  n.width + width,
		height: n.height,
		down:   n,
		right: &node{
			x:      n.width,
			y:      0,
			width:  width,
			height: n.height,
		},
	}

	node := newRoot.find(width, height)
	if node != nil {
		return newRoot, node.split(width, height)
	}
	return nil, nil
}

func (n *node) growDown(width, height int) (root, grown *node) {
	newRoot := &node{
		x:      0,
		y:      0,
		width:  n.width,
		height: n.height + height,
		down: &node{
			x:      0,
			y:      n.height,
			width:  n.width,
			height: height,
		},
		right: n,
	}

	node := newRoot.find(width, height)
	if node != nil {
		return newRoot, node.split(width, height)
	}
	return nil, nil
}
