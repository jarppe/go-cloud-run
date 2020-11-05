help:
  @just --list


# Build Docker image
build-image:
  docker build -t talsu-server:latest .


# Push image to Google repo
push-image:
  docker tag talsu-server:latest eu.gcr.io/talsu-291313/server:latest
  docker push eu.gcr.io/talsu-291313/server:latest
