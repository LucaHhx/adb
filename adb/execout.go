package adb

import (
	"fmt"
	"regexp"
	"strings"
)

// Execout 执行 adb exec-out 命令并返回输出。
// exec-out 命令类似于 shell 命令，但输出是二进制安全的，不会进行转换。
//
// 参数：
//   - command: 要在设备上执行的命令字符串
//
// 返回值：
//   - string: 命令的原始输出
//   - error: 如果执行失败，返回 error 对象
//
// 使用场景：
//   - 需要二进制安全输出的场景
//   - 执行 uiautomator 等需要原始输出的命令
//   - 获取截图等二进制数据
//
// 注意事项：
//   - 输出不经过 shell 处理，保持原始格式
//   - 适合处理二进制或特殊格式的数据
//
// 示例：
//
//	// 执行 uiautomator dump
//	output, err := device.Execout("uiautomator dump /dev/tty")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (d *Device) Execout(command string) (string, error) {
	return d.execCommand("exec-out", command)
}

// errorMap 存储常见错误信息的翻译映射。
// 用于将设备返回的非中文错误信息转换为中文，便于理解。
var errorMap = map[string]string{
	"Mohon periksa koneksi internet Anda.": "请检查您的互联网连接。",
}

// UiautomatorDump 获取设备当前屏幕的 UI 层次结构 XML 数据。
// 该方法使用 UIAutomator 工具导出屏幕布局信息。
//
// 返回值：
//   - string: 包含完整 UI 层次结构的 XML 字符串
//   - error: 如果获取失败或包含已知错误信息，返回 error 对象
//
// 工作原理：
//  1. 执行 'uiautomator dump /dev/tty' 命令
//  2. UIAutomator 分析当前屏幕的 UI 结构
//  3. 将结构导出为 XML 格式并输出到 /dev/tty（标准输出）
//  4. 检查输出中是否包含已知的错误信息
//  5. 返回 XML 字符串或错误
//
// XML 结构说明：
//   - 根元素为 <hierarchy>
//   - 每个 UI 元素表示为 <node> 标签
//   - 包含丰富的属性信息（class, text, bounds, clickable 等）
//
// 使用场景：
//   - UI 自动化测试的元素定位基础
//   - 屏幕内容分析
//   - UI 调试和分析
//   - 元素查找的数据源
//
// 注意事项：
//   - 每次调用都会重新分析屏幕，有一定性能开销
//   - 某些动态内容可能捕获不完整
//   - 返回的 XML 可能很大（复杂界面）
//   - 会检测并翻译已知的错误信息
//
// 示例：
//
//	// 获取屏幕 UI 结构
//	xmlData, err := device.UiautomatorDump()
//	if err != nil {
//	    log.Fatal("获取 UI 结构失败:", err)
//	}
//	fmt.Println(xmlData)
//
//	// 保存到文件
//	xmlData, err := device.UiautomatorDump()
//	if err == nil {
//	    os.WriteFile("screen_layout.xml", []byte(xmlData), 0644)
//	}
func (d *Device) UiautomatorDump() (string, error) {
	// 执行 uiautomator dump 命令，输出到 /dev/tty（标准输出）
	command, err := d.execCommand("exec-out", "uiautomator dump /dev/tty")
	if err != nil {
		return "", err
	}

	// 检查输出中是否包含已知的错误信息
	for originalErr, translatedErr := range errorMap {
		if strings.Contains(command, originalErr) {
			return "", fmt.Errorf(translatedErr)
		}
	}

	return command, nil
}

// ExistElement 检查屏幕上是否存在包含指定内容的 UI 元素。
// 该方法通过获取屏幕 UI 结构并进行字符串匹配来判断元素是否存在。
//
// 参数：
//   - content: 要查找的内容字符串
//     可以是元素的 text、content-desc、resource-id 等任何属性值
//     也可以是 XML 中的任意文本片段
//
// 返回值：
//   - bool: 如果找到包含该内容的元素返回 true，否则返回 false
//     注意：获取 UI 结构失败时也会返回 false
//
// 匹配方式：
//   - 简单的字符串包含匹配（不区分属性）
//   - 只要 XML 中包含该字符串即可
//   - 大小写敏感
//
// 使用场景：
//   - 快速检查特定文本是否显示在屏幕上
//   - 验证元素是否存在
//   - 条件判断（如等待某个元素出现）
//   - 简单的内容验证
//
// 注意事项：
//   - 这是一个简单的字符串匹配，不够精确
//   - 可能会匹配到 XML 标签名或其他非预期内容
//   - 建议使用更精确的 FindNode 方法进行元素查找
//   - 获取 UI 结构失败时返回 false（无法区分元素不存在还是获取失败）
//
// 示例：
//
//	// 检查"登录"按钮是否存在
//	if device.ExistElement("登录") {
//	    fmt.Println("登录按钮存在")
//	} else {
//	    fmt.Println("登录按钮不存在")
//	}
//
//	// 等待元素出现
//	for i := 0; i < 10; i++ {
//	    if device.ExistElement("加载完成") {
//	        fmt.Println("页面加载完成")
//	        break
//	    }
//	    time.Sleep(1 * time.Second)
//	}
//
//	// 检查特定 resource-id 是否存在
//	if device.ExistElement("com.example:id/submit_button") {
//	    fmt.Println("提交按钮存在")
//	}
func (d *Device) ExistElement(content string) bool {
	// 获取屏幕 UI 结构
	command, err := d.UiautomatorDump()
	if err != nil {
		return false
	}
	// 检查 XML 中是否包含指定内容
	return strings.Contains(command, content)
}

