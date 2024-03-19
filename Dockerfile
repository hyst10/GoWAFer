# 使用官方Go镜像作为构建环境
FROM golang:1.20-alpine AS builder

# 设置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .


# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# 使用scratch作为最小运行时容器
FROM scratch

# 从构建者镜像中复制编译好的应用到当前目录
COPY --from=builder /app/myapp .


COPY --from=builder /app/config.yaml .
COPY --from=builder /app/static /static
COPY --from=builder /app/templates /templates
COPY --from=builder /app/pages /pages



# 运行应用
CMD ["./myapp"]
