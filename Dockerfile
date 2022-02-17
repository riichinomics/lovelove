FROM node:16.12 AS node-build
WORKDIR /build/client
COPY ./client/package.json ./client/yarn.lock ./
RUN yarn install --frozen-lockfile --no-cache --production

COPY ./client/src/rpc/proto ./src/rpc/proto/
COPY ./server/proto/*.proto ../server/proto/
RUN yarn run proto:generate

COPY ./client/ ./
RUN yarn run webpack --mode=production

FROM nginx AS client
COPY --from=node-build /build/client/dist /dist
COPY ./client/nginx.conf /etc/nginx/nginx.conf
