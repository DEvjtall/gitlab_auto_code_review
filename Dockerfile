FROM  golang:1.23.1 AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
        && go env -w GOPROXY=https://goproxy.cn,direct \
        && go env -w CGO_ENABLED=0 \
        && go env \
        && go mod tidy \
        && go build -o code-review main.go

# 执行过程
FROM apline:latest
WORKDIR /app
COPY --from=builder code-review .
COPY --from=builder config/config.ini ./config/
CMD ["./code-review"]