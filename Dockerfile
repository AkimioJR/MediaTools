# 前端构建阶段
FROM --platform=linux/amd64 node:24-alpine AS web-builder
WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm run build

# Go 后端构建阶段
# 使用低版本通过 GOTOOLCHAIN 自动下载更高版本
FROM golang:1.23-alpine AS go-builder 

ARG APP_VERSION="dev-docker"
ARG APP_GIT_COMMIT=""
ARG APP_BUILT_AT="unknown"

WORKDIR /build

RUN apk update && \
    apk add git build-base

# 设置环境变量以启用工具链自动下载
ENV GOTOOLCHAIN=auto
ENV GO111MODULE=on
ENV CGO_ENABLED=1 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 复制前端构建产物
COPY --from=web-builder /web/dist ./web/dist

RUN go build \
    -tags=onlyServer \
    -o MediaTools \
    -ldflags="-s -w \
              -X MediaTools/internal/version.appVersion=$APP_VERSION \
              -X MediaTools/internal/version.commitHash=$APP_GIT_COMMIT \
              -X MediaTools/internal/version.buildTime=$APP_BUILT_AT" \
    .

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