// Regexp 使用正则表达式从屏幕 UI 结构中提取匹配的内容。
// 该方法获取屏幕 UI 的 XML 结构，并应用正则表达式提取第一个捕获组的内容。
//
// 参数：
//   - rex: 正则表达式字符串，必须包含至少一个捕获组 ()
//     捕获组用于指定要提取的内容部分
//
// 返回值：
//   - string: 正则表达式第一个捕获组匹配的内容
//   - error: 如果获取 UI 结构失败或没有匹配项，返回 error 对象
//
// 使用场景：
//   - 提取特定格式的文本内容
//   - 从复杂的 UI 结构中提取特定信息
//   - 需要灵活匹配的场景
//   - 提取动态内容
//
// 注意事项：
//   - 正则表达式必须包含捕获组 ()
//   - 只返回第一个捕获组的内容
//   - 如果有多个匹配项，只返回第一个
//   - 正则表达式语法错误会导致失败
//
// 示例：
//
//	// 提取 content-desc 属性值
//	// XML: <node content-desc="用户名: admin" ... />
//	username, err := device.Regexp(`content-desc="用户名: ([^"]+)"`)
//	if err == nil {
//	    fmt.Println("用户名:", username) // 输出: admin
//	}
//
//	// 提取 resource-id
//	// XML: <node resource-id="com.example:id/status_text" text="在线" ... />
//	status, err := device.Regexp(`resource-id="com.example:id/status_text"[^>]*text="([^"]+)"`)
//	if err == nil {
//	    fmt.Println("状态:", status) // 输出: 在线
//	}
//
//	// 提取数字
//	// XML: <node text="余额: 1234.56元" ... />
//	balance, err := device.Regexp(`text="余额: ([\d.]+)元"`)
//	if err == nil {
//	    fmt.Println("余额:", balance) // 输出: 1234.56
//	}
func (d *Device) Regexp(rex string) (string, error) {
	// 获取屏幕 UI 结构
	command, err := d.UiautomatorDump()
	if err != nil {
		return "", err
	}

	// 编译正则表达式
	re := regexp.MustCompile(rex)

	// 查找匹配项（包括捕获组）
	matches := re.FindStringSubmatch(command)

	// 检查是否有捕获组匹配
	if len(matches) > 1 {
		// 返回第一个捕获组的内容
		return matches[1], nil
	}

	// 没有找到匹配项
	return "", fmt.Errorf("not found")
}

// FindDesc 根据元素的边界坐标（bounds）查找其 content-desc 属性值。
// 该方法先定位指定边界的节点，然后提取其 content-desc 属性。
//
// 参数：
//   - bounds: 元素的边界坐标字符串，格式为 "[x1,y1][x2,y2]"
//     例如："[100,200][300,400]" 表示左上角 (100,200)，右下角 (300,400)
//
// 返回值：
//   - string: 该节点的 content-desc 属性值
//   - error: 如果查找失败，返回 error 对象
//
// 工作流程：
//  1. 获取屏幕 UI 结构的 XML
//  2. 使用正则表达式查找具有指定 bounds 的 <node> 标签
//  3. 从找到的节点中提取 content-desc 属性值
//  4. 返回属性值或错误
//
// 使用场景：
//   - 已知元素位置，需要获取其描述信息
//   - 验证特定位置的元素内容
//   - 通过坐标反查元素信息
//   - UI 测试中的断言验证
//
// 注意事项：
//   - bounds 参数必须精确匹配（包括格式）
//   - 如果节点没有 content-desc 属性，会打印提示并返回空字符串
//   - 坐标必须完全相同，差一个像素都不会匹配
//   - 如果找不到节点，会打印提示并返回空字符串
//
// 示例：
//
//	// 获取指定位置元素的描述
//	desc, err := device.FindDesc("[100,200][300,400]")
//	if err != nil {
//	    log.Fatal("查找失败:", err)
//	}
//	if desc != "" {
//	    fmt.Println("元素描述:", desc)
//	} else {
//	    fmt.Println("该元素没有 content-desc 属性")
//	}
//
//	// 验证特定位置的按钮文本
//	desc, err := device.FindDesc("[500,1000][700,1100]")
//	if err == nil && desc == "确定" {
//	    fmt.Println("确定按钮在预期位置")
//	}
func (d *Device) FindDesc(bounds string) (string, error) {
	// 获取屏幕 UI 结构
	data, err := d.UiautomatorDump()
	if err != nil {
		return "", err
	}

	// 构建正则表达式：查找具有指定 bounds 的 <node> 标签
	// regexp.QuoteMeta 用于转义 bounds 中的特殊字符
	nodeRe := regexp.MustCompile(`<node\b[^>]*\bbounds="` + regexp.QuoteMeta(bounds) + `"[^>]*/>`)
	nodeMatch := nodeRe.FindString(data)

	// 检查是否找到目标节点
	if nodeMatch == "" {
		fmt.Println("未找到目标节点")
		return "", nil
	}

	// 从节点标签中提取 content-desc 属性
	contentRe := regexp.MustCompile(`content-desc="([^"]*)"`)
	content := contentRe.FindStringSubmatch(nodeMatch)

	// 检查是否找到 content-desc 属性
	if len(content) > 1 {
		return content[1], nil
	}

	// 节点存在但没有 content-desc 属性
	fmt.Println("未找到 content-desc 属性")
	return "", nil
}
