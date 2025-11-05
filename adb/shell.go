package adb

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Shell 在 Android 设备上执行 shell 命令并返回输出结果。
// 该方法允许直接在设备的 Linux shell 环境中执行任意命令。
//
// 参数：
//   - command: 要在设备上执行的 shell 命令字符串
//     支持所有标准的 Linux shell 命令，如 ls, cat, ps, getprop 等
//     可以使用管道、重定向等 shell 特性
//
// 返回值：
//   - string: 命令的标准输出（已去除首尾空白字符）
//   - error: 如果命令执行失败，返回 error 对象
//
// 使用场景：
//   - 获取设备属性信息（如系统版本、设备型号）
//   - 执行系统命令（如列出文件、查看进程）
//   - 修改系统设置
//   - 执行自定义脚本
//
// 注意事项：
//   - 某些命令可能需要 root 权限才能执行
//   - 命令中的特殊字符（如引号、反斜杠）需要正确转义
//   - 长时间运行的命令可能会导致超时
//   - 返回值会自动去除首尾空白字符
//
// 示例：
//
//	// 获取 Android 版本
//	version, err := device.Shell("getprop ro.build.version.release")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Android 版本:", version)
//
//	// 列出目录内容
//	output, err := device.Shell("ls -la /sdcard/")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
//
//	// 查看当前运行的进程
//	processes, err := device.Shell("ps | grep com.example")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("相关进程:", processes)
//
//	// 获取设备型号
//	model, err := device.Shell("getprop ro.product.model")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("设备型号:", model)
func (d *Device) Shell(command string) (string, error) {
	return d.execCommand("shell", command)
}

// Tap 在设备屏幕的指定坐标处模拟点击操作。
// 该方法通过 'input tap' 命令实现屏幕点击，可用于自动化测试和UI交互。
//
// 参数：
//   - x: 点击位置的 X 坐标（像素值，从屏幕左边开始）
//   - y: 点击位置的 Y 坐标（像素值，从屏幕顶部开始）
//
// 返回值：
//   - error: 如果点击操作失败，返回 error 对象
//
// 坐标系说明：
//   - 原点 (0,0) 位于屏幕左上角
//   - X 轴从左到右递增
//   - Y 轴从上到下递增
//   - 坐标值取决于设备的屏幕分辨率
//
// 使用场景：
//   - 自动化测试中模拟用户点击
//   - UI 自动化操作
//   - 批量操作应用界面
//   - 游戏脚本自动化
//
// 注意事项：
//   - 坐标必须在设备屏幕范围内，否则点击无效
//   - 某些应用可能会检测并阻止模拟点击
//   - 如果屏幕被锁定，点击可能无效
//   - 建议在点击后添加适当的延时，等待UI响应
//
// 示例：
//
//	// 点击屏幕中心（假设屏幕分辨率为 1080x1920）
//	err := device.Tap(540, 960)
//	if err != nil {
//	    log.Fatal("点击失败:", err)
//	}
//
//	// 点击指定按钮位置
//	buttonX, buttonY := 200, 500
//	err = device.Tap(buttonX, buttonY)
//	if err != nil {
//	    log.Fatal("点击按钮失败:", err)
//	}
//
//	// 连续点击多个位置
//	positions := [][2]int{{100, 200}, {300, 400}, {500, 600}}
//	for _, pos := range positions {
//	    device.Tap(pos[0], pos[1])
//	    time.Sleep(500 * time.Millisecond) // 每次点击后等待
//	}
func (d *Device) Tap(x, y int) error {
	// 构建 input tap 命令
	command := fmt.Sprintf("input tap %d %d", x, y)
	_, err := d.Shell(command)
	return err
}

