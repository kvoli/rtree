package rtree

import "fmt"

// OTreeNode Represents a treenode for use with org tree (no frac cascase)
type OTreeNode struct {
	XKey        int
	Left, Right *OTreeNode
	YTree       *YTreeNode
}

// CTreeNode represents a fractional cascading Tree Node
type CTreeNode struct {
	XKey        int
	Left, Right *CTreeNode
	YArray      []*CYTreeNode
}

// YTreeNode represents a treenode of the internal Y subtre within the OTree
type YTreeNode struct {
	Pnt         *Point
	Left, Right *YTreeNode
}

// CYTreeNode represents a treenode within an array for casc
type CYTreeNode struct {
	Pnt         *Point
	Left, Right int
}

func newOTreeNode(p *Point) *OTreeNode {
	return &OTreeNode{
		XKey:  p.X,
		Left:  nil,
		Right: nil,
		YTree: nil,
	}
}

func newYTreeNode(p *Point) *YTreeNode {
	return &YTreeNode{
		Pnt:   p,
		Left:  nil,
		Right: nil,
	}
}

func newCTreeNode(p *Point) *CTreeNode {
	return &CTreeNode{
		XKey:   p.X,
		Left:   nil,
		Right:  nil,
		YArray: nil,
	}
}

func newCYTreeNode(p *Point) *CYTreeNode {
	return &CYTreeNode{
		Pnt:   p,
		Left:  -1,
		Right: -1,
	}
}

func mergeYTree(l, r *YTreeNode) *YTreeNode {
	lArray := storeInOrder(l)
	rArray := storeInOrder(r)
	res := merge(lArray, rArray)
	return arrayToTree(res, 0, len(res)-1)
}

func arrayToTree(array []*Point, start, end int) *YTreeNode {
	if start > end {
		return nil
	}
	m := (start + end) / 2
	tree := newYTreeNode(array[m])
	tree.Left = arrayToTree(array, start, m-1)
	tree.Right = arrayToTree(array, m+1, end)
	return tree
}

func merge(left, right []*Point) []*Point {
	size, i, j := len(left)+len(right), 0, 0
	slice := make([]*Point, size, size)

	for k := 0; k < size; k++ {
		if i > len(left)-1 && j <= len(right)-1 {
			slice[k] = right[j]
			j++
		} else if j > len(right)-1 && i <= len(left)-1 {
			slice[k] = left[i]
			i++
		} else if left[i].Y < right[j].Y {
			slice[k] = left[i]
			i++
		} else {
			slice[k] = right[j]
			j++
		}
	}
	return slice
}

func insert(root *YTreeNode, x *Point) *YTreeNode {
	if root == nil {
		return newYTreeNode(x)
	}
	if x.Y > root.Pnt.Y {
		root.Right = insert(root.Right, x)
	} else {
		root.Left = insert(root.Left, x)
	}
	return root
}

func naiveMerge(l, r *YTreeNode) *YTreeNode {
	lArray := storeInOrder(l)
	rArray := storeInOrder(r)
	var res *YTreeNode
	var ins []*Point
	if len(lArray) > len(rArray) {
		res = arrayToTree(lArray, 0, len(lArray)-1)
		ins = rArray
	} else {
		res = arrayToTree(rArray, 0, len(rArray)-1)
		ins = lArray
	}
	for _, v := range ins {
		res = insert(res, v)
	}
	return res
}

func storeInOrder(root *YTreeNode) []*Point {
	if root != nil {
		L := storeInOrder(root.Left)
		R := storeInOrder(root.Right)
		res := append(L, root.Pnt)
		res = append(res, R...)
		return res
	}
	return nil
}

func mapPoints(A []*CYTreeNode) []*Point {
	res := make([]*Point, len(A), len(A))
	for i, v := range A {
		res[i] = v.Pnt
	}
	return res
}

func mapYNodes(A []*Point) []*CYTreeNode {
	res := make([]*CYTreeNode, len(A), len(A))
	for i, v := range A {
		res[i] = newCYTreeNode(v)
	}
	return res
}

