package uixml

import "fmt"

// FindButton 根据 content-desc 查找可点击的按钮节点。
// 该方法是 Find 方法的便捷封装，专门用于查找按钮元素。
//
// 参数：
//   - name: 按钮的 content-desc 属性值
//
// 返回值：
//   - Node: 第一个匹配的按钮节点
//   - error: 如果没有找到按钮，返回 "not found" 错误
//
// 查找条件：
//   - 元素的 ContentDesc 必须等于 name
//   - 元素的 Clickable 属性必须为 "true"
//
// 使用场景：
//   - 快速查找按钮元素
//   - 获取按钮信息（位置、属性等）
//   - 验证按钮是否存在
//
// 注意事项：
//   - 只匹配 content-desc，不匹配 text 属性
//   - 返回第一个匹配的按钮
//   - 如果按钮使用 text 而非 content-desc，需使用 Find 方法
//
// 示例：
//
//	// 查找"确定"按钮
//	button, err := xml.FindButton("确定")
//	if err != nil {
//	    log.Fatal("未找到按钮:", err)
//	}
//	fmt.Printf("按钮位置: %s\n", button.Bounds)
//
//	// 获取按钮中心坐标
//	button, _ := xml.FindButton("提交")
//	x, y := button.Middle()
func (x *Xml) FindButton(name string) (Node, error) {
	return x.Find(func(n, pn Node) bool {
		return n.ContentDesc == name && n.Clickable == "true"
	})
}

// Find 使用自定义条件函数查找第一个匹配的 UI 节点。
// 该方法遍历整个 UI 树，返回第一个满足条件的节点。
//
// 参数：
//   - fn: 自定义的节点匹配函数
//     函数接收当前节点 n 和父节点 pn，返回 bool 表示是否匹配
//
// 返回值：
//   - Node: 第一个匹配的节点对象
//   - error: 如果没有找到匹配的节点，返回 "not found" 错误
//
// 查找过程：
//  1. 遍历根节点下的所有直接子节点
//  2. 对每个子节点递归查找匹配的节点
//  3. 返回第一个找到的非空节点（Class != ""）
//
// 使用场景：
//   - 复杂的节点查找条件
//   - 基于多个属性的组合查找
//   - 需要访问父节点信息的查找
//   - 自定义元素定位逻辑
//
// 注意事项：
//   - 深度优先遍历 UI 树
//   - 只返回第一个匹配的节点
//   - Class 为空的节点会被跳过
//
// 示例：
//
//	// 查找可点击且包含特定文本的元素
//	node, err := xml.Find(func(n, pn Node) bool {
//	    return n.Clickable == "true" && strings.Contains(n.Text, "登录")
//	})
//
//	// 查找特定 resource-id 的元素
//	node, err := xml.Find(func(n, pn Node) bool {
//	    return n.ResourceID == "com.example:id/submit_button"
//	})
//
//	// 查找特定父节点下的可编辑元素
//	node, err := xml.Find(func(n, pn Node) bool {
//	    return pn.ResourceID == "com.example:id/form" &&
//	           n.Class == "android.widget.EditText"
//	})
func (x *Xml) Find(fn func(n, pn Node) bool) (Node, error) {
	// 遍历根节点下的所有直接子节点
	for _, node := range x.Nodes {
		// 在每个子节点及其后代中查找匹配的节点
		if nodes := FindAll(node, Node{}, fn); len(nodes) > 0 && nodes[0].Class != "" {
			// 返回第一个找到的非空节点
			return nodes[0], nil
		}
	}
	// 没有找到匹配的节点
	return Node{}, fmt.Errorf("not found")
}

// FindAll 使用自定义条件函数查找所有匹配的 UI 节点。
// 该方法遍历整个 UI 树，返回所有满足条件的节点列表。
//
// 参数：
//   - fn: 自定义的节点匹配函数
//     函数接收当前节点 n 和父节点 pn，返回 bool 表示是否匹配
//
// 返回值：
//   - []Node: 所有匹配的节点对象数组（可能为空）
//
// 查找过程：
//  1. 遍历根节点下的所有直接子节点
//  2. 对每个子节点递归查找所有匹配的节点
//  3. 合并所有找到的节点并返回
//
// 使用场景：
//   - 查找所有同类型的元素
//   - 统计特定元素的数量
//   - 批量操作多个元素
//   - 列表项的查找
//
// 注意事项：
//   - 返回所有匹配的节点（可能很多）
//   - 如果没有匹配的节点，返回空数组（不是 nil）
//   - 节点顺序与 UI 树遍历顺序相同（深度优先）
//
// 示例：
//
//	// 查找所有复选框
//	checkboxes := xml.FindAll(func(n, pn Node) bool {
//	    return n.Class == "android.widget.CheckBox"
//	})
//	fmt.Printf("找到 %d 个复选框\n", len(checkboxes))
//
//	// 查找所有可点击的按钮
//	buttons := xml.FindAll(func(n, pn Node) bool {
//	    return n.Class == "android.widget.Button" && n.Clickable == "true"
//	})
//
//	// 查找所有包含文本的元素
//	textNodes := xml.FindAll(func(n, pn Node) bool {
//	    return n.Text != ""
//	})
//	for _, node := range textNodes {
//	    fmt.Println("文本:", node.Text)
//	}
func (x *Xml) FindAll(fn func(n, pn Node) bool) []Node {
	var out []Node
	// 遍历根节点下的所有直接子节点
	for _, node := range x.Nodes {
		// 在每个子节点及其后代中查找所有匹配的节点
		// 将找到的节点追加到结果数组
		out = append(out, FindAll(node, Node{}, fn)...)
	}
	return out
}

// FindAll 在指定的节点树中查找所有满足条件的节点。
// 该函数递归遍历节点树，收集所有匹配的节点。
//
// 参数：
//   - root: 要搜索的根节点
//   - pRoot: 根节点的父节点（用于传递给条件函数）
//   - predicate: 节点匹配条件函数
//     函数接收当前节点 n 和父节点 pn，返回 bool 表示是否匹配
//
// 返回值：
//   - []Node: 所有匹配的节点对象数组
//
// 工作原理：
//  1. 使用 Walk 函数递归遍历节点树（深度优先）
//  2. 对每个节点调用 predicate 函数判断是否匹配
//  3. 将所有匹配的节点收集到数组中返回
//
// 使用场景：
//   - 在特定子树中查找元素
//   - 递归收集满足条件的所有节点
//   - 内部实现，通常不直接调用
//
// 注意事项：
//   - 这是一个包级函数，不是方法
//   - 返回的数组可能为空（没有匹配的节点）
//   - 遍历顺序为深度优先
//
// 示例：
//
//	// 在特定节点下查找所有按钮
//	buttons := uixml.FindAll(containerNode, parentNode, func(n, pn Node) bool {
//	    return n.Class == "android.widget.Button"
//	})
//
//	// 查找所有启用的输入框
//	editTexts := uixml.FindAll(formNode, Node{}, func(n, pn Node) bool {
//	    return n.Class == "android.widget.EditText" && n.Enabled == "true"
//	})
func FindAll(root, pRoot Node, predicate func(n, pn Node) bool) []Node {
	var out []Node
	// 使用 Walk 递归遍历节点树
	Walk(root, pRoot, func(n, pn Node) {
		// 如果节点满足条件，添加到结果数组
		if predicate(n, pn) {
			out = append(out, n)
		}
	})
	return out
}
