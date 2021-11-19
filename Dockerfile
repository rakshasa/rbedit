FROM ubuntu:focal AS build-env

WORKDIR /build

RUN set -xe; \
  apt-get update; apt-get upgrade -y; apt-get --purge autoremove -y; apt-get install -y \
    curl \
    gcc \
    linux-libc-dev \
    make; \
  apt-get clean

RUN set -xe; \
  curl -LSs https://www.musl-libc.org/releases/musl-latest.tar.gz -o musl.tar.gz; \
  tar -xzf musl.tar.gz; \
  ( cd musl-*/; \
    ./configure; \
    make; \
    make install; \
  ); \
  rm -rf musl*; \
  \
  cd /usr/local/musl/include; \
  ln -s /usr/include/linux .; \
  ln -s /usr/include/asm-generic .; \
  ln -s /usr/include/x86_64-linux-gnu .; \
  ln -s x86_64-linux-gnu/asm asm

ENV CC=/usr/local/musl/bin/musl-gcc

RUN set -xe; \
  curl -LSs https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz -o go.tar.gz; \
  tar -C /usr/local/ -xzf go.tar.gz; \
  rm -f go.tar.gz

ENV GOPATH=/build/go
ENV GO111MODULE=on
ENV PATH=${GOPATH}/bin:/usr/local/go/bin/:${PATH}


FROM build-env AS rbedit-builder

COPY ./ ./

RUN go build

RUN ls ./


FROM rbedit-builder AS rbedit

COPY --from=rbedit-builder /build/rbedit /

ENTRYPOINT ["/rbedit"]
