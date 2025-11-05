#!/bin/bash

set -e

echo "开始初始化 Android 设备..."

# 检查设备连接
echo "检查设备连接..."
adb devices

# 安装 APK
echo "安装 geto.apk..."
adb install -r ./geto.apk

echo "安装 adbkey.apk..."
adb install -r ./adbkey.apk

echo "安装 clipper.apk (使用 --bypass-low-target-sdk-block)..."
adb install -r --bypass-low-target-sdk-block ./clipper.apk

# 配置系统设置
echo "启用触摸显示..."
adb shell settings put system show_touches 1

echo "启用指针位置..."
adb shell settings put system pointer_location 1

# 授予权限
echo "授予 Geto 写入安全设置权限..."
adb shell pm grant com.android.geto android.permission.WRITE_SECURE_SETTINGS

# 配置输入法
echo "启用 ADB 键盘..."
adb shell ime enable com.android.adbkeyboard/.AdbIME

echo "设置 ADB 键盘为默认输入法..."
adb shell ime set com.android.adbkeyboard/.AdbIME

echo "初始化完成！"
