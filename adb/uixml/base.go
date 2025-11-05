// Package uixml 提供对 Android UIAutomator dump XML 的解析与查询工具。
// 该包实现了对 UI 层次结构 XML 的完整解析和查询功能，支持灵活的元素定位。
package uixml

import (
	"encoding/xml"
	"io"
	"strings"
)

// Xml 是 UIAutomator XML 数据的封装结构。
// 它包含了解析后的 UI 层次结构，并提供便捷的查询方法。
//
// 字段说明：
//   - *Hierarchy: 嵌入的层次结构对象，包含所有 UI 节点
//
// 使用方式：
//   - 通过 NewXml() 创建实例
//   - 使用 Find() 或 FindAll() 查找节点
type Xml struct {
	*Hierarchy
}

// NewXml 从 XML 字符串创建 Xml 对象。
// 该函数解析 UIAutomator 导出的 XML 数据，构建可查询的对象结构。
//
// 参数：
//   - data: UIAutomator dump 命令输出的 XML 字符串
//
// 返回值：
//   - *Xml: 解析后的 Xml 对象，可用于查询 UI 元素
//   - error: 如果 XML 解析失败，返回 error 对象
//
// 使用场景：
//   - 将设备返回的 UI XML 数据转换为可查询的对象
//   - UI 自动化测试的数据准备
//   - UI 结构分析
//
// 示例：
//
//	// 获取 UI XML 数据并解析
//	xmlData, err := device.UiautomatorDump()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	xml, err := uixml.NewXml(xmlData)
//	if err != nil {
//	    log.Fatal("解析失败:", err)
//	}
//
//	// 查找元素
//	button, _ := xml.FindButton("确定")
func NewXml(data string) (*Xml, error) {
	// 解析 XML 字符串为 Hierarchy 结构
	xmlData, err := ParseHierarchyFromString(data)
	if err != nil {
		return nil, err
	}
	// 返回包装后的 Xml 对象
	return &Xml{xmlData}, nil
}

// Walk 递归遍历 UI 节点树，对每个节点执行指定的函数。
// 该函数实现深度优先遍历，先处理当前节点，再递归处理子节点。
//
// 参数：
//   - n: 当前要处理的节点
//   - pn: 当前节点的父节点（用于在处理函数中访问父节点信息）
//   - fn: 对每个节点执行的函数，接收当前节点和父节点作为参数
//
// 遍历顺序：
//   - 深度优先（先序遍历）
//   - 先处理当前节点，再递归处理所有子节点
//
// 使用场景：
//   - 遍历所有 UI 元素
//   - 统计特定类型的元素数量
//   - 收集满足条件的所有节点
//   - UI 树的自定义分析
//
// 注意事项：
//   - 这是一个递归函数，复杂 UI 可能导致较深的递归
//   - fn 函数会被调用多次（每个节点一次）
//   - 修改节点不会影响原始 XML 数据
//
// 示例：
//
//	// 打印所有节点的文本内容
//	uixml.Walk(rootNode, uixml.Node{}, func(n, pn uixml.Node) {
//	    if n.Text != "" {
//	        fmt.Println("文本:", n.Text)
//	    }
//	})
//
//	// 统计可点击元素数量
//	count := 0
//	uixml.Walk(rootNode, uixml.Node{}, func(n, pn uixml.Node) {
//	    if n.Clickable == "true" {
//	        count++
//	    }
//	})
//	fmt.Printf("可点击元素数量: %d\n", count)
func Walk(n, pn Node, fn func(n, pn Node)) {
	// 先处理当前节点
	fn(n, pn)
	// 递归处理所有子节点
	for _, c := range n.Children {
		Walk(c, n, fn)
	}
}

// Hierarchy 表示 UI 层次结构的根节点。
// 对应 UIAutomator XML 中的 <hierarchy> 标签。
//
// 字段说明：
//   - XMLName: XML 标签名（固定为 "hierarchy"）
//   - Rotation: 屏幕旋转角度（0, 90, 180, 270）
//   - Nodes: 根节点下的所有直接子节点数组
//
// XML 示例：
//
//	<hierarchy rotation="0">
//	  <node index="0" text="..." class="..." ...>
//	    ...
//	  </node>
//	</hierarchy>
type Hierarchy struct {
	XMLName  xml.Name `xml:"hierarchy"`
	Rotation string   `xml:"rotation,attr"`
	Nodes    []Node   `xml:"node"`
}

