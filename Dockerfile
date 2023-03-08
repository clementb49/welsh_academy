# syntax=docker/dockerfile:1.5
# Build the go source file in a specific image
FROM golang:1.20-bullseye AS go_build
WORKDIR /usr/src/app
# Use this line to allow apt to cache downloaded package
RUN <<eot sh
    set -e
    rm -f /etc/apt/apt.conf.d/docker-clean
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
eot
# Cache the apt downloaded package and index
# Use gosu to change the user info at runtime
RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt <<eot bash
  set -e
  apt-get update
  apt-get --no-install-recommends install -y gosu
eot
# change entrypoint to use a dynamic user based on user info specified in env file (required only for dev)
COPY --chown=root:root .devcontainer/docker-entrypoint.sh /usr/local/bin/entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint"]
# build the go source code
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build <<eot bash
  set -e
  CGO_ENABLED=0 go build -v
eot
# production image 
FROM gcr.io/distroless/static-debian11:nonroot AS app
# get the binary produced previously
COPY --from=go_build --chown=nonroot:nonroot /usr/src/app/welsh_academy /app
# launch the app
CMD ["/app"]
