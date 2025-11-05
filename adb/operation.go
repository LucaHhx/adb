package adb

import (
	"fmt"

	"github.com/LucaHhx/adb/adb/uixml"
)

// XML 获取设备当前屏幕的 UI 层次结构（XML 格式）。
// 该方法通过 UIAutomator 导出屏幕布局，并解析为可查询的 XML 对象。
//
// 返回值：
//   - *uixml.Xml: 解析后的 XML 对象，包含完整的 UI 树结构
//   - error: 如果获取或解析失败，返回 error 对象
//
// 工作原理：
//  1. 调用 UiautomatorDump() 获取屏幕 UI 的 XML 字符串
//  2. 解析 XML 字符串为结构化对象
//  3. 返回可查询的 XML 对象
//
// 使用场景：
//   - 分析当前屏幕的 UI 结构
//   - 查找特定的 UI 元素
//   - 自动化测试中的元素定位
//   - UI 调试和分析
//
// 注意事项：
//   - 每次调用都会重新获取当前屏幕状态
//   - 获取 XML 需要一定时间，频繁调用会影响性能
//   - XML 结构会随着屏幕内容变化而变化
//   - 某些动态内容可能捕获不到
//
// 示例：
//
//	// 获取当前屏幕的 UI 结构
//	xml, err := device.XML()
//	if err != nil {
//	    log.Fatal("获取 UI 结构失败:", err)
//	}
//
//	// 查找特定按钮
//	button, err := xml.FindButton("确定")
//	if err != nil {
//	    log.Fatal("未找到按钮:", err)
//	}
//
//	// 使用自定义条件查找节点
//	node, err := xml.Find(func(n, pn uixml.Node) bool {
//	    return n.Text == "登录" && n.Clickable == "true"
//	})
func (d *Device) XML() (*uixml.Xml, error) {
	// 获取 UI 层次结构的 XML 数据
	data, err := d.UiautomatorDump()
	if err != nil {
		return nil, err
	}
	// 解析 XML 数据为结构化对象
	return uixml.NewXml(data)
}

// ClickButton 根据按钮的 content-desc 属性查找并点击按钮。
// 该方法会自动查找可点击的按钮元素，并点击其中心位置。
//
// 参数：
//   - name: 按钮的 content-desc 属性值（内容描述）
//     这是 Android 无障碍功能中用于描述 UI 元素的属性
//
// 返回值：
//   - error: 如果查找或点击失败，返回 error 对象
//
// 查找条件：
//   - 元素的 content-desc 属性必须等于指定的 name
//   - 元素的 clickable 属性必须为 "true"
//
// 使用场景：
//   - 点击有明确 content-desc 的按钮
//   - 自动化测试中的按钮点击操作
//   - 无障碍功能测试
//
// 注意事项：
//   - 如果找到多个匹配的按钮，会点击第一个
//   - 并非所有按钮都有 content-desc 属性
//   - 如果按钮使用 text 属性而非 content-desc，此方法无法找到
//   - 建议优先使用 ClickNode 方法（更灵活）
//
// 示例：
//
//	// 点击 "确定" 按钮
//	err := device.ClickButton("确定")
//	if err != nil {
//	    log.Fatal("点击按钮失败:", err)
//	}
//
//	// 点击 "提交" 按钮
//	err = device.ClickButton("提交")
//
//	// 处理按钮不存在的情况
//	err = device.ClickButton("不存在的按钮")
//	if err != nil {
//	    log.Println("按钮不存在，跳过点击")
//	}
func (d *Device) ClickButton(name string) error {
	// 查找指定名称的按钮
	button, err := d.FindButton(name)
	if err != nil {
		return err
	}
	// 点击按钮的中心位置
	return d.Tap(button.Middle())
}

