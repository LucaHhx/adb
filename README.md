# ADB 自动化工具库

这是一个用 Go 语言编写的 Android 设备自动化工具库，提供了对 ADB（Android Debug Bridge）的高级封装，让 Android 自动化测试和操作变得更加简单。

## 功能特性

### 核心功能

- **设备管理**：连接和管理 Android 设备
- **UI 自动化**：通过 UIAutomator 解析和操作 UI 元素
- **触摸操作**：点击、滑动等手势模拟
- **输入控制**：文本输入、按键事件
- **应用控制**：启动/停止应用、Activity 管理
- **文件传输**：设备与本地之间的文件传输
- **剪贴板操作**：读取设备剪贴板内容

### UI 元素查找

- 支持按 class、text、content-desc 等属性查找元素
- 提供灵活的自定义查找函数
- 自动计算元素中心点坐标
- 支持查找可点击的按钮元素

## 项目结构

```
.
├── adb/                    # 核心库代码
│   ├── adb.go             # 设备管理
│   ├── shell.go           # Shell 命令封装
│   ├── operation.go       # UI 操作封装
│   ├── utils.go           # 工具函数
│   └── uixml/             # UI XML 解析
│       ├── base.go        # 基础结构定义
│       ├── find.go        # 元素查找
│       └── utils.go       # XML 工具函数
├── example/               # 示例代码
│   ├── main.go           # 基础示例
│   └── find_node/        # 元素查找示例
├── init/                  # 初始化脚本和工具
│   ├── init.sh           # 设备初始化脚本
│   ├── geto.apk          # 权限管理工具
│   ├── adbkey.apk        # ADB 键盘
│   └── clipper.apk       # 剪贴板工具
└── go.mod                # Go 模块定义
```

## 安装

### 在其他项目中使用

在你的 Go 项目中安装这个库：

```bash
go get github.com/LucaHhx/adb
```

### 创建新项目使用示例

1. 创建一个新的 Go 项目：

```bash
mkdir my-adb-project
cd my-adb-project
go mod init my-adb-project
```

2. 安装此库：

```bash
go get github.com/LucaHhx/adb
```

3. 创建 main.go 文件：

```go
package main

import (
    "fmt"
    "github.com/LucaHhx/adb/adb"
)

func main() {
    // 创建设备实例
    dev := adb.NewDevice()

    // 点击屏幕
    err := dev.Tap(500, 1000)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

4. 运行：

```bash
go run main.go
```

## 快速开始

### 环境准备

1. 安装 ADB 工具
2. 安装 Go 1.21 或更高版本
3. 连接 Android 设备并启用 USB 调试

### 设备初始化

首次使用前，需要初始化 Android 设备：

```bash
cd init
./init.sh
```

初始化脚本会自动完成以下配置：

- 安装 Geto（权限管理工具）
- 安装 ADB Keyboard（允许通过 ADB 输入文本）
- 安装 Clipper（剪贴板工具）
- 启用触摸显示和指针位置显示
- 配置必要的权限
- 设置 ADB Keyboard 为默认输入法

### 基础使用

```go
package main

import (
    "fmt"
    "github.com/LucaHhx/adb/adb"
)

func main() {
    // 创建设备实例
    dev := adb.NewDevice()

    // 点击屏幕坐标
    dev.Tap(500, 1000)

    // 输入文本
    dev.Input("Hello World")

    // 启动应用
    dev.StartActivity("com.example.app", "com.example.app.MainActivity")

    // 查找并点击按钮
    dev.ClickButton("登录")

    // 获取剪贴板内容
    clipboardText, _ := dev.GetClipper()
    fmt.Println(clipboardText)
}
```

### UI 元素查找示例

```go
// 通过自定义函数查找元素
node, err := dev.FindNode(func(n, pn uixml.Node) bool {
    return n.Class == "android.widget.Button" && n.Text == "提交"
})

// 点击找到的元素
dev.ClickNodeBy(node)

// 查找所有匹配的元素
nodes, err := dev.FindNodes(func(n, pn uixml.Node) bool {
    return n.Clickable == "true"
})
```

## API 文档

### 设备操作

- `NewDevice(serial ...string)` - 创建设备实例
- `Shell(command string)` - 执行 Shell 命令
- `Connect(address string)` - 连接到网络设备

### 触摸和输入

- `Tap(x, y int)` - 点击指定坐标
- `Swipe(x1, y1, x2, y2, duration int32)` - 滑动操作
- `Input(text string)` - 输入文本
- `KeyEvent(keyCode int)` - 发送按键事件
- `PressBack()` - 按返回键
- `PressHome()` - 按主屏幕键
- `PressEnter()` - 按回车键

### 应用管理

- `StartActivity(packageName, activityName string)` - 启动 Activity
- `ForceStopApp(packageName string)` - 强制停止应用

### UI 操作

- `XML()` - 获取当前屏幕的 UI XML 结构
- `FindNode(fn FindNodeFunc)` - 查找单个元素
- `FindNodes(fn FindNodeFunc)` - 查找多个元素
- `ClickButton(name string)` - 点击按钮
- `ClickNode(class, desc string)` - 点击指定元素

### 文件操作

- `Pull(devicePath, localPath string)` - 从设备拉取文件
- `Push(localPath, devicePath string)` - 推送文件到设备

### 工具功能

- `GetClipper()` - 获取剪贴板内容
- `UiautomatorDump()` - 导出 UI 层级结构

## 依赖项

本项目仅依赖 Go 标准库，无需额外的第三方依赖。

## 使用场景

- Android 应用自动化测试
- UI 功能测试
- 重复性操作自动化
- 设备批量操作
- 应用压力测试

## 注意事项

1. 使用前确保设备已启用 USB 调试模式
2. 某些操作需要设备具有 ROOT 权限
3. 建议使用 ADB Keyboard 进行文本输入，支持中文和特殊字符
4. 使用 `init.sh` 脚本配置设备环境以获得最佳体验

## License

MIT
