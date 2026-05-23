FROM node:20-alpine
WORKDIR /app
COPY package.json ./
COPY angular.json ./
COPY tsconfig.json ./
COPY tsconfig.app.json ./
COPY src ./src
EXPOSE 9960
CMD ["sh", "-lc", "npm install && npm run start"]
