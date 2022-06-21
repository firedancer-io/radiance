FROM docker.io/node:18.4.0-bullseye

WORKDIR /work
ENTRYPOINT /bin/bash

RUN npm install -g pnpm
RUN apt-get update && apt-get install -y vim
