help:
  @just --list


# Build Docker image
build-image:
  docker build -t go-cloud-run:latest .


# Build Docker image
run-image:
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


sh:
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
