FROM node:20-alpine
WORKDIR /app
COPY package.json ./
COPY angular.json ./
COPY tsconfig.json ./
COPY tsconfig.app.json ./
COPY src ./src
EXPOSE 3000
CMD ["sh", "-lc", "corepack enable && pnpm install && pnpm start"]

