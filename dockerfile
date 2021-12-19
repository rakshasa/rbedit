ARG ALPINE_VERSION=3.15


FROM alpine:${ALPINE_VERSION} AS build-env

ARG GO_VERSION=1.17.2

ARG BUILD_OS=linux
ARG BUILD_ARCH=amd64
ARG TARGET_OS=linux
ARG TARGET_ARCH=amd64

WORKDIR /build

RUN set -eux; \
  apk add --no-cache \
    bash \
    libc6-compat

RUN set -eux; \
  wget -O go.tar.gz "https://dl.google.com/go/go${GO_VERSION}.${BUILD_OS}-${BUILD_ARCH}.tar.gz" ; \
  tar -C /usr/local/ -xzf go.tar.gz; \
  rm -f go.tar.gz

ENV GOPATH=/go
ENV GOOS="${TARGET_OS}"
ENV GOARCH="${TARGET_ARCH}"
ENV GOFLAGS="-v -mod=readonly -mod=vendor"
ENV GO111MODULE=on
ENV CGO_ENABLED=0

ENV PATH="${GOPATH}/bin:/usr/local/go/bin:${PATH}"

RUN go version


FROM build-env AS rbedit-builder

ARG TARGET_OS=linux
ARG TARGET_ARCH=amd64
ARG BUILD_MARKDOWN=no

COPY ./ ./

RUN set -eux; \
  go build -ldflags "-s -w -extldflags '-static  -fno-PIC'" -o "/rbedit-${TARGET_OS}-${TARGET_ARCH}" ./cmd/rbedit; \
  \
  if [ "${BUILD_MARKDOWN}" == "yes" ]; then \
    go build -ldflags "-s -w -extldflags '-static  -fno-PIC'" -o "/rbedit-markdown-${TARGET_OS}-${TARGET_ARCH}" ./cmd/rbedit-markdown; \
  fi


FROM scratch AS rbedit

ARG TARGET_OS=linux
ARG TARGET_ARCH=amd64

COPY --from=rbedit-builder "/rbedit-${TARGET_OS}-${TARGET_ARCH}" /rbedit

ENTRYPOINT ["/rbedit"]
