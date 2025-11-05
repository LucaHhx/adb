package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

// Device 代表通过 ADB 连接的 Android 设备
type Device struct {
	Serial string
}

// NewDevice 创建一个新的 Device 实例
// 如果 serial 为空，将使用默认设备
func NewDevice(serial ...string) *Device {
	if len(serial) == 0 {
		return &Device{
			Serial: "",
		}
	}
	return &Device{
		Serial: serial[0],
	}
}

// execCommand 执行 ADB 命令并返回输出
func (d *Device) execCommand(args ...string) (string, error) {
	cmdArgs := []string{}

	// 如果指定了设备序列号，则添加
	if d.Serial != "" {
		cmdArgs = append(cmdArgs, "-s", d.Serial)
	}

	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("adb", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("adb command failed: %w, output: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
