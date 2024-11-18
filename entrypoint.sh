#!/bin/sh

# 动态设置时区
if [ -n "$TZ" ]; then
    echo "Setting timezone to $TZ"
    # 确保时区文件存在并设置
    if [ -f "/usr/share/zoneinfo/$TZ" ]; then
        cp /usr/share/zoneinfo/$TZ /etc/localtime
        echo "$TZ" > /etc/timezone
    else
        echo "Invalid timezone: $TZ" >&2
        exit 1
    fi
else
    echo "No TZ environment variable provided, using default timezone: Asia/Shanghai"
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
    echo "Asia/Shanghai" > /etc/timezone
fi

# 执行主程序
exec "$@"
