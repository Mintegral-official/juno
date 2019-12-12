package main

import "fmt"

//
//import (
//	"bufio"
//	"fmt"
//	"io"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//)
//
//func main() {
//
//	fi, err := os.Open("/Users/tangye/go/src/juno/count/main/message.txt")
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		return
//	}
//	defer fi.Close()
//
//	var bst intBinarySearchTree
//	br := bufio.NewReader(fi)
//
//	t := time.Now()
//	for {
//		a, _, c := br.ReadLine()
//		if c == io.EOF {
//			break
//		}
//		a1 := strings.Split(string(a), ":")
//		res, _ := strconv.Atoi(strings.Trim(a1[len(a1)-1], " "))
//		bst.Insert(res, res)
//	}
//	fmt.Println(time.Since(t))
//
//	t = time.Now()
//	var result []int
//	bst.InOrderTraverse(func(i int) {
//		result = append(result, i)
//	})
//
//	fmt.Println(time.Since(t))
//
//	fmt.Println(len(result))
//	fmt.Println(result[len(result)/2])
//	fmt.Println(result[len(result)*3/4])
//	fmt.Println(result[len(result)*9/10])
//	fmt.Println(result[len(result)*95/100])
//	fmt.Println(result[len(result)*99/100])
//
//}
//
//// 节点
//type Node struct {
//	Key   int
//	Value int
//	left  *Node
//	right *Node
//}
//
//type intBinarySearchTree struct {
//	Root *Node
//}
//
//func (bst *intBinarySearchTree) Insert(key int, value int) {
//	n := &Node{key, value, nil, nil}
//
//	if bst.Root == nil {
//		bst.Root = n
//	} else {
//		insertNode(bst.Root, n)
//	}
//}
//
//func insertNode(node, newNode *Node) {
//	if newNode.Key < node.Key {
//		if node.left == nil {
//			node.left = newNode
//		} else {
//			insertNode(node.left, newNode)
//		}
//	} else {
//		if node.right == nil {
//			node.right = newNode
//		} else {
//			insertNode(node.right, newNode)
//		}
//	}
//}
//
//func (bst *intBinarySearchTree) InOrderTraverse(f func(int int)) {
//	inOrderTraver(bst.Root, f)
//}
//
//func inOrderTraver(n *Node, f func(int int)) {
//	if n != nil {
//		inOrderTraver(n.left, f)
//		f(n.Value)
//		inOrderTraver(n.right, f)
//	}
//}
type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseBetween(head *ListNode, m int, n int) *ListNode {
	var prev, tmp *ListNode
	e, cur, count := head, head, 0
	for count < n {
		count++
		if count > n {
			break
		}
		if count == m {
			e = cur
		}
		if m <= count {
			cur, cur.Next, prev = cur.Next, prev, cur
			continue
		}
		tmp = cur
		cur = cur.Next
	}
	if tmp != nil {
		tmp.Next, prev = prev, head
	}

	if e != nil {
		e.Next = cur
	}
	return prev
}

func main() {
	l := &ListNode{
		Val: 0,
		Next: &ListNode{
			Val:  1,
			Next: nil,
		},
	}
	fmt.Println(reverseBetween(l, 1, 2))
}

