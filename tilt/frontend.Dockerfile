FROM node:23-alpine
WORKDIR /app

ARG node_env=production

ENV NODE_ENV $node_env

ADD package.json .
ADD yarn.lock .
RUN yarn install

ADD . .

EXPOSE 3000