// Node 表示 UI 层次结构中的一个节点（UI 元素）。
// 对应 UIAutomator XML 中的 <node> 标签，包含所有常用的 UI 属性。
//
// 字段说明（Android UI 属性）：
//   - NAF: Not Accessibility Friendly，是否对无障碍不友好
//   - Index: 元素在父节点中的索引位置
//   - Text: 元素的文本内容
//   - ResourceID: 元素的资源 ID（如 "com.example:id/button"）
//   - Class: 元素的类名（如 "android.widget.Button"）
//   - Package: 元素所属的应用包名
//   - ContentDesc: 内容描述，用于无障碍功能
//   - Checkable: 元素是否可选中（"true" 或 "false"）
//   - Checked: 元素是否已选中
//   - Clickable: 元素是否可点击
//   - Enabled: 元素是否启用
//   - Focusable: 元素是否可获得焦点
//   - Focused: 元素是否已获得焦点
//   - Scrollable: 元素是否可滚动
//   - LongClickable: 元素是否可长按
//   - Password: 元素是否是密码输入框
//   - Selected: 元素是否被选中
//   - Bounds: 元素的边界坐标，格式为 "[x1,y1][x2,y2]"
//   - Children: 该节点的所有子节点数组
//
// XML 示例：
//
//	<node index="0" text="登录" resource-id="com.example:id/login_btn"
//	      class="android.widget.Button" package="com.example"
//	      content-desc="登录按钮" clickable="true" enabled="true"
//	      bounds="[100,200][300,400]" />
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

// Middle 计算并返回节点边界的中心点坐标。
// 该方法解析节点的 Bounds 属性，计算矩形区域的中心位置。
//
// 返回值：
//   - x: 中心点的 X 坐标（像素）
//   - y: 中心点的 Y 坐标（像素）
//   - 如果 Bounds 解析失败，返回 (0, 0)
//
// 计算方式：
//   - X 坐标 = (右边界 - 左边界) / 2 + 左边界
//   - Y 坐标 = (下边界 - 上边界) / 2 + 上边界
//
// 使用场景：
//   - 点击元素时需要元素的中心坐标
//   - 计算元素之间的距离
//   - UI 布局分析
//
// 注意事项：
//   - 如果 Bounds 属性为空或格式错误，返回 (0, 0)
//   - 返回的坐标是屏幕绝对坐标
//
// 示例：
//
//	// 获取按钮中心坐标
//	button, _ := xml.FindButton("确定")
//	x, y := button.Middle()
//	fmt.Printf("按钮中心: (%d, %d)\n", x, y)
//
//	// 点击元素中心
//	x, y := node.Middle()
//	device.Tap(x, y)
func (n *Node) Middle() (x, y int) {
	// 解析 Bounds 字符串为矩形结构
	bounds, err := ParseBounds(n.Bounds)
	if err != nil {
		// 解析失败返回原点
		return 0, 0
	}
	// 计算并返回中心点坐标
	// X = (宽度 / 2) + 左边界
	// Y = (高度 / 2) + 上边界
	return (bounds.X2-bounds.X1)/2 + bounds.X1, (bounds.Y2-bounds.Y1)/2 + bounds.Y1
}

// ---------- 解析入口 ----------

// ParseHierarchy 从 io.Reader 解析 UI 层次结构 XML。
// 该函数读取 XML 数据流并解析为 Hierarchy 结构。
//
// 参数：
//   - r: 包含 XML 数据的 io.Reader
//
// 返回值：
//   - *Hierarchy: 解析后的层次结构对象
//   - error: 如果解析失败，返回 error 对象
//
// 使用场景：
//   - 从文件读取 XML
//   - 从网络流读取 XML
//   - 从任意数据源读取 XML
//
// 示例：
//
//	// 从文件解析
//	file, _ := os.Open("ui_dump.xml")
//	defer file.Close()
//	hierarchy, err := uixml.ParseHierarchy(file)
//
//	// 从字符串解析（通常使用 ParseHierarchyFromString）
//	reader := strings.NewReader(xmlString)
//	hierarchy, err := uixml.ParseHierarchy(reader)
func ParseHierarchy(r io.Reader) (*Hierarchy, error) {
	var h Hierarchy
	// 创建 XML 解码器
	dec := xml.NewDecoder(r)
	// 解码 XML 数据到 Hierarchy 结构
	if err := dec.Decode(&h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ParseHierarchyFromString 从字符串解析 UI 层次结构 XML。
// 该函数是 ParseHierarchy 的便捷封装，直接接受字符串参数。
//
// 参数：
//   - s: 包含 XML 数据的字符串
//
// 返回值：
//   - *Hierarchy: 解析后的层次结构对象
//   - error: 如果解析失败，返回 error 对象
//
// 使用场景：
//   - 解析 UIAutomator dump 命令的输出
//   - 处理从设备获取的 XML 字符串
//   - 测试和调试
//
// 示例：
//
//	// 解析设备返回的 XML
//	xmlData, _ := device.UiautomatorDump()
//	hierarchy, err := uixml.ParseHierarchyFromString(xmlData)
//	if err != nil {
//	    log.Fatal("解析失败:", err)
//	}
//
//	// 访问根节点
//	for _, node := range hierarchy.Nodes {
//	    fmt.Println("根节点类名:", node.Class)
//	}
func ParseHierarchyFromString(s string) (*Hierarchy, error) {
	// 将字符串转换为 Reader 并调用 ParseHierarchy
	return ParseHierarchy(strings.NewReader(s))
}
