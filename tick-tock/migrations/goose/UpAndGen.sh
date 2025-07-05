#!/bin/bash

# Goose 数据库迁移
goose mysql "root:root@tcp(127.0.0.1:33060)/micro_tt?parseTime=true&loc=UTC" up
if [ $? -ne 0 ]; then
    echo "Goose migration failed"
    exit 1
fi

# 进入 gentool 目录并运行代码生成
cd "$(dirname "$0")/../gentool" || { echo "Failed to enter gentool directory"; exit 1; }
gentool -c "db2struct.yaml"
if [ $? -ne 0 ]; then
    echo "gentool execution failed"
    exit 1
fi
