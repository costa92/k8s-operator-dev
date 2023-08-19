#!/bin/bash

code_server_name="code-server-4.11.0-linux-arm64.tar.gz"
golang_version_name="go1.19.4.linux-amd64.tar.gz"

code_server_name_version=$(echo "$code_server_name" | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+')
echo ${code_server_name_version}
# 判断文件是否存在
if [ ! -f "$code_server_name" ]; then
	echo "start download code-server version:"$code_server_name_version
	wget "https://github.com/coder/code-server/releases/download/v"$code_server_name_version"/"$code_server_name
	echo "download code-server success"
else
	echo "code-server is exist"
fi

# 判断golang_version_name是否存在
if [[ ! -f "$golang_version_name" ]]; then
	echo "start download golang"
	wget "https://studygolang.com/dl/golang/"$golang_version_name
	echo "download golang success"
else
	echo "golang is exist"
fi
