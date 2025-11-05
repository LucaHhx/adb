package adb

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Shell 在设备上执行 shell 命令
func (d *Device) Shell(command string) (string, error) {
	return d.execCommand("shell", command)
}

// Tap 在指定坐标处执行点击操作
func (d *Device) Tap(x, y int) error {
	command := fmt.Sprintf("input tap %d %d", x, y)
	_, err := d.Shell(command)
	return err
}

// Swipe 在指定坐标之间执行滑动操作
func (d *Device) Swipe(x1, y1, x2, y2, duration int32) error {
	command := fmt.Sprintf("input swipe %d %d %d %d %d", x1, y1, x2, y2, duration)
	_, err := d.Shell(command)
	return err
}

// Input 向设备发送文本输入
func (d *Device) Input(text string) error {
	// 将空格替换为 %s 以适配 ADB input 命令
	escapedText := strings.ReplaceAll(text, " ", "%s")
	command := fmt.Sprintf("am broadcast -a ADB_INPUT_TEXT --es msg '%s'", escapedText)
	//adb shell am broadcast -a ADB_INPUT_TEXT --es msg 'Machinery@999'
	_, err := d.Shell(command)
	return err
}

// KeyEvent 向设备发送按键事件
func (d *Device) KeyEvent(keyCode int) error {
	command := fmt.Sprintf("input keyevent %d", keyCode)
	_, err := d.Shell(command)
	return err
}

// PressBack 模拟按下返回键
func (d *Device) PressBack() error {
	return d.KeyEvent(4) // 返回键代码
}

// PressHome 模拟按下主屏幕键
func (d *Device) PressHome() error {
	return d.KeyEvent(3) // 主屏幕键代码
}

// PressEnter 模拟按下回车键
func (d *Device) PressEnter() error {
	return d.KeyEvent(66) // 回车键代码
}

// StartActivity 启动一个活动
func (d *Device) StartActivity(packageName, activityName string) error {
	command := fmt.Sprintf("am start -n %s/%s", packageName, activityName)
	_, err := d.Shell(command)
	return err
}

// ForceStopApp 强制停止一个应用
func (d *Device) ForceStopApp(packageName string) error {
	command := fmt.Sprintf("am force-stop %s", packageName)
	_, err := d.Shell(command)
	return err
}

// Pull 从设备拉取文件到本地
func (d *Device) Pull(devicePath, localPath string) error {
	_, err := d.execCommand("pull", devicePath, localPath)
	return err
}

// Push 从本地推送文件到设备
func (d *Device) Push(localPath, devicePath string) error {
	_, err := d.execCommand("push", localPath, devicePath)
	return err
}

// Connect 连接到指定的设备
func Connect(address string) error {
	_, err := exec.Command("adb", "connect", address).CombinedOutput()
	return err
}

func (d *Device) GetClipper() (string, error) {
	err := d.StartActivity("ca.zgrs.clipper", "ca.zgrs.clipper.Main")
	if err != nil {
		return "", err
	}
	time.Sleep(1 * time.Second) // 等待应用启动
	output, err := d.Shell("am broadcast -a clipper.get")
	if err != nil {
		return "", err
	}
	parts := strings.Split(output, "data=")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output: %s", output)
	}
	clipText := strings.TrimSpace(parts[1])
	return clipText, nil
}
