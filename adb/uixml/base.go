// Package uixml 提供对 Android UIAutomator dump XML 的解析与查询工具。
package uixml

import (
	"encoding/xml"
	"io"
	"strings"
)

type Xml struct {
	*Hierarchy
}

func NewXml(data string) (*Xml, error) {
	xmlDta, err := ParseHierarchyFromString(data)
	if err != nil {
		return nil, err
	}
	return &Xml{xmlDta}, nil
}

func Walk(n, pn Node, fn func(n, pn Node)) {
	fn(n, pn)
	for _, c := range n.Children {
		Walk(c, n, fn)
	}
}

// Hierarchy 对应 <hierarchy ...> 根标签
type Hierarchy struct {
	XMLName  xml.Name `xml:"hierarchy"`
	Rotation string   `xml:"rotation,attr"`
	Nodes    []Node   `xml:"node"`
}

// Node 对应 <node ...>，列出常用属性
type Node struct {
	NAF           string `xml:"NAF,attr"`
	Index         string `xml:"index,attr"`
	Text          string `xml:"text,attr"`
	ResourceID    string `xml:"resource-id,attr"`
	Class         string `xml:"class,attr"`
	Package       string `xml:"package,attr"`
	ContentDesc   string `xml:"content-desc,attr"`
	Checkable     string `xml:"checkable,attr"`
	Checked       string `xml:"checked,attr"`
	Clickable     string `xml:"clickable,attr"`
	Enabled       string `xml:"enabled,attr"`
	Focusable     string `xml:"focusable,attr"`
	Focused       string `xml:"focused,attr"`
	Scrollable    string `xml:"scrollable,attr"`
	LongClickable string `xml:"long-clickable,attr"`
	Password      string `xml:"password,attr"`
	Selected      string `xml:"selected,attr"`
	Bounds        string `xml:"bounds,attr"`

	Children []Node `xml:"node"`
}

func (n *Node) Middle() (x, y int) {
	bounds, err := ParseBounds(n.Bounds)
	if err != nil {
		return 0, 0
	}
	return (bounds.X2-bounds.X1)/2 + bounds.X1, (bounds.Y2-bounds.Y1)/2 + bounds.Y1
}

// ---------- 解析入口 ----------

// ParseHierarchy 从 io.Reader 解析
func ParseHierarchy(r io.Reader) (*Hierarchy, error) {
	var h Hierarchy
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ParseHierarchyFromString 从字符串解析
func ParseHierarchyFromString(s string) (*Hierarchy, error) {
	return ParseHierarchy(strings.NewReader(s))
}
