FROM node:23-alpine3.19 AS builder

# Build stage
WORKDIR /app
COPY package*.json .
RUN yarn install --immutable
COPY . .
RUN yarn run build

FROM node:23-alpine
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .
EXPOSE 3000
ENV NODE_ENV=production
CMD [ "node", "build" ]