package datastruct

import "fmt"

type TreeNode struct {
	Data  interface{}
	Right *TreeNode
	Left  *TreeNode
}

var res []interface{}

func (treeNode *TreeNode) Print() {
	if treeNode != nil {
		treeNode.Left.Print()
		treeNode.Right.Print()
		fmt.Printf("****%v****\n", treeNode.Data)
	}
}

func (treeNode *TreeNode) to() []interface{} {
	if treeNode != nil {
		treeNode.Left.to()
		treeNode.Right.to()
		res = append(res, treeNode.Data)
	}
	return res
}

func (treeNode *TreeNode) To() Stack {
	s, r := NewStack(), treeNode.to()
	for i := len(r) - 1; i >= 0; i-- {
		s.Push(r[i])
	}
	return *s
}
