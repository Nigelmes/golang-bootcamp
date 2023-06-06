package areToysBalanced

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func CreateNode(v bool) *TreeNode {
	return &TreeNode{v, nil, nil}
}

func GetTrueToyNum(root *TreeNode) int {
	if root == nil {
		return 0
	} else {
		if root.HasToy {
			return GetTrueToyNum(root.Left) + GetTrueToyNum(root.Right) + 1
		}
		return GetTrueToyNum(root.Left) + GetTrueToyNum(root.Right)
	}
}

func areToysBalanced(root *TreeNode) bool {
	if root == nil {
		return false
	}
	left := GetTrueToyNum(root.Left)
	right := GetTrueToyNum(root.Right)
	return left == right
}
