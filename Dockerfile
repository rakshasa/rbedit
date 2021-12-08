FROM alpine:3.15 AS build-env

WORKDIR /build

ARG TARGET_ARCH

RUN set -xe; \
  apk add --no-cache \
    curl \
    gcc \
    libc6-compat \
    make

RUN set -xe; \
  curl -LSs "https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz" -o go.tar.gz; \
  tar -C /usr/local/ -xzf go.tar.gz; \
  rm -f go.tar.gz

ENV GOPATH=/go
ENV GOOS="${TARGET_ARCH}"
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CGO_ENABLED=0

ENV PATH=${GOPATH}/bin:/usr/local/go/bin/:${PATH}


FROM build-env AS rbedit-builder

COPY ./ ./

RUN go build \
    -o ./rbedit \
    -v \
    -mod=readonly \
    -mod=vendor \
    -ldflags '-s -w -extldflags "-static"'


FROM scratch AS rbedit

COPY --from=rbedit-builder /build/rbedit /

ENTRYPOINT ["/rbedit"]
