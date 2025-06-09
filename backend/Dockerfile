# ---- Stage 1: Builder ----
# 使用官方 Go image 作為建置環境
FROM golang:1.24-alpine AS builder

# 安裝 ca-certificates 套件確保憑證存在
RUN apk add --no-cache ca-certificates

# 設定工作目錄
WORKDIR /app

# 複製 go module 檔案並下載依賴，以利用 Docker 的 layer cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ---- Stage 2: Development ----

# RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/backend/

# FROM scratch

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# # 將最終執行檔從 builder 階段複製過來
# COPY --from=builder /app/main /main

# # 將整個 config 資料夾從 builder 階段複製過來
# # 這樣你的執行檔就能找到設定檔了
# COPY --from=builder /app/config /config

# # 暴露你的應用程式監聽的 port
# EXPOSE 8080

# # 設定容器啟動時執行的命令
# ENTRYPOINT ["/main"]


# ---- Stage 2: Production ----

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# 複製本地已編譯好的執行檔
COPY main /main

# 複製 config 資料夾（如果有設定檔需求）
COPY config /config

EXPOSE 8080

ENTRYPOINT ["/main"]