ARG NODE_VERSION=24.7.0-alpine

FROM node:${NODE_VERSION} AS builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN --mount=type=cache,target=/root/.npm npm ci

EXPOSE 4200

CMD [ "npm", "run", "start" ]