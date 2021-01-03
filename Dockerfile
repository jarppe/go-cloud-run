#
# tini:
#

FROM debian:buster-slim AS tini

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

#
# Server:
#

FROM golang:1.15-buster AS server

WORKDIR /app

COPY ./go/go.mod  .
COPY ./go/go.sum  .

RUN go mod download

COPY ./go  ./

RUN CGO_ENABLED=0                      \
    GOOS=linux                         \
    GOARCH=amd64                       \
    go build -a -o server

#
# Front assets:
#

FROM node:15-buster-slim AS assets

WORKDIR /app

COPY ./js/package.json                 \
     ./js/yarn.lock                    \
     ./

RUN yarn install --production=false    \
                 --frozen-lockfile     \
                 --no-progress

COPY ./js/assets              ./assets
COPY ./js/src                 ./src
COPY ./js/webpack.config.js   ./

RUN yarn prod

RUN for f in ./dist/*; do             \
      gzip -9 "$f";                   \
    done

#
# Dist:
#

FROM gcr.io/distroless/static-debian10 AS dist

COPY --from=tini    /tini            /tini
COPY --from=server  /app/server      /app/server
COPY --from=server  /app/templates   /app/templates
COPY --from=assets  /app/dist        /app/assets

ENV HOST=0.0.0.0
ENV PORT=8080
ENV MODE=production
ENV ASSETS=/app/assets

WORKDIR /app
ENTRYPOINT ["/tini", "--"]
CMD ["/app/server"]