// Swipe 在设备屏幕上执行滑动操作，从起始坐标滑动到结束坐标。
// 该方法通过 'input swipe' 命令实现，可用于模拟滚动、翻页等手势操作。
//
// 参数：
//   - x1: 滑动起始位置的 X 坐标（像素值）
//   - y1: 滑动起始位置的 Y 坐标（像素值）
//   - x2: 滑动结束位置的 X 坐标（像素值）
//   - y2: 滑动结束位置的 Y 坐标（像素值）
//   - duration: 滑动持续时间（毫秒），控制滑动速度
//
// 返回值：
//   - error: 如果滑动操作失败，返回 error 对象
//
// 滑动方向说明：
//   - 向上滑动: y1 > y2（从下往上）
//   - 向下滑动: y1 < y2（从上往下）
//   - 向左滑动: x1 > x2（从右往左）
//   - 向右滑动: x1 < x2（从左往右）
//
// 持续时间说明：
//   - duration 越小，滑动速度越快（快速滑动/快速翻页）
//   - duration 越大，滑动速度越慢（慢速滑动/精确控制）
//   - 典型值：300-1000 毫秒为正常滑动速度
//
// 使用场景：
//   - 页面滚动查看更多内容
//   - 左右滑动切换页面
//   - 下拉刷新操作
//   - 游戏中的滑动操作
//   - 列表浏览
//
// 注意事项：
//   - 坐标必须在屏幕范围内
//   - duration 为 0 时表示瞬间完成滑动
//   - 滑动距离太短可能被识别为点击
//   - 某些应用可能检测并限制模拟滑动
//
// 示例：
//
//	// 向上滑动（滚动查看下方内容）
//	// 从屏幕底部向上滑动到顶部，持续 500 毫秒
//	err := device.Swipe(500, 1500, 500, 500, 500)
//	if err != nil {
//	    log.Fatal("向上滑动失败:", err)
//	}
//
//	// 向下滑动（滚动查看上方内容或下拉刷新）
//	err = device.Swipe(500, 500, 500, 1500, 500)
//
//	// 向左滑动（切换到下一页）
//	err = device.Swipe(800, 900, 200, 900, 300)
//
//	// 向右滑动（切换到上一页）
//	err = device.Swipe(200, 900, 800, 900, 300)
//
//	// 快速滑动（快速翻页）
//	err = device.Swipe(500, 1500, 500, 300, 200)
//
//	// 慢速滑动（精确控制）
//	err = device.Swipe(500, 1500, 500, 500, 1000)
func (d *Device) Swipe(x1, y1, x2, y2, duration int32) error {
	// 构建 input swipe 命令
	command := fmt.Sprintf("input swipe %d %d %d %d %d", x1, y1, x2, y2, duration)
	_, err := d.Shell(command)
	return err
}

// Input 向当前焦点的输入框发送文本内容。
// 该方法通过广播机制实现文本输入，支持包含空格和特殊字符的文本。
//
// 参数：
//   - text: 要输入的文本字符串
//     支持中文、英文、数字和特殊字符
//     空格会自动转换为 %s 以符合 ADB 命令格式
//
// 返回值：
//   - error: 如果输入操作失败，返回 error 对象
//
// 工作原理：
//   - 使用 Android 的广播机制 (am broadcast) 发送文本
//   - 需要设备上安装支持 ADB_INPUT_TEXT 广播的应用或服务
//   - 文本会被发送到当前具有输入焦点的输入框
//
// 使用场景：
//   - 自动填写表单
//   - 登录账号密码
//   - 搜索框输入
//   - 聊天消息发送
//   - 自动化测试中的文本输入
//
// 注意事项：
//   - 需要确保目标输入框已获得焦点
//   - 某些设备或 ROM 可能不支持此广播方式
//   - 如果不支持广播方式，建议使用 'input text' 命令（但不支持中文）
//   - 空格会被自动替换为 %s
//   - 特殊字符可能需要额外转义
//
// 示例：
//
//	// 输入简单文本
//	err := device.Input("Hello World")
//	if err != nil {
//	    log.Fatal("输入失败:", err)
//	}
//
//	// 输入包含空格的文本
//	err = device.Input("This is a test message")
//
//	// 输入密码
//	err = device.Input("MyPassword@123")
//
//	// 自动登录示例
//	// 1. 点击用户名输入框
//	device.Tap(500, 600)
//	time.Sleep(500 * time.Millisecond)
//	// 2. 输入用户名
//	device.Input("myusername")
//	// 3. 点击密码输入框
//	device.Tap(500, 800)
//	time.Sleep(500 * time.Millisecond)
//	// 4. 输入密码
//	device.Input("mypassword")
//	// 5. 点击登录按钮
//	device.Tap(500, 1000)
func (d *Device) Input(text string) error {
	// 将空格替换为 %s 以适配 ADB input 命令格式
	escapedText := strings.ReplaceAll(text, " ", "%s")
	// 构建广播命令发送文本
	// am broadcast: 发送广播
	// -a: 指定 action（ADB_INPUT_TEXT）
	// --es: 附加字符串数据（msg 为 key，escapedText 为 value）
	command := fmt.Sprintf("am broadcast -a ADB_INPUT_TEXT --es msg '%s'", escapedText)
	_, err := d.Shell(command)
	return err
}

