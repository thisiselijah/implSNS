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

# 將 backend/ 的所有內容複製到容器的 /app/ 中
# 這會包含你的 cmd/, internal/, 以及 config/ 等資料夾
COPY . .

# 編譯 Go 應用程式，假設你的主程式入口在 ./cmd/
# 將編譯結果輸出到 /app/main
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/backend/

# ---- Stage 2: Production ----
# 使用一個極度輕量的空白 image
FROM scratch

# # 【*** 修正點 ***】
# # 從 builder 階段複製 CA 憑證到正式環境中
# # Go 應用程式會自動在這些預設路徑尋找憑證
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

# 保留 CA 憑證，讓 Go 程式能正確驗證 HTTPS
COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 複製本地已編譯好的執行檔
COPY main /main

# 複製 config 資料夾（如果有設定檔需求）
COPY config /config

EXPOSE 8080

ENTRYPOINT ["/main"]