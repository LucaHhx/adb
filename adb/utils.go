package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetDevices 获取所有已连接设备的序列号列表。
// 该函数会执行 'adb devices' 命令，解析输出并返回所有状态为 "device" 的设备序列号。
//
// 返回值：
//   - []string: 已连接设备的序列号列表（数组）
//   - error: 如果命令执行失败或解析出错，返回 error 对象
//
// 设备状态说明：
//   - device: 设备已连接并正常工作（会被返回）
//   - offline: 设备离线或未响应（不会被返回）
//   - unauthorized: 设备未授权，需要在设备上确认 USB 调试（不会被返回）
//   - no devices: 没有设备连接
//
// 使用场景：
//   - 检查是否有设备连接
//   - 列出所有可用设备供用户选择
//   - 在多设备环境下进行设备管理
//
// 注意事项：
//   - 返回的序列号可能是设备的物理序列号，也可能是网络地址（如 192.168.1.100:5555）
//   - 只返回状态为 "device" 的设备，其他状态的设备不会被返回
//   - 如果没有设备连接，会返回空数组（不是 nil）
//
// 示例：
//
//	// 获取所有连接的设备
//	devices, err := adb.GetDevices()
//	if err != nil {
//	    log.Fatal("获取设备列表失败:", err)
//	}
//
//	if len(devices) == 0 {
//	    fmt.Println("没有连接的设备")
//	} else {
//	    fmt.Printf("发现 %d 个设备:\n", len(devices))
//	    for i, serial := range devices {
//	        fmt.Printf("%d. %s\n", i+1, serial)
//	    }
//	}
//
//	// 使用第一个设备
//	if len(devices) > 0 {
//	    device := adb.NewDevice(devices[0])
//	}
func GetDevices() ([]string, error) {
	// 执行 'adb devices' 命令
	cmd := exec.Command("adb", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	// 按行分割输出
	lines := strings.Split(string(output), "\n")
	var devices []string

	// 遍历每一行输出
	for i, line := range lines {
		// 跳过第一行（标题："List of devices attached"）和空行
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// 将行按空白字符分割为多个字段
		// 正常格式: "序列号  状态"（例如："emulator-5554  device"）
		parts := strings.Fields(line)

		// 只添加状态为 "device" 的设备（正常连接的设备）
		if len(parts) >= 2 && parts[1] == "device" {
			devices = append(devices, parts[0])
		}
	}

	return devices, nil
}

// WaitForDevice 等待指定设备连接并进入就绪状态。
// 该函数会阻塞执行，直到设备连接成功或发生错误。
//
// 参数：
//   - serial: 设备序列号
//     如果为空字符串，则等待任意设备连接
//     如果指定序列号，则等待该特定设备连接
//
// 返回值：
//   - error: 如果等待过程中发生错误，返回 error 对象；成功时返回 nil
//
// 使用场景：
//   - 脚本启动时等待设备连接
//   - 设备重启后等待其重新连接
//   - 自动化测试中确保设备就绪后再执行操作
//   - 批处理脚本中同步设备状态
//
// 工作原理：
//   - 该函数调用 'adb wait-for-device' 命令
//   - 命令会阻塞直到设备进入 "device" 状态
//   - 如果设备已连接，命令会立即返回
//   - 如果设备未连接，会持续等待直到设备连接
//
// 注意事项：
//   - 该函数会一直阻塞，直到设备连接或出现错误
//   - 没有超时机制，如果设备一直未连接，函数会一直等待
//   - 建议在使用时设置自己的超时控制机制
//   - 该函数不会验证设备是否已授权 USB 调试
//
// 示例：
//
//	// 等待任意设备连接
//	fmt.Println("等待设备连接...")
//	err := adb.WaitForDevice("")
//	if err != nil {
//	    log.Fatal("等待设备失败:", err)
//	}
//	fmt.Println("设备已连接!")
//
//	// 等待特定设备连接
//	serial := "emulator-5554"
//	fmt.Printf("等待设备 %s 连接...\n", serial)
//	err = adb.WaitForDevice(serial)
//	if err != nil {
//	    log.Fatal("等待设备失败:", err)
//	}
//	fmt.Printf("设备 %s 已连接!\n", serial)
//
//	// 带超时控制的等待
//	done := make(chan error, 1)
//	go func() {
//	    done <- adb.WaitForDevice("")
//	}()
//
//	select {
//	case err := <-done:
//	    if err != nil {
//	        log.Fatal("等待设备失败:", err)
//	    }
//	    fmt.Println("设备已连接!")
//	case <-time.After(30 * time.Second):
//	    log.Fatal("等待设备超时（30秒）")
//	}
func WaitForDevice(serial string) error {
	// 构建命令参数
	args := []string{"wait-for-device"}

	// 如果指定了设备序列号，添加 "-s serial" 参数
	if serial != "" {
		args = append([]string{"-s", serial}, args...)
	}

	// 执行命令并等待完成
	cmd := exec.Command("adb", args...)
	return cmd.Run()
}
