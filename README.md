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

1. 安装 ADB 工具（见下方安装指南）
2. 安装 Go 1.21 或更高版本
3. 连接 Android 设备并启用 USB 调试

### ADB 安装指南

#### Windows 系统

**方法一：使用 Chocolatey（推荐）**

```bash
# 安装 Chocolatey（如果未安装）
# 以管理员身份运行 PowerShell，执行以下命令
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装 ADB
choco install adb
```

**方法二：使用 Scoop**

```bash
# 安装 Scoop（如果未安装）
# 在 PowerShell 中执行
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# 安装 ADB
scoop install adb
```

**方法三：手动安装**

1. 访问 [Android SDK Platform Tools](https://developer.android.com/tools/releases/platform-tools) 下载页面
2. 下载 Windows 版本的 ZIP 文件
3. 解压到任意目录（如 `C:\platform-tools`）
4. 将解压目录添加到系统环境变量 PATH：
   - 右键点击"此电脑" → 属性 → 高级系统设置 → 环境变量
   - 在系统变量中找到 Path，点击编辑
   - 添加新路径：`C:\platform-tools`（替换为你的实际路径）
5. 重启命令提示符/PowerShell

**验证安装**

```bash
adb version
```

#### macOS 系统

**方法一：使用 Homebrew（推荐）**

```bash
# 安装 Homebrew（如果未安装）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 ADB
brew install android-platform-tools
```

**方法二：手动安装**

1. 访问 [Android SDK Platform Tools](https://developer.android.com/tools/releases/platform-tools) 下载页面
2. 下载 macOS 版本的 ZIP 文件
3. 解压到 `/usr/local/bin/` 或其他 PATH 目录
4. 添加到 PATH（如果需要）：
   ```bash
   # 编辑 ~/.zshrc 或 ~/.bash_profile
   echo 'export PATH=$PATH:/path/to/platform-tools' >> ~/.zshrc
   source ~/.zshrc
   ```

**验证安装**

```bash
adb version
```

#### Linux 系统

**Debian/Ubuntu 系列**

```bash
# 方法一：使用 apt（推荐）
sudo apt update
sudo apt install android-tools-adb android-tools-fastboot

# 方法二：添加 udev 规则（解决权限问题）
sudo apt install adb
# 添加当前用户到 plugdev 组
sudo usermod -aG plugdev $USER
```

**Fedora/RHEL/CentOS 系列**

```bash
# 使用 dnf
sudo dnf install android-tools

# 或使用 yum（旧版本）
sudo yum install android-tools
```

**Arch Linux 系列**

```bash
# 使用 pacman
sudo pacman -S android-tools

# 使用 AUR（最新版本）
yay -S android-platform-tools
```

**手动安装（适用于所有 Linux 发行版）**

1. 访问 [Android SDK Platform Tools](https://developer.android.com/tools/releases/platform-tools) 下载页面
2. 下载 Linux 版本的 ZIP 文件
3. 解压并添加到 PATH：
   ```bash
   unzip platform-tools-linux.zip
   sudo mv platform-tools /usr/local/
   echo 'export PATH=$PATH:/usr/local/platform-tools' >> ~/.bashrc
   source ~/.bashrc
   ```

**配置 udev 规则（解决权限问题）**

创建 udev 规则文件以允许非 root 用户访问 Android 设备：

```bash
# 创建规则文件
sudo nano /etc/udev/rules.d/51-android.rules

# 添加以下内容（适用于大多数 Android 设备）
SUBSYSTEM=="usb", ATTR{idVendor}=="0502", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="0b05", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="413c", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="0489", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="04c5", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="18d1", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="22b8", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="054c", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="04e8", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="0955", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="2717", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTR{idVendor}=="2a45", MODE="0666", GROUP="plugdev"

# 保存并重新加载规则
sudo chmod a+r /etc/udev/rules.d/51-android.rules
sudo udevadm control --reload-rules
sudo udevadm trigger
```

**验证安装**

```bash
adb version
```

#### 常见问题排查

**Windows**

- 如果提示 "adb 不是内部或外部命令"，请检查环境变量 PATH 配置
- 需要安装 [USB 驱动](https://developer.android.com/studio/run/oem-usb)（某些设备）

**macOS**

- 如果遇到 "无法验证开发者" 错误：
  ```bash
  xattr -d com.apple.quarantine /path/to/adb
  ```

**Linux**

- 如果遇到权限问题，确保已添加 udev 规则并重启或重新插拔设备
- 确保用户在 `plugdev` 组中：
  ```bash
  groups
  # 如果没有 plugdev，执行：
  sudo usermod -aG plugdev $USER
  # 然后注销并重新登录
  ```

#### 验证 ADB 连接

安装完成后，连接 Android 设备并执行：

```bash
# 列出已连接的设备
adb devices

# 如果看到设备序列号和 "device" 状态，说明连接成功
# 示例输出：
# List of devices attached
# ABC123DEF456    device
```

**首次连接设备时**：
1. 确保设备已启用 **USB 调试**（开发者选项中）
2. 连接设备后，会在手机上弹出授权提示
3. 点击"允许 USB 调试"并勾选"始终允许"
4. 再次运行 `adb devices` 确认连接成功

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
