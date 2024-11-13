# 第一階段：構建階段
FROM golang:1.23.3 AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum（如果有的話）
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製其餘代碼文件
COPY . .

# 編譯 Go 應用
RUN go build -o CloudflareSpeedTest

# 第二階段：運行階段
FROM alpine:3.20

# 設置字符編碼和時區（可選）
ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

# 設置工作目錄
WORKDIR /app

# 複製編譯後的二進制文件到運行映像
COPY --from=builder /app/CloudflareSpeedTest .

# 複製默認配置文件到 /config 目錄
COPY config.json /config/config.json

# 設置可映射的配置和數據目錄
VOLUME ["/config", "/data"]

# 設置默認執行命令
CMD ["./CloudflareSpeedTest", "-c", "/config/config.json"]