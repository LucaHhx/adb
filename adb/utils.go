package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetDevices 返回已连接设备的序列号列表
func GetDevices() ([]string, error) {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var devices []string

	for i, line := range lines {
		// 跳过第一行（标题）和空行
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "device" {
			devices = append(devices, parts[0])
		}
	}

	return devices, nil
}

// WaitForDevice 等待设备连接
func WaitForDevice(serial string) error {
	args := []string{"wait-for-device"}
	if serial != "" {
		args = append([]string{"-s", serial}, args...)
	}

	cmd := exec.Command("adb", args...)
	return cmd.Run()
}
