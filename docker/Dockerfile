# 前端构建阶段
FROM node:24-alpine AS web-builder
WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm run build

# Go 后端构建阶段
# 使用低版本通过 GOTOOLCHAIN 自动下载更高版本
FROM golang:1.24.0-alpine AS go-builder 


WORKDIR /build

RUN apk update && \
    apk add git build-base gcc musl-dev binutils-gold

# 设置环境变量以启用工具链自动下载
ENV GOTOOLCHAIN=auto
ENV GO111MODULE=on
ENV CGO_ENABLED=1 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制前端构建产物
COPY --from=web-builder /web/dist ./web/dist

RUN go run scripts/build.go -output MediaTools

# 最终运行阶段
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=go-builder /build/MediaTools .
RUN chmod +x ./MediaTools

EXPOSE 5000
VOLUME ["/app/data", "/app/logs"]

ENTRYPOINT ["./MediaTools"]
CMD ["-server"]