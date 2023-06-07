package ex01

import (
	"reflect"
	"testing"
)

func TestCase1(t *testing.T) {
	root := CreateNode(true)
	root.Left = CreateNode(true)
	root.Right = CreateNode(false)
	root.Left.Left = CreateNode(true)
	root.Left.Right = CreateNode(false)
	root.Right.Left = CreateNode(true)
	root.Right.Right = CreateNode(true)
	want := []bool{true, true, false, true, true, false, true}
	got := unrollGarland(root)
	t.Run(
		"TestCase1", func(t *testing.T) {
			if !reflect.DeepEqual(got, want) {
				t.Error("error: ")
			}
		},
	)
}
