package uixml

import (
	"fmt"
	"regexp"
	"strconv"
)

// boundsRe 是用于解析 Android UI 元素边界坐标的正则表达式。
// 匹配格式：[x1,y1][x2,y2]
// 其中：
//   - [x1,y1] 是元素的左上角坐标
//   - [x2,y2] 是元素的右下角坐标
//
// 示例：
//   - "[100,200][300,400]" 表示左上角 (100,200)，右下角 (300,400)
//   - 元素宽度 = x2 - x1 = 300 - 100 = 200
//   - 元素高度 = y2 - y1 = 400 - 200 = 200
var boundsRe = regexp.MustCompile(`\[(\d+),(\d+)]\[(\d+),(\d+)]`)

// Rect 表示一个矩形区域，用于描述 UI 元素的边界。
// 该结构体使用左上角和右下角的坐标来定义矩形。
//
// 字段说明：
//   - X1: 矩形左边界的 X 坐标（左上角 X）
//   - Y1: 矩形上边界的 Y 坐标（左上角 Y）
//   - X2: 矩形右边界的 X 坐标（右下角 X）
//   - Y2: 矩形下边界的 Y 坐标（右下角 Y）
//
// 坐标系统：
//   - 原点 (0,0) 位于屏幕左上角
//   - X 轴向右递增，Y 轴向下递增
//   - 所有坐标值单位为像素
//
// 计算方法：
//   - 宽度 = X2 - X1
//   - 高度 = Y2 - Y1
//   - 面积 = (X2 - X1) * (Y2 - Y1)
//   - 中心点 X = (X1 + X2) / 2
//   - 中心点 Y = (Y1 + Y2) / 2
//
// 使用场景：
//   - 存储 UI 元素的位置和大小
//   - 计算元素的中心点（用于点击）
//   - 判断元素是否重叠
//   - UI 布局分析
//
// 示例：
//
//	// 创建矩形
//	rect := uixml.Rect{X1: 100, Y1: 200, X2: 300, Y2: 400}
//
//	// 计算宽高
//	width := rect.X2 - rect.X1   // 200
//	height := rect.Y2 - rect.Y1  // 200
//
//	// 计算中心点
//	centerX := (rect.X1 + rect.X2) / 2  // 200
//	centerY := (rect.Y1 + rect.Y2) / 2  // 300
type Rect struct {
	X1 int // 左上角 X 坐标
	Y1 int // 左上角 Y 坐标
	X2 int // 右下角 X 坐标
	Y2 int // 右下角 Y 坐标
}

// ParseBounds 解析 Android UI 元素的边界坐标字符串。
// 该函数将 UIAutomator XML 中的 bounds 属性解析为 Rect 结构体。
//
// 参数：
//   - b: 边界坐标字符串，格式为 "[x1,y1][x2,y2]"
//     例如："[100,200][300,400]"
//
// 返回值：
//   - Rect: 解析后的矩形结构体，包含四个坐标值
//   - error: 如果字符串格式错误，返回错误信息
//
// 解析格式：
//   - 必须包含两对方括号
//   - 每对方括号内包含两个整数，用逗号分隔
//   - 第一对表示左上角坐标 (x1, y1)
//   - 第二对表示右下角坐标 (x2, y2)
//   - 坐标值必须是非负整数
//
// 使用场景：
//   - 解析 UIAutomator XML 的 bounds 属性
//   - 计算 UI 元素的位置和大小
//   - 获取元素的中心点坐标
//
// 错误处理：
//   - 如果格式不匹配，返回带有原始字符串的错误
//   - 支持任意大小的整数坐标值
//
// 注意事项：
//   - 格式必须严格匹配，多余的空格会导致解析失败
//   - 坐标值会被转换为 int 类型
//   - X2 应该大于 X1，Y2 应该大于 Y1（但函数不验证这一点）
//
// 示例：
//
//	// 解析标准格式的边界
//	rect, err := uixml.ParseBounds("[100,200][300,400]")
//	if err != nil {
//	    log.Fatal("解析失败:", err)
//	}
//	fmt.Printf("左上角: (%d,%d), 右下角: (%d,%d)\n",
//	    rect.X1, rect.Y1, rect.X2, rect.Y2)
//	// 输出: 左上角: (100,200), 右下角: (300,400)
//
//	// 计算元素尺寸
//	rect, _ := uixml.ParseBounds("[50,100][250,300]")
//	width := rect.X2 - rect.X1   // 200
//	height := rect.Y2 - rect.Y1  // 200
//	fmt.Printf("元素大小: %dx%d\n", width, height)
//	// 输出: 元素大小: 200x200
//
//	// 错误格式示例
//	_, err := uixml.ParseBounds("[100,200]")  // 缺少第二对坐标
//	if err != nil {
//	    fmt.Println(err)  // 输出: bad bounds format: "[100,200]"
//	}
//
//	// 从 UI 节点解析边界
//	node, _ := xml.FindButton("确定")
//	rect, err := uixml.ParseBounds(node.Bounds)
//	if err == nil {
//	    // 计算中心点用于点击
//	    centerX := (rect.X1 + rect.X2) / 2
//	    centerY := (rect.Y1 + rect.Y2) / 2
//	    device.Tap(centerX, centerY)
//	}
func ParseBounds(b string) (Rect, error) {
	// 使用正则表达式匹配边界字符串
	// 捕获组：(\d+) 匹配一个或多个数字
	// 预期匹配：[数字,数字][数字,数字]
	m := boundsRe.FindStringSubmatch(b)

	// 检查匹配结果
	// m[0] 是完整匹配，m[1-4] 是四个捕获组（x1, y1, x2, y2）
	if len(m) != 5 {
		return Rect{}, fmt.Errorf("bad bounds format: %q", b)
	}

	// 将字符串坐标转换为整数
	// strconv.Atoi 的错误被忽略，因为正则已确保是数字
	x1, _ := strconv.Atoi(m[1])
	y1, _ := strconv.Atoi(m[2])
	x2, _ := strconv.Atoi(m[3])
	y2, _ := strconv.Atoi(m[4])

	// 返回解析后的矩形结构
	return Rect{X1: x1, Y1: y1, X2: x2, Y2: y2}, nil
}