// ClickNodeBy 点击指定的 UI 节点对象。
// 该方法接受一个已查找到的节点对象，并点击其中心位置。
//
// 参数：
//   - node: 要点击的 UI 节点对象（uixml.Node）
//     通常通过 FindNode 或 FindNodes 方法获取
//
// 返回值：
//   - error: 如果点击失败，返回 error 对象
//
// 使用场景：
//   - 已经通过其他方法找到节点，需要点击它
//   - 在批量操作中点击预先查找好的节点
//   - 需要对同一节点进行多次操作
//
// 注意事项：
//   - 节点对象必须包含有效的 bounds 属性
//   - 节点的位置可能随时间变化，建议及时点击
//   - 如果节点已不在屏幕上，点击会失效
//
// 示例：
//
//	// 查找并点击节点
//	node, err := device.FindNode(func(n, pn uixml.Node) bool {
//	    return n.Text == "登录"
//	})
//	if err == nil {
//	    device.ClickNodeBy(node)
//	}
//
//	// 批量点击多个节点
//	nodes, err := device.FindNodes(func(n, pn uixml.Node) bool {
//	    return n.Class == "android.widget.CheckBox"
//	})
//	if err == nil {
//	    for _, node := range nodes {
//	        device.ClickNodeBy(node)
//	        time.Sleep(500 * time.Millisecond)
//	    }
//	}
func (d *Device) ClickNodeBy(node uixml.Node) error {
	// 点击节点的中心位置
	return d.Tap(node.Middle())
}

// ClickNode 根据类名和描述/文本查找并点击 UI 节点。
// 该方法提供了灵活的节点查找方式，支持按类名和内容查找。
//
// 参数：
//   - class: UI 元素的类名（例如："android.widget.Button"）
//     如果为空字符串，则不限制类名
//   - desc: UI 元素的 content-desc 或 text 属性值
//     如果为空字符串，则只匹配类名
//     该参数会同时匹配 content-desc 和 text 属性
//
// 返回值：
//   - error: 如果查找或点击失败，返回 error 对象
//
// 查找逻辑：
//   - 必须匹配指定的 class（如果提供）
//   - content-desc 或 text 必须等于 desc（如果提供）
//   - 如果 desc 为空，则只匹配 class
//
// 常用类名：
//   - android.widget.Button: 按钮
//   - android.widget.TextView: 文本视图
//   - android.widget.EditText: 输入框
//   - android.widget.ImageView: 图片视图
//   - android.widget.CheckBox: 复选框
//   - android.view.View: 通用视图
//
// 使用场景：
//   - 按类名和文本内容定位元素
//   - 更精确的元素定位
//   - 处理动态 ID 的元素
//   - 跨应用的通用元素操作
//
// 注意事项：
//   - class 参数通常需要完整的类名（包含包路径）
//   - 如果有多个匹配的元素，会点击第一个
//   - desc 参数会同时匹配 content-desc 和 text，可能导致意外匹配
//
// 示例：
//
//	// 点击类名为 Button 且文本为 "登录" 的按钮
//	err := device.ClickNode("android.widget.Button", "登录")
//	if err != nil {
//	    log.Fatal("点击失败:", err)
//	}
//
//	// 点击任意包含 "确定" 文本的按钮（不限类名）
//	err = device.ClickNode("", "确定")
//
//	// 点击特定类名的第一个元素（不限文本）
//	err = device.ClickNode("android.widget.ImageButton", "")
//
//	// 点击 TextView 中显示 "设置" 的元素
//	err = device.ClickNode("android.widget.TextView", "设置")
func (d *Device) ClickNode(class, desc string) error {
	// 使用自定义条件查找节点
	node, err := d.FindNode(func(n, pn uixml.Node) bool {
		// 匹配类名（如果指定）
		classMatch := class == "" || n.Class == class
		// 匹配 content-desc 或 text（如果指定）
		descMatch := desc == "" || n.ContentDesc == desc || n.Text == desc
		return classMatch && descMatch
	})
	if err != nil {
		return err
	}
	// 点击找到的节点
	return d.Tap(node.Middle())
}

