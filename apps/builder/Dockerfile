FROM node:20-alpine
WORKDIR /workspace
COPY package.json package-lock.json ./
COPY apps/builder/package.json ./apps/builder/package.json
RUN npm ci
COPY apps/builder ./apps/builder
COPY packages/templates/supromail ./apps/builder/templates/supromail
WORKDIR /workspace/apps/builder
EXPOSE 4321
CMD ["npm", "run", "dev"]
