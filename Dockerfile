FROM ubuntu:focal AS build-env

WORKDIR /build

ARG TARGET_ARCH

RUN set -xe; \
  apt-get update; apt-get upgrade -y; apt-get --purge autoremove -y; apt-get install -y \
    curl \
    gcc \
    linux-libc-dev \
    make; \
  apt-get clean

RUN set -xe; \
  curl -LSs "https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz" -o go.tar.gz; \
  tar -C /usr/local/ -xzf go.tar.gz; \
  rm -f go.tar.gz

ENV PATH=${GOPATH}/bin:/usr/local/go/bin/:${PATH}

ENV GOPATH=/build/go
ENV GOOS="${TARGET_ARCH}"
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CGO_ENABLED=0


FROM build-env AS rbedit-builder

COPY ./ ./

RUN go build \
    -o ./rbedit \
    -v \
    -mod=readonly \
    -mod=vendor \
    -ldflags "-s -w"


FROM scratch AS rbedit

COPY --from=rbedit-builder /build/rbedit /

ENTRYPOINT ["/rbedit"]