// KeyEvent 向设备发送指定的按键事件。
// 该方法通过 'input keyevent' 命令实现，可以模拟各种硬件按键和软键盘按键。
//
// 参数：
//   - keyCode: Android 系统定义的按键代码（整数）
//     常用按键代码见 Android KeyEvent 类文档
//
// 返回值：
//   - error: 如果发送按键失败，返回 error 对象
//
// 常用按键代码（KeyEvent常量）：
//   - 3: HOME（主屏幕键）
//   - 4: BACK（返回键）
//   - 24: VOLUME_UP（音量增加）
//   - 25: VOLUME_DOWN（音量减少）
//   - 26: POWER（电源键）
//   - 66: ENTER（回车键）
//   - 67: DEL（删除键）
//   - 82: MENU（菜单键）
//   - 84: SEARCH（搜索键）
//   - 85: MEDIA_PLAY_PAUSE（播放/暂停）
//   - 86: MEDIA_STOP（停止）
//   - 87: MEDIA_NEXT（下一曲）
//   - 88: MEDIA_PREVIOUS（上一曲）
//   - 122: MOVE_HOME（光标移至开头）
//   - 123: MOVE_END（光标移至末尾）
//   - 187: APP_SWITCH（应用切换/最近任务）
//
// 使用场景：
//   - 模拟硬件按键操作
//   - 控制媒体播放
//   - 导航操作
//   - 自动化测试
//
// 注意事项：
//   - 某些按键在特定应用中可能被拦截或禁用
//   - 部分按键需要特定权限
//   - 不同 Android 版本支持的按键可能有差异
//
// 示例：
//
//	// 发送返回键
//	err := device.KeyEvent(4)
//
//	// 发送回车键
//	err = device.KeyEvent(66)
//
//	// 音量增加
//	err = device.KeyEvent(24)
//
//	// 打开最近任务
//	err = device.KeyEvent(187)
func (d *Device) KeyEvent(keyCode int) error {
	command := fmt.Sprintf("input keyevent %d", keyCode)
	_, err := d.Shell(command)
	return err
}

// PressBack 模拟按下返回键（Back 键）。
// 该方法是 KeyEvent(4) 的便捷封装，用于执行返回操作。
//
// 返回值：
//   - error: 如果操作失败，返回 error 对象
//
// 使用场景：
//   - 返回上一级页面
//   - 关闭当前对话框或弹窗
//   - 退出当前应用到桌面（在应用首页时）
//   - 自动化导航操作
//
// 注意事项：
//   - 某些应用可能拦截返回键事件
//   - 在主屏幕按返回键通常无效果
//   - 部分全屏应用可能不响应返回键
//
// 示例：
//
//	// 简单的返回操作
//	err := device.PressBack()
//	if err != nil {
//	    log.Fatal("返回操作失败:", err)
//	}
//
//	// 连续返回多次
//	for i := 0; i < 3; i++ {
//	    device.PressBack()
//	    time.Sleep(500 * time.Millisecond)
//	}
func (d *Device) PressBack() error {
	return d.KeyEvent(4) // Android KeyEvent.KEYCODE_BACK = 4
}

// PressHome 模拟按下主屏幕键（Home 键）。
// 该方法是 KeyEvent(3) 的便捷封装，用于返回设备主屏幕。
//
// 返回值：
//   - error: 如果操作失败，返回 error 对象
//
// 使用场景：
//   - 快速返回主屏幕
//   - 退出当前应用
//   - 重置应用状态前的操作
//   - 自动化测试中的初始化操作
//
// 注意事项：
//   - 按下 Home 键会将当前应用置于后台
//   - 应用的生命周期会触发 onPause/onStop
//   - 某些系统可能有自定义的 Home 键行为
//
// 示例：
//
//	// 返回主屏幕
//	err := device.PressHome()
//	if err != nil {
//	    log.Fatal("返回主屏幕失败:", err)
//	}
//
//	// 退出应用后启动另一个应用
//	device.PressHome()
//	time.Sleep(1 * time.Second)
//	device.StartActivity("com.example.app", ".MainActivity")
func (d *Device) PressHome() error {
	return d.KeyEvent(3) // Android KeyEvent.KEYCODE_HOME = 3
}

