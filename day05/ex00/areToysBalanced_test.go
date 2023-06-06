package areToysBalanced

import "testing"

func TestCase1(t *testing.T) {
	root := CreateNode(false)
	root.Left = CreateNode(false)
	root.Right = CreateNode(true)
	root.Left.Left = CreateNode(false)
	root.Left.Right = CreateNode(true)
	t.Run(
		"TestCase1", func(t *testing.T) {
			if areToysBalanced(root) != true {
				t.Error("TestCase1 must be true")
			}
		},
	)
}

func TestCase2(t *testing.T) {
	root := CreateNode(true)
	root.Left = CreateNode(true)
	root.Right = CreateNode(false)
	root.Left.Left = CreateNode(true)
	root.Left.Right = CreateNode(false)
	root.Right.Left = CreateNode(true)
	root.Right.Right = CreateNode(true)
	t.Run(
		"TestCase2", func(t *testing.T) {
			if areToysBalanced(root) != true {
				t.Error("TestCase2 must be true")
			}
		},
	)
}

func TestCase3(t *testing.T) {
	root := CreateNode(true)
	root.Left = CreateNode(true)
	root.Right = CreateNode(false)
	t.Run(
		"TestCase3", func(t *testing.T) {
			if areToysBalanced(root) != false {
				t.Error("TestCase3 must be false")
			}
		},
	)
}

func TestCase4(t *testing.T) {
	root := CreateNode(false)
	root.Left = CreateNode(true)
	root.Right = CreateNode(false)
	root.Right.Left = CreateNode(true)
	root.Right.Right = CreateNode(true)
	t.Run(
		"TestCase4", func(t *testing.T) {
			if areToysBalanced(root) != false {
				t.Error("TestCase4 must be false")
			}
		},
	)
}
