# ---- Stage 1: Build ----
# 使用官方 Node.js image 作為建置環境
FROM node:18-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 package.json 和 package-lock.json
COPY package*.json ./

# 安裝依賴
RUN npm install

# 複製所有原始碼
COPY . .

# 執行建置命令
RUN npm run build

# ---- Stage 2: Production ----
# 使用一個更小的 Node.js image
FROM node:18-alpine

WORKDIR /app

# 從建置階段複製必要的檔案
# next.config.mjs, public, 和 .next (standalone 模式)
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# 暴露 Next.js 預設的端口
EXPOSE 3000

# 啟動 Node.js 伺服器
CMD ["node", "server.js"]