// PressEnter 模拟按下回车键（Enter 键）。
// 该方法是 KeyEvent(66) 的便捷封装，常用于提交表单或确认输入。
//
// 返回值：
//   - error: 如果操作失败，返回 error 对象
//
// 使用场景：
//   - 提交表单
//   - 确认文本输入
//   - 发送消息
//   - 执行搜索
//   - 触发输入框的确认事件
//
// 注意事项：
//   - 需要确保输入框或可接收回车事件的控件获得焦点
//   - 在某些应用中，回车键可能触发换行而非提交
//   - 不同输入法对回车键的处理可能不同
//
// 示例：
//
//	// 输入文本后按回车提交
//	device.Tap(500, 600) // 点击输入框
//	time.Sleep(500 * time.Millisecond)
//	device.Input("search query")
//	device.PressEnter() // 提交搜索
//
//	// 登录表单提交
//	device.Input("username")
//	device.Tap(500, 800) // 切换到密码框
//	device.Input("password")
//	device.PressEnter() // 提交登录
func (d *Device) PressEnter() error {
	return d.KeyEvent(66) // Android KeyEvent.KEYCODE_ENTER = 66
}

// StartActivity 启动指定应用的指定 Activity（活动）。
// 该方法通过 'am start' 命令实现，可以直接启动应用的特定页面。
//
// 参数：
//   - packageName: 应用的包名（例如："com.android.settings"）
//   - activityName: Activity 的名称（例如：".Settings" 或完整类名）
//     如果以 "." 开头，会自动拼接包名
//     也可以使用完整的类名（包含包路径）
//
// 返回值：
//   - error: 如果启动失败，返回 error 对象
//
// Activity 名称格式：
//   - 简写形式: ".MainActivity"（自动拼接包名）
//   - 完整形式: "com.example.app.activities.MainActivity"
//   - 相对路径: "com.example.app/.MainActivity"（也是有效的）
//
// 使用场景：
//   - 启动特定应用
//   - 打开应用的特定页面
//   - 自动化测试中的应用初始化
//   - 深度链接测试
//   - 快速访问系统设置页面
//
// 注意事项：
//   - 应用必须已安装在设备上
//   - Activity 必须在 AndroidManifest.xml 中声明
//   - 某些 Activity 可能需要特定权限或参数
//   - 如果 Activity 名称错误，启动会失败
//
// 示例：
//
//	// 启动系统设置
//	err := device.StartActivity("com.android.settings", ".Settings")
//	if err != nil {
//	    log.Fatal("启动失败:", err)
//	}
//
//	// 启动自定义应用
//	err = device.StartActivity("com.example.myapp", ".MainActivity")
//
//	// 使用完整类名
//	err = device.StartActivity("com.example.myapp", "com.example.myapp.ui.SplashActivity")
//
//	// 启动浏览器打开网页（需要额外参数，使用 Shell 方法）
//	device.Shell("am start -a android.intent.action.VIEW -d https://www.example.com")
func (d *Device) StartActivity(packageName, activityName string) error {
	// 构建 am start 命令
	// -n: 指定组件名称（packageName/activityName）
	command := fmt.Sprintf("am start -n %s/%s", packageName, activityName)
	_, err := d.Shell(command)
	return err
}

