FROM node:22-alpine

WORKDIR /work/app

# 開発時は起動コマンドを docker-compose の command で上書きする
CMD ["sh", "-c", "npm install && npm run dev -- --host 0.0.0.0 --port 3000"]
