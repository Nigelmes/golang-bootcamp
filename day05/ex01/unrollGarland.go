package ex01

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

func unrollGarland(root *TreeNode) []bool {
	if root == nil {
		return []bool{}
	}
	queue := []*TreeNode{root}
	result := []bool{root.HasToy}
	level := 0
	for len(queue) > 0 {
		levelSize := len(queue)

		for i, j := 0, levelSize; i < levelSize; i, j = i+1, j-1 {
			var node *TreeNode
			if level%2 != 0 {
				node = queue[j-1]
			} else {
				node = queue[i]
			}

			if level%2 == 0 {
				if node.Left != nil {
					queue = append(queue, node.Left)
					result = append(result, node.Left.HasToy)
				}
				if node.Right != nil {
					queue = append(queue, node.Right)
					result = append(result, node.Right.HasToy)
				}

			} else {
				if node.Right != nil {
					queue = append(queue, node.Right)
					result = append(result, node.Right.HasToy)
				}
				if node.Left != nil {
					queue = append(queue, node.Left)
					result = append(result, node.Left.HasToy)
				}
			}
		}
		level++
		queue = queue[levelSize:]
	}

	return result
}
