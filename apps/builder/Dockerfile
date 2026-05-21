FROM node:20-alpine
WORKDIR /app
COPY package.json ./
COPY astro.config.mjs ./
COPY tsconfig.json ./
COPY src ./src
COPY public ./public
EXPOSE 4321
CMD ["sh", "-lc", "corepack enable && pnpm install && pnpm dev"]