// ForceStopApp 强制停止指定应用的所有进程。
// 该方法通过 'am force-stop' 命令实现，会完全终止应用。
//
// 参数：
//   - packageName: 要停止的应用包名（例如："com.example.app"）
//
// 返回值：
//   - error: 如果停止操作失败，返回 error 对象
//
// 工作原理：
//   - 杀死应用的所有进程
//   - 清除应用的后台任务
//   - 触发应用的完整生命周期结束
//   - 类似于在系统设置中"强行停止"应用
//
// 使用场景：
//   - 测试前重置应用状态
//   - 清理应用内存
//   - 停止异常运行的应用
//   - 自动化测试的清理操作
//   - 释放系统资源
//
// 注意事项：
//   - 应用的所有数据会丢失（未保存的状态）
//   - 后台服务会被终止
//   - 可能影响应用的通知、定时任务等
//   - 系统应用某些情况下无法被强制停止
//   - 停止后应用不会自动重启
//
// 示例：
//
//	// 停止应用
//	err := device.ForceStopApp("com.example.myapp")
//	if err != nil {
//	    log.Fatal("停止应用失败:", err)
//	}
//
//	// 重启应用前的清理
//	device.ForceStopApp("com.example.myapp")
//	time.Sleep(1 * time.Second) // 等待进程完全终止
//	device.StartActivity("com.example.myapp", ".MainActivity")
//
//	// 批量停止多个应用
//	apps := []string{"com.app1", "com.app2", "com.app3"}
//	for _, pkg := range apps {
//	    device.ForceStopApp(pkg)
//	}
func (d *Device) ForceStopApp(packageName string) error {
	// 构建 am force-stop 命令
	command := fmt.Sprintf("am force-stop %s", packageName)
	_, err := d.Shell(command)
	return err
}

// Pull 从 Android 设备拉取文件到本地计算机。
// 该方法通过 'adb pull' 命令实现文件传输。
//
// 参数：
//   - devicePath: 设备上的文件路径（例如："/sdcard/Download/file.txt"）
//   - localPath: 本地保存路径（例如："./downloads/file.txt"）
//
// 返回值：
//   - error: 如果拉取失败，返回 error 对象
//
// 使用场景：
//   - 备份设备文件到电脑
//   - 获取应用日志文件
//   - 下载截图或录屏
//   - 导出数据库文件
//   - 获取测试结果文件
//
// 注意事项：
//   - 需要有读取设备文件的权限
//   - 某些系统目录可能需要 root 权限
//   - 本地路径的目录必须存在
//   - 大文件传输可能需要较长时间
//   - 文件已存在时会被覆盖
//
// 示例：
//
//	// 拉取单个文件
//	err := device.Pull("/sdcard/screenshot.png", "./screenshot.png")
//	if err != nil {
//	    log.Fatal("拉取文件失败:", err)
//	}
//
//	// 拉取应用数据库（需要 root）
//	err = device.Pull("/data/data/com.example.app/databases/app.db", "./app.db")
//
//	// 拉取整个目录
//	err = device.Pull("/sdcard/DCIM/Camera/", "./camera_photos/")
//
//	// 拉取日志文件
//	err = device.Pull("/sdcard/Android/data/com.example.app/files/logs/app.log", "./app.log")
func (d *Device) Pull(devicePath, localPath string) error {
	// 执行 adb pull 命令
	_, err := d.execCommand("pull", devicePath, localPath)
	return err
}

// Push 从本地计算机推送文件到 Android 设备。
// 该方法通过 'adb push' 命令实现文件传输。
//
// 参数：
//   - localPath: 本地文件路径（例如："./test.txt"）
//   - devicePath: 设备上的目标路径（例如："/sdcard/test.txt"）
//
// 返回值：
//   - error: 如果推送失败，返回 error 对象
//
// 使用场景：
//   - 上传测试文件到设备
//   - 安装配置文件
//   - 传输媒体文件
//   - 部署测试数据
//   - 更新应用资源文件
//
// 注意事项：
//   - 需要有写入设备路径的权限
//   - 某些系统目录可能需要 root 权限
//   - 本地文件必须存在
//   - 设备存储空间必须充足
//   - 文件已存在时会被覆盖
//   - /sdcard/ 通常是普通应用可写的位置
//
// 示例：
//
//	// 推送单个文件
//	err := device.Push("./test_data.json", "/sdcard/Download/test_data.json")
//	if err != nil {
//	    log.Fatal("推送文件失败:", err)
//	}
//
//	// 推送配置文件
//	err = device.Push("./config.xml", "/sdcard/Android/data/com.example.app/files/config.xml")
//
//	// 推送图片
//	err = device.Push("./photo.jpg", "/sdcard/Pictures/photo.jpg")
//
//	// 推送整个目录
//	err = device.Push("./test_files/", "/sdcard/test_files/")
func (d *Device) Push(localPath, devicePath string) error {
	// 执行 adb push 命令
	_, err := d.execCommand("push", localPath, devicePath)
	return err
}

