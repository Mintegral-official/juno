package datastruct

import "fmt"

type TreeNode struct {
	Data  interface{}
	Right *TreeNode
	Left  *TreeNode
}

func (treeNode *TreeNode) Print() {
	if treeNode != nil {
		treeNode.Left.Print()
		treeNode.Right.Print()
		fmt.Println(treeNode.Data)
	}
}
