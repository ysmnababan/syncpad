package main

type VectorNode []int

func NewVectorNode(l int) *VectorNode {
	if l == 0 {
		panic("length can't be")
	}
	var x VectorNode = make([]int, l)
	return &x
}

// AddLocalEvent
//
//	`i` must be an slice index, starts from 0.
func (n *VectorNode) IncrementNode(idx int) *VectorNode {
	out := NewVectorNode(len(*n))
	copy(*out, *n)
	if idx >= 0 && idx < len(*n) {
		val := (*n)[idx]
		val++
		(*out)[idx] = val
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (n *VectorNode) MaxVector(v *VectorNode) *VectorNode {
	out := NewVectorNode(len(*v))
	for i, val := range *n {
		maxVal := max((*v)[i], val)
		// (*n)[i] = maxVal
		(*out)[i] = maxVal
	}
	return out
}