// Connect 通过 TCP/IP 网络连接到指定地址的 Android 设备。
// 该方法通过 'adb connect' 命令实现无线 ADB 连接。
//
// 参数：
//   - address: 设备的网络地址，格式为 "IP:端口"（例如："192.168.1.100:5555"）
//     默认端口是 5555，如果使用默认端口可以省略（但建议明确指定）
//
// 返回值：
//   - error: 如果连接失败，返回 error 对象
//
// 前置条件：
//  1. 设备必须开启 ADB 网络调试功能
//  2. 设备和电脑必须在同一网络中
//  3. 设备的 ADB 守护进程必须监听指定端口
//
// 开启方式：
//   - 通过 USB 先连接设备，然后执行：
//     adb tcpip 5555  // 开启网络调试
//     adb connect <设备IP>:5555  // 连接设备
//
// 使用场景：
//   - 无线调试和测试
//   - USB 接口不可用时的替代方案
//   - 多设备并行测试
//   - 远程设备控制
//
// 注意事项：
//   - 网络连接不如 USB 稳定
//   - 传输速度较慢
//   - 需要确保防火墙允许连接
//   - 安全性较低，仅在可信网络中使用
//
// 示例：
//
//	// 连接局域网设备
//	err := adb.Connect("192.168.1.100:5555")
//	if err != nil {
//	    log.Fatal("连接设备失败:", err)
//	}
//
//	// 连接后创建设备实例
//	err = adb.Connect("192.168.1.100:5555")
//	if err == nil {
//	    device := adb.NewDevice("192.168.1.100:5555")
//	}
//
//	// 连接多个设备
//	addresses := []string{"192.168.1.100:5555", "192.168.1.101:5555"}
//	for _, addr := range addresses {
//	    if err := adb.Connect(addr); err != nil {
//	        log.Printf("连接 %s 失败: %v", addr, err)
//	    }
//	}
func Connect(address string) error {
	// 执行 adb connect 命令
	_, err := exec.Command("adb", "connect", address).CombinedOutput()
	return err
}

// GetClipper 从设备剪贴板获取文本内容。
// 该方法依赖第三方应用 Clipper (ca.zgrs.clipper) 来读取剪贴板内容。
//
// 返回值：
//   - string: 剪贴板中的文本内容
//   - error: 如果获取失败，返回 error 对象
//
// 前置条件：
//   - 设备上必须安装 Clipper 应用（包名：ca.zgrs.clipper）
//   - Clipper 应用必须有剪贴板读取权限
//
// 工作原理：
//  1. 启动 Clipper 应用
//  2. 等待应用启动（1秒）
//  3. 通过广播获取剪贴板内容
//  4. 解析输出并返回文本
//
// 使用场景：
//   - 读取设备剪贴板内容
//   - 验证复制操作是否成功
//   - 自动化测试中的数据验证
//   - 获取应用间传递的数据
//
// 注意事项：
//   - 需要先安装 Clipper 应用
//   - 启动应用需要时间，已内置 1 秒等待
//   - 如果剪贴板为空，可能返回空字符串
//   - 输出格式必须包含 "data=" 字段
//
// 示例：
//
//	// 获取剪贴板内容
//	text, err := device.GetClipper()
//	if err != nil {
//	    log.Fatal("获取剪贴板失败:", err)
//	}
//	fmt.Println("剪贴板内容:", text)
//
//	// 验证复制操作
//	// 1. 执行复制操作
//	device.Tap(500, 600) // 点击文本
//	device.KeyEvent(278) // 模拟复制操作（KEYCODE_COPY）
//	time.Sleep(500 * time.Millisecond)
//	// 2. 验证剪贴板内容
//	text, err := device.GetClipper()
//	if err == nil && strings.Contains(text, "expected content") {
//	    fmt.Println("复制成功")
//	}
func (d *Device) GetClipper() (string, error) {
	// 启动 Clipper 应用
	err := d.StartActivity("ca.zgrs.clipper", "ca.zgrs.clipper.Main")
	if err != nil {
		return "", err
	}

	// 等待应用完全启动
	time.Sleep(1 * time.Second)

	// 通过广播获取剪贴板内容
	output, err := d.Shell("am broadcast -a clipper.get")
	if err != nil {
		return "", err
	}

	// 解析输出，提取 data= 后面的内容
	parts := strings.Split(output, "data=")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output: %s", output)
	}

	// 返回剪贴板文本（去除首尾空白）
	clipText := strings.TrimSpace(parts[1])
	return clipText, nil
}
