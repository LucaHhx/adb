// Package adb 提供了对 Android Debug Bridge (ADB) 的 Go 语言封装。
// 该包允许开发者通过编程方式控制 Android 设备，执行各种自动化操作。
package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

// Device 代表通过 ADB 连接的 Android 设备实例。
// 每个 Device 实例可以通过设备序列号（Serial）来标识唯一的设备。
// 如果 Serial 为空字符串，则默认操作系统中唯一连接的设备。
//
// 字段说明：
//   - Serial: 设备的序列号，可通过 'adb devices' 命令查看
//
// 示例：
//
//	// 创建默认设备（系统唯一连接的设备）
//	device := adb.NewDevice()
//
//	// 创建指定序列号的设备
//	device := adb.NewDevice("emulator-5554")
type Device struct {
	Serial string // 设备序列号，为空时使用默认设备
}

// NewDevice 创建一个新的 Device 实例。
//
// 参数：
//   - serial: 可选的设备序列号（可变参数）
//     如果不传入任何参数，将创建一个默认设备实例（操作唯一连接的设备）
//     如果传入序列号，将创建一个指定设备的实例
//
// 返回值：
//   - *Device: 新创建的设备实例指针
//
// 使用场景：
//   - 当系统只连接一个设备时，使用 NewDevice() 即可
//   - 当系统连接多个设备时，需要指定序列号，如 NewDevice("emulator-5554")
//
// 示例：
//
//	// 创建默认设备
//	device := adb.NewDevice()
//
//	// 创建指定序列号的设备
//	device := adb.NewDevice("192.168.1.100:5555")
//
//	// 从设备列表中选择
//	devices, _ := adb.GetDevices()
//	if len(devices) > 0 {
//	    device := adb.NewDevice(devices[0])
//	}
func NewDevice(serial ...string) *Device {
	// 如果没有传入序列号参数，返回默认设备
	if len(serial) == 0 {
		return &Device{
			Serial: "",
		}
	}
	// 返回指定序列号的设备
	return &Device{
		Serial: serial[0],
	}
}

// execCommand 执行 ADB 命令并返回输出结果。
// 这是一个内部方法，用于封装所有 ADB 命令的执行逻辑。
//
// 参数：
//   - args: 要执行的 ADB 命令参数（可变参数）
//     例如：执行 "adb shell ls" 命令时，传入 "shell", "ls"
//
// 返回值：
//   - string: 命令执行的标准输出（已去除首尾空白字符）
//   - error: 如果命令执行失败，返回包含错误信息的 error 对象
//
// 执行流程：
//  1. 检查是否指定了设备序列号，如果指定了则添加 "-s serial" 参数
//  2. 拼接完整的命令参数
//  3. 执行命令并获取输出（包括标准输出和标准错误）
//  4. 如果执行失败，返回错误信息
//  5. 返回去除首尾空白的输出结果
//
// 错误处理：
//   - 命令执行失败时，错误信息中会包含原始错误和命令输出
//   - 这有助于诊断问题，如设备未连接、权限不足等
//
// 注意事项：
//   - 该方法使用 CombinedOutput()，同时捕获标准输出和标准错误
//   - 输出会自动去除首尾的空白字符（空格、换行符等）
//   - 如果设备未连接或 ADB 未安装，会返回相应错误
//
// 示例：
//
//	// 内部调用示例（一般不直接调用此方法）
//	output, err := device.execCommand("shell", "getprop", "ro.build.version.release")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Android version:", output)
func (d *Device) execCommand(args ...string) (string, error) {
	// 初始化命令参数切片
	cmdArgs := []string{}

	// 如果指定了设备序列号，添加 "-s serial" 参数
	// 这确保命令在正确的设备上执行
	if d.Serial != "" {
		cmdArgs = append(cmdArgs, "-s", d.Serial)
	}

	// 将传入的命令参数追加到参数列表
	cmdArgs = append(cmdArgs, args...)

	// 创建 ADB 命令
	cmd := exec.Command("adb", cmdArgs...)

	// 执行命令并获取输出（包括 stdout 和 stderr）
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 命令执行失败，返回详细的错误信息
		return "", fmt.Errorf("adb command failed: %w, output: %s", err, string(output))
	}

	// 返回去除首尾空白字符的输出结果
	return strings.TrimSpace(string(output)), nil
}