// FindNodeFunc 是用于查找 UI 节点的自定义条件函数类型。
// 该函数接受当前节点和父节点作为参数，返回是否匹配的布尔值。
//
// 参数说明：
//   - n: 当前正在检查的节点（uixml.Node）
//   - pn: 当前节点的父节点（uixml.Node）
//
// 返回值：
//   - bool: 如果节点满足条件返回 true，否则返回 false
//
// 使用场景：
//   - 定义复杂的节点查找条件
//   - 基于父节点的关系进行查找
//   - 组合多个条件进行精确定位
//
// 示例：
//
//	// 查找可点击的按钮
//	fn := func(n, pn uixml.Node) bool {
//	    return n.Clickable == "true" && n.Class == "android.widget.Button"
//	}
//
//	// 查找特定父节点下的文本
//	fn := func(n, pn uixml.Node) bool {
//	    return pn.ResourceID == "com.example:id/container" && n.Text != ""
//	}
type FindNodeFunc func(n, pn uixml.Node) bool

// FindNode 使用自定义条件函数查找第一个匹配的 UI 节点。
// 该方法遍历整个 UI 树，返回第一个满足条件的节点。
//
// 参数：
//   - fn: 自定义的节点匹配函数（FindNodeFunc 类型）
//     该函数会被应用到每个节点上，返回 true 表示匹配
//
// 返回值：
//   - uixml.Node: 第一个匹配的节点对象
//   - error: 如果没有找到匹配的节点或发生错误，返回 error 对象
//
// 查找顺序：
//   - 深度优先遍历 UI 树
//   - 返回第一个匹配的节点
//
// 使用场景：
//   - 需要复杂查找条件的场景
//   - 基于多个属性组合查找
//   - 需要访问父节点信息的查找
//   - 自定义的元素定位逻辑
//
// 注意事项：
//   - 每次调用都会重新获取屏幕 UI 结构
//   - 只返回第一个匹配的节点
//   - 如果没有匹配的节点，返回错误
//
// 示例：
//
//	// 查找可点击且包含特定文本的元素
//	node, err := device.FindNode(func(n, pn uixml.Node) bool {
//	    return n.Clickable == "true" && strings.Contains(n.Text, "提交")
//	})
//	if err != nil {
//	    log.Fatal("未找到节点:", err)
//	}
//
//	// 查找特定 resource-id 的元素
//	node, err = device.FindNode(func(n, pn uixml.Node) bool {
//	    return n.ResourceID == "com.example:id/login_button"
//	})
//
//	// 查找特定父节点下的可编辑元素
//	node, err = device.FindNode(func(n, pn uixml.Node) bool {
//	    return pn.ResourceID == "com.example:id/form" && n.Class == "android.widget.EditText"
//	})
//
//	// 查找启用的且非密码的输入框
//	node, err = device.FindNode(func(n, pn uixml.Node) bool {
//	    return n.Class == "android.widget.EditText" &&
//	           n.Enabled == "true" &&
//	           n.Password == "false"
//	})
func (d *Device) FindNode(fn FindNodeFunc) (uixml.Node, error) {
	// 获取当前屏幕的 UI 结构
	xml, err := d.XML()
	if err != nil {
		return uixml.Node{}, err
	}
	// 使用自定义条件查找节点
	return xml.Find(fn)
}

