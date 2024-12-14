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

# 編譯 Go 應用，設置目標平台為 linux/amd64
# RUN GOARCH=amd64 GOOS=linux go build -o CloudflareSpeedTest
RUN GOARCH=${TARGETARCH} GOOS=linux go build -o CloudflareSpeedTest

# 第二階段：運行階段
FROM alpine:3.20

# 安裝 glibc 兼容庫（如果需要）
RUN apk --no-cache add libc6-compat

RUN apk --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 安裝 traceroute 工具
RUN apk --no-cache add traceroute

USER root
RUN echo "Use Root User"
# 安裝 curl 工具
RUN apk update
RUN echo "Use apk update"
RUN apk add --no-cache curl
RUN echo "Install  curl"
RUN apt-add --no-cache install -y bash
RUN echo "Install  bash"

# 安裝 nxtrace 工具RUN curl -sSL nxtrace.org/nt | bash || echo "nxtrace installation failed, skipping."
RUN curl -sSL nxtrace.org/nt | bash 
RUN echo "Install  nxtrace"

# > /app/nxtrace.install.log 2>&1 || echo "nxtrace installation failed, skipping." >> /app/nxtrace.install.log


# 設置字符編碼和時區（可選）
ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

# 创建设置时区的入口脚本
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# 設置工作目錄
WORKDIR /app

# 複製編譯後的二進制文件到運行映像
COPY --from=builder /app/CloudflareSpeedTest /app/CloudflareSpeedTest

# 複製 static 目錄到運行映像
COPY --from=builder /app/static /app/static

# 設置執行權限
RUN chmod +x /app/CloudflareSpeedTest

# 創建目錄
RUN mkdir /config/

# 複製默認配置文件到 /config 目錄
COPY config.json /config/config.json

# 添加当前时间到 versionDate.txt
RUN date '+%Y-%m-%d %H:%M:%S' > /app/versionDate.txt

# 設置可映射的配置和數據目錄
VOLUME ["/config", "/data"]

# 曝露端口 8080
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

# 設置默認執行命令
CMD ["/app/CloudflareSpeedTest", "-c", "/config/config.json"]
