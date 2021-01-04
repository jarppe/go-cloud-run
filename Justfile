help:
  @just --list


# Build image
build-server:
  docker build -t go-cloud-run:latest .


# Run image
run-server:
  docker run --rm -p 8080:8080 go-cloud-run:latest


dev +args="":
  docker run                                         \
    --rm                                             \
    --name dev                                       \
    -p 8080:8080                                     \
    -e HOST=0.0.0.0                                  \
    -e PORT=8080                                     \
    -e MODE=develop                                  \
    -e ASSETS=/app/assets                            \
    -w /app/src                                      \
    -v $(pwd)/src:/app/src:cached                    \
    -v $(pwd)/assets:/app/assets:cached              \
    go-cloud-run:dev {{ args }}


sh +args="":
  docker run                                         \
    --rm                                             \
    -it                                              \
    -p 8081:8080                                     \
    -e HOST=0.0.0.0                                  \
    -e PORT=8080                                     \
    -e MODE=develop                                  \
    -e ASSETS=/app/assets                            \
    -w /app/src                                      \
    -v $(pwd)/src:/app/src:cached                    \
    -v $(pwd)/assets:/app/assets:cached              \
    go-cloud-run:dev {{ args }}


build-assets-image:
  docker volume create go-cloud-run-assets 2> /dev/null
  docker build -t go-cloud-run:assets                \
                  -f ./js/Dockerfile-dev             \
                  ./js


assets +args="":
  @docker run                                        \
    --rm                                             \
    --init                                           \
    -it                                              \
    -w /app                                          \
    -v $(pwd)/js:/app:cached                         \
    -v go-cloud-run-assets:/app/dist                 \
    go-cloud-run:assets {{ args }}
