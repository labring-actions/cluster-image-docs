#!/bin/bash

# 仓库名称
repository="labring-actions/repos"

# 获取最新release的版本号
latest_release=$(curl -s "https://api.github.com/repos/$repository/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
# 构建下载链接
download_url="https://github.com/$repository/releases/download/$latest_release/repos_${latest_release#v}_linux_amd64.tar.gz"

# 下载最新release
wget $download_url

# 解压缩下载的文件（如果是tar.gz格式）
tar -zxvf repos_${latest_release#v}_linux_amd64.tar.gz repos

# 删除压缩包
rm -rf repos_${latest_release#v}_linux_amd64.tar.gz

chmod a+x repos

mkdir "/tmp"

mv srepos "/tmp"
