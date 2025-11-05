package adb

import (
	"fmt"
	"regexp"
	"strings"
)

// Execout 执行 adb exec-out 命令
func (d *Device) Execout(command string) (string, error) {
	return d.execCommand("exec-out", command)
}

var errorMap = map[string]string{
	"Mohon periksa koneksi internet Anda.": "请检查您的互联网连接。",
}

func (d *Device) UiautomatorDump() (string, error) {
	command, err := d.execCommand("exec-out", "uiautomator dump /dev/tty")
	if err != nil {
		return "", err
	}
	for s, s2 := range errorMap {
		if strings.Contains(command, s) {
			return "", fmt.Errorf(s2)
		}
	}
	return command, nil
}

// ExistElement 检查 UI 元素是否存在
func (d *Device) ExistElement(content string) bool {
	command, err := d.UiautomatorDump()
	if err != nil {
		return false
	}
	return strings.Contains(command, content)
}

func (d *Device) Regexp(rex string) (string, error) {
	command, err := d.UiautomatorDump()
	if err != nil {
		return "", err
	}
	// 构造正则：匹配含有该 bounds 的 node 并提取 content-desc
	re := regexp.MustCompile(rex)

	matches := re.FindStringSubmatch(command)
	if len(matches) > 1 {
		return matches[1], nil
	} else {
		return "", fmt.Errorf("not found")
	}
}

func (d *Device) FindDesc(bounds string) (string, error) {
	data, err := d.UiautomatorDump()
	if err != nil {
		return "", err
	}
	nodeRe := regexp.MustCompile(`<node\b[^>]*\bbounds="` + regexp.QuoteMeta(bounds) + `"[^>]*/>`)
	nodeMatch := nodeRe.FindString(data)
	if nodeMatch == "" {
		fmt.Println("未找到目标节点")
		return "", nil
	}

	// ② 再从该节点中提取 content-desc
	contentRe := regexp.MustCompile(`content-desc="([^"]*)"`)
	content := contentRe.FindStringSubmatch(nodeMatch)
	if len(content) > 1 {
		return content[1], nil
	} else {
		fmt.Println("未找到 content-desc 属性")
	}
	return "", nil
}
