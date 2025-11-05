package uixml

import "fmt"

func (x *Xml) FindButton(name string) (Node, error) {
	return x.Find(func(n, pn Node) bool {
		return n.ContentDesc == name && n.Clickable == "true"
	})
}
func (x *Xml) Find(fn func(n, pn Node) bool) (Node, error) {
	for _, node := range x.Nodes {
		if nodes := FindAll(node, Node{}, fn); len(nodes) > 0 && nodes[0].Class != "" {
			return nodes[0], nil
		}
	}
	return Node{}, fmt.Errorf("not found")
}

func (x *Xml) FindAll(fn func(n, pn Node) bool) []Node {
	var out []Node
	for _, node := range x.Nodes {
		out = append(out, FindAll(node, Node{}, fn)...)
	}
	return out
}

func FindAll(root, pRoot Node, predicate func(n, pn Node) bool) []Node {
	var out []Node
	Walk(root, pRoot, func(n, pn Node) {
		if predicate(n, pn) {
			out = append(out, n)
		}
	})
	return out
}
