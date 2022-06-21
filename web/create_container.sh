#!/usr/bin/env bash
set -euo pipefail

podman build -t radiance-web:latest .

! podman rm -f radiance-web
podman run -it \
    -v $(pwd):$(pwd):z \
    -w $(pwd) \
    --net=host \
    --name radiance-web \
    radiance-web
