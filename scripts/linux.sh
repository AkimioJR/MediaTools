#!/bin/bash

set -e  # 如果任何命令失败就退出

echo "正在更新软件包列表..."
apt-get update

echo "正在安装必要的开发库..."
apt-get install -y \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    pkg-config