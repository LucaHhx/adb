package main

import (
	"github.com/LucaHhx/adb/adb"
)

// 强制停止应用：adb shell am force-stop id.co.bri.biz.app.token
// 启动应用：adb shell am start -n id.co.bri.biz.app.token/id.co.bri.biz.app.nyx.MainActivity
// 点击坐标：adb shell input tap 271 1000
// 输入文本：adb shell input text "Cosmetic999"
// 点击坐标：adb shell input tap 313 959
// 点击坐标：adb shell input tap 362 1452
// 点击坐标：adb shell input tap 540 1487
// 输入文本：adb shell input text "Cosmetic999"
func main() {
	dev := adb.NewDevice()
	//str, err := dev.UiautomatorDump()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//fmt.Println(str)
	clipper, err := dev.GetClipper()
	if err != nil {
		return
	}

	println(clipper)
}