func linkPointers(left bool, cur, ref []*CYTreeNode) {
	i := 0
	for _, v := range cur {
		for i < len(ref) && v.Pnt.Y > ref[i].Pnt.Y {
			i++
		}
		if i < len(ref) {
			if left {
				v.Left = i
			} else {
				v.Right = i
			}
		}
	}
}

func mergeFracCascade(l, r []*CYTreeNode) []*CYTreeNode {
	// needs to alloc a pointer on each new YTreeNode to the succ YTreeNode
	lPoints := mapPoints(l)
	rPoints := mapPoints(r)
	combPoints := merge(lPoints, rPoints)
	res := mapYNodes(combPoints)
	linkPointers(true, res, l)
	linkPointers(false, res, r)
	return res
}

// ContrSorted constructs a 2D Range tree in O(nlogn) time by using merge sort principles for O(n) Y Tree
func ContrSorted(points []*Point) *OTreeNode {
	n := len(points)
	if n == 1 {
		oNode := newOTreeNode(points[0])
		oNode.YTree = newYTreeNode(points[0])
		return oNode
	} else if n > 1 {
		m := (n - 1) / 2
		oNode := newOTreeNode(points[m])
		oNode.Left = ContrSorted(points[0 : m+1])
		oNode.Right = ContrSorted(points[m+1 : n])
		var l, r *YTreeNode
		if oNode.Left != nil {
			l = oNode.Left.YTree
		}
		if oNode.Right != nil {
			r = oNode.Right.YTree
		}
		oNode.YTree = mergeYTree(l, r)
		return oNode
	}
	return nil
}

// ContrNaive constructs a 2D Range Tree in O(nlog^2n) time by naively creating Y Tree in O(nlogn)
func ContrNaive(points []*Point) *OTreeNode {
	n := len(points)
	if n == 1 {
		oNode := newOTreeNode(points[0])
		oNode.YTree = newYTreeNode(points[0])
		return oNode
	} else if n > 1 {
		m := (n - 1) / 2
		oNode := newOTreeNode(points[m])
		oNode.Left = ContrSorted(points[0 : m+1])
		oNode.Right = ContrSorted(points[m+1 : n])
		var l, r *YTreeNode
		if oNode.Left != nil {
			l = oNode.Left.YTree
		}
		if oNode.Right != nil {
			r = oNode.Right.YTree
		}
		oNode.YTree = naiveMerge(l, r)
		return oNode
	}
	return nil
}

// ContrFC constructs a 2D Range Tree using fractional cascading in O(nlogn) time
func ContrFC(points []*Point) *CTreeNode {
	n := len(points)
	if n == 1 {
		cNode := newCTreeNode(points[0])
		yArray := make([]*CYTreeNode, 1, 1)
		yArray[0] = newCYTreeNode(points[0])
		cNode.YArray = yArray
		return cNode
	} else if n > 1 {
		m := (n - 1) / 2
		cNode := newCTreeNode(points[m])
		cNode.Left = ContrFC(points[0 : m+1])
		cNode.Right = ContrFC(points[m+1 : n])
		var l, r []*CYTreeNode
		if cNode.Left != nil {
			l = cNode.Left.YArray
		}
		if cNode.Right != nil {
			r = cNode.Right.YArray
		}
		cNode.YArray = mergeFracCascade(l, r)
		return cNode
	}
	return nil
}

// Traverse prints out vals of OTree in sorted order
func Traverse(root *OTreeNode) {
	if root == nil {
		return
	}

	Traverse(root.Left)
	printYTree(root)
	Traverse(root.Right)
}

func printYTree(root *OTreeNode) {
	fmt.Printf("\nX KEY: %d, Count: %d, MaxDepthX %d, MaxDepthY %d @ ", root.XKey, nodeCount(root.YTree), xDepth(root), yDepth(root.YTree))
	printYTreeHelper(root.YTree)
}

func printXTree(root *OTreeNode) {
	if root != nil {
		printXTree(root.Left)
	}
}

func printSpacing(level int) {
	for v := 0; v < level; v++ {
		fmt.Printf(" ")
	}
}

// PrintLevelOrder prints level order traaversal
func PrintLevelOrder(root *OTreeNode) {
	h := xDepth(root)
	for v := 0; v < h; v++ {
		printSpacing(h - v)
		printGivenLevel(root, v)
		printSpacing(h - v)
		fmt.Printf("\n")
	}
}