// FindNodes 使用自定义条件函数查找所有匹配的 UI 节点。
// 该方法遍历整个 UI 树，返回所有满足条件的节点列表。
//
// 参数：
//   - fn: 自定义的节点匹配函数（FindNodeFunc 类型）
//     该函数会被应用到每个节点上，返回 true 表示匹配
//
// 返回值：
//   - []uixml.Node: 所有匹配的节点对象数组
//   - error: 如果没有找到任何匹配的节点或发生错误，返回 error 对象
//
// 使用场景：
//   - 需要操作多个相似元素
//   - 统计特定类型的元素数量
//   - 批量验证元素属性
//   - 列表项的批量处理
//
// 注意事项：
//   - 每次调用都会重新获取屏幕 UI 结构
//   - 返回所有匹配的节点（可能很多）
//   - 如果没有匹配的节点，返回错误（不是空数组）
//   - 节点顺序与 UI 树遍历顺序相同
//
// 示例：
//
//	// 查找所有复选框
//	checkboxes, err := device.FindNodes(func(n, pn uixml.Node) bool {
//	    return n.Class == "android.widget.CheckBox"
//	})
//	if err == nil {
//	    fmt.Printf("找到 %d 个复选框\n", len(checkboxes))
//	}
//
//	// 查找所有包含特定文本的按钮
//	buttons, err := device.FindNodes(func(n, pn uixml.Node) bool {
//	    return n.Class == "android.widget.Button" &&
//	           strings.Contains(n.Text, "删除")
//	})
//
//	// 批量点击所有匹配的元素
//	nodes, err := device.FindNodes(func(n, pn uixml.Node) bool {
//	    return n.Clickable == "true" && n.Checked == "false"
//	})
//	if err == nil {
//	    for _, node := range nodes {
//	        device.ClickNodeBy(node)
//	        time.Sleep(300 * time.Millisecond)
//	    }
//	}
//
//	// 统计可见的文本元素数量
//	textViews, err := device.FindNodes(func(n, pn uixml.Node) bool {
//	    return n.Class == "android.widget.TextView" && n.Text != ""
//	})
//	if err == nil {
//	    fmt.Printf("屏幕上有 %d 个非空文本元素\n", len(textViews))
//	}
func (d *Device) FindNodes(fn FindNodeFunc) ([]uixml.Node, error) {
	// 获取当前屏幕的 UI 结构
	xml, err := d.XML()
	if err != nil {
		return nil, err
	}
	// 使用自定义条件查找所有匹配的节点
	list := xml.FindAll(fn)
	// 如果没有找到任何节点，返回错误
	if len(list) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return list, nil
}

// FindButton 根据 content-desc 查找可点击的按钮节点。
// 该方法专门用于查找带有特定 content-desc 属性的可点击按钮。
//
// 参数：
//   - name: 按钮的 content-desc 属性值（内容描述）
//
// 返回值：
//   - uixml.Node: 找到的按钮节点对象
//   - error: 如果没有找到按钮或发生错误，返回 error 对象
//
// 查找条件：
//   - 元素的 content-desc 属性必须完全匹配指定的 name
//   - 元素的 clickable 属性必须为 "true"
//
// 使用场景：
//   - 查找带有无障碍描述的按钮
//   - 需要获取按钮信息而不立即点击
//   - 验证按钮是否存在
//   - 获取按钮的位置和属性信息
//
// 注意事项：
//   - 只查找 content-desc 属性，不匹配 text 属性
//   - 必须是可点击的元素（clickable="true"）
//   - 返回第一个匹配的按钮
//   - 并非所有按钮都有 content-desc 属性
//
// 示例：
//
//	// 查找按钮
//	button, err := device.FindButton("确定")
//	if err != nil {
//	    log.Fatal("未找到按钮:", err)
//	}
//
//	// 获取按钮的位置信息
//	button, err := device.FindButton("提交")
//	if err == nil {
//	    x, y := button.Middle()
//	    fmt.Printf("按钮中心位置: (%d, %d)\n", x, y)
//	}
//
//	// 检查按钮是否存在（不点击）
//	_, err = device.FindButton("删除")
//	if err != nil {
//	    fmt.Println("删除按钮不存在")
//	} else {
//	    fmt.Println("删除按钮存在")
//	}
//
//	// 查找并获取按钮的所有属性
//	button, err := device.FindButton("设置")
//	if err == nil {
//	    fmt.Printf("按钮类名: %s\n", button.Class)
//	    fmt.Printf("按钮边界: %s\n", button.Bounds)
//	    fmt.Printf("是否启用: %s\n", button.Enabled)
//	}
func (d *Device) FindButton(name string) (uixml.Node, error) {
	// 使用自定义条件查找按钮
	return d.FindNode(func(n, pn uixml.Node) bool {
		// content-desc 必须匹配，且元素必须可点击
		return n.ContentDesc == name && n.Clickable == "true"
	})
}