func printGivenLevel(root *OTreeNode, level int) {
	if root != nil {
		if level == 1 {
			fmt.Printf("%v ", root.XKey)
		}
		printGivenLevel(root.Left, level-1)
		printGivenLevel(root.Right, level-1)
	}
}

func printYTreeHelper(root *YTreeNode) {
	if root == nil {
		return
	}
	printYTreeHelper(root.Left)
	fmt.Printf("(X: %d, Y: %d),  ", root.Pnt.X, root.Pnt.Y)
	printYTreeHelper(root.Right)
}

func nodeCount(root *YTreeNode) int {
	if root == nil {
		return 0
	}

	return nodeCount(root.Left) + nodeCount(root.Right) + 1
}

func yDepth(root *YTreeNode) int {
	if root == nil {
		return 0
	}

	return max(yDepth(root.Left), yDepth(root.Right)) + 1
}

func xDepth(root *OTreeNode) int {
	if root == nil {
		return 0
	}

	return max(xDepth(root.Left), xDepth(root.Right)) + 1
}

// TraverseC prints out vals of OTree in sorted order
func TraverseC(root *CTreeNode) {
	if root == nil {
		return
	}

	TraverseC(root.Left)
	printCYTree(root)
	TraverseC(root.Right)
}

func printCYTree(root *CTreeNode) {
	fmt.Printf("\nX KEY: %d, Count: %d, MaxDepthX %d @ [", root.XKey, len(root.YArray), cDepth(root))

	for _, v := range root.YArray {
		fmt.Printf("{pnt: (%d,%d): (", v.Pnt.X, v.Pnt.Y)
		if v.Left != -1 {
			fmt.Printf("%v,", root.Left.YArray[v.Left].Pnt.Y)
		} else {
			fmt.Printf("-,")
		}
		if v.Right != -1 {
			fmt.Printf("%v", root.Right.YArray[v.Right].Pnt.Y)
		} else {
			fmt.Printf("-")
		}
		fmt.Printf(")} ")
	}
	fmt.Printf("]\n")
}

func cDepth(root *CTreeNode) int {
	if root == nil {
		return 0
	}

	return max(cDepth(root.Left), cDepth(root.Right)) + 1
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func findSplitNode(x1, x2 int, root *OTreeNode) *OTreeNode {
	if root != nil {
		if isLeaf(root) {
			return root
		}
		if root.XKey >= x1 && root.XKey <= x2 {
			return root
		}
		if root.XKey < x1 {
			return findSplitNode(x1, x2, root.Right)
		}
		if root.XKey > x2 {
			return findSplitNode(x1, x2, root.Left)
		}
	}
	return nil
}

// QueryOTree returns the list of points in range [p1.x,p2.x] x [p1.y, p2.y]
func QueryOTree(p1, p2 *Point, root *OTreeNode) []*Point {
	vsplit := findSplitNode(p1.X, p2.X, root)
	// vsplit := root
	// fmt.Printf("\n[%v,%v]x[%v,%v]\n", p1.X, p2.X, p1.Y, p2.Y)
	// fmt.Printf("SPLIT NODE %+v", vsplit)
	if isLeaf(vsplit) && vsplit.XKey >= p1.X && vsplit.XKey <= p2.X {
		return reportYTree(root.YTree, p1.Y, p2.Y)
	}
	res := handlePath(vsplit.Left, p1, p2, true)
	res = append(res, handlePath(vsplit.Right, p1, p2, false)...)
	return res
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func isLeaf(root *OTreeNode) bool {
	return root.Left == nil && root.Right == nil
}

func handlePath(root *OTreeNode, p1, p2 *Point, isLeft bool) []*Point {
	if root == nil {
		return make([]*Point, 0)
	}
	if isLeaf(root) && root.XKey >= p1.X && root.XKey <= p2.X {
		return reportYTree(root.YTree, p1.Y, p2.Y)
	}

	if isLeft {
		if root.XKey < p1.X {
			return handlePath(root.Right, p1, p2, true)
		}
		return append(reportYTree(root.Right.YTree, p1.Y, p2.Y), handlePath(root.Left, p1, p2, true)...)
	}

	if root.XKey > p2.X {
		return handlePath(root.Left, p1, p2, false)
	}
	return append(reportYTree(root.Left.YTree, p1.Y, p2.Y), handlePath(root.Right, p1, p2, false)...)
}

func reportYTree(root *YTreeNode, a, b int) []*Point {
	if root != nil {
		res := make([]*Point, 0)
		if a < root.Pnt.Y {
			res = append(res, reportYTree(root.Left, a, b)...)
		}
		if a <= root.Pnt.Y && b >= root.Pnt.Y {
			res = append(res, root.Pnt)
		}
		if b > root.Pnt.Y {
			res = append(res, reportYTree(root.Right, a, b)...)
		}
		return res
	}
	return make([]*Point, 0)
}

func findCSplitNode(x1, x2 int, root *CTreeNode) *CTreeNode {
	if root != nil {
		if isCLeaf(root) {
			return root
		}
		if root.XKey >= x1 && root.XKey <= x2 {
			return root
		}
		if root.XKey < x1 {
			return findCSplitNode(x1, x2, root.Right)
		}
		if root.XKey > x2 {
			return findCSplitNode(x1, x2, root.Left)
		}
	}
	return nil
}

func findYNode(A []*CYTreeNode, x int) int {
	if A != nil {
		l, r := 0, len(A)-1
		m := l + (r-l)/2
		for r >= l {
			m = l + (r-l)/2
			if A[m].Pnt.Y == x {
				return m
			}
			if A[m].Pnt.Y > x {
				r = m - 1
			} else {
				l = m + 1
			}
		}
		return l
	}
	return -1
}

func reportCYTree(A []*CYTreeNode, a, b, i int) []*Point {
	if A == nil || i > len(A)-1 || i < 0 {
		return make([]*Point, 0)
	}
	res := make([]*Point, 0)
	for i < len(A) && inRange(A[i].Pnt.Y, a, b) {
		res = append(res, A[i].Pnt)
		i++
	}
	return res
}

func inRange(x, a, b int) bool {
	return x >= a && x <= b
}

func isCLeaf(root *CTreeNode) bool {
	return root.Left == nil && root.Right == nil
}

func handleLCPath(root *CTreeNode, p1, p2 *Point, artic int) []*Point {
	if root == nil || artic < 0 {
		return make([]*Point, 0)
	}
	if isCLeaf(root) && inRange(root.XKey, p1.X, p2.X) {
		return reportCYTree(root.YArray, p1.Y, p2.Y, artic)
	}
	if root.XKey < p1.X {
		return handleLCPath(root.Right, p1, p2, root.YArray[artic].Right)
	}
	return append(reportCYTree(root.Right.YArray, p1.Y, p2.Y, root.YArray[artic].Right), handleLCPath(root.Left, p1, p2, root.YArray[artic].Left)...)
}

func handleRCPath(root *CTreeNode, p1, p2 *Point, artic int) []*Point {
	if root == nil || artic < 0 {
		return make([]*Point, 0)
	}
	if isCLeaf(root) && inRange(root.XKey, p1.X, p2.X) {
		return reportCYTree(root.YArray, p1.Y, p2.Y, artic)
	}
	if root.XKey > p2.X {
		return handleRCPath(root.Left, p1, p2, root.YArray[artic].Left)
	}
	return append(reportCYTree(root.Left.YArray, p1.Y, p2.Y, root.YArray[artic].Left), handleRCPath(root.Right, p1, p2, root.YArray[artic].Right)...)
}

// QueryCTree queries with [p1.x,p2.x] x [p1.y,p2.y]
func QueryCTree(p1, p2 *Point, root *CTreeNode) []*Point {
	vSplit := findCSplitNode(p1.X, p2.X, root)
	articNode := findYNode(vSplit.YArray, p1.Y)
	if isCLeaf(vSplit) && inRange(vSplit.XKey, p1.X, p2.X) {
		reportCYTree(vSplit.YArray, p1.Y, p2.Y, articNode)
	}
	res := handleLCPath(vSplit.Left, p1, p2, vSplit.YArray[articNode].Left)
	return append(res, handleRCPath(vSplit.Right, p1, p2, vSplit.YArray[articNode].Right)...)

}
