rbEdit
======

A statically compiled and dependency-free Bencode editor in Go, useful for command line use and scripts.


Quick Start
-----------

```bash
# Compile:
./scripts/build.sh

# Compile for darwin arch:
TARGET_OS=darwin ./scripts/build.sh

# Compile for linux arch:
TARGET_OS=linux ./scripts/build.sh
```

The only build requirement is Docker and the binary has no runtime dependencies.

```bash
./build/rbedit announce get --input ./test.torrent
```

Print announce url.

```bash
./build/rbedit announce put --input ./test.torrent --inplace http://example.com/announce
```

Change announce url.

```bash
./build/rbedit get --input ./test.torrent info length
```

Get value of the `length` entry in the `info` map.

```bash
./build/rbedit put --input ./test.torrent --inplace --bencode d3:bari2e3:bazi3e3:fooi1ee foo-info
```

Write a custom bencoded object to `foo-info` entry in the torrent root.


Batch Operations
----------------

Generate 10,000 torrents with unique info hashes:

```bash
$ RBEDIT_PATH=./build/rbedit-darwin-amd64 COUNT=10000 ./scripts/generate-torrents.sh /tmp/slackware-14.2-install-d1.torrent /tmp/slackware-torrents
Generating torrent test files

RBEDIT_PATH: ./build/rbedit
COUNT: 10000
SRC-TORRENT: /tmp/slackware-14.2-install-d1.torrent
DEST-DIR: /tmp/slackware-torrents
PREFIX-DEPTH: 1

generating..................................

Finished generating 10000 torrents

$ find /tmp/slackware-torrents -type f | wc -l
   10000

# du -h -d0 /tmp/slackware-torrents
273M    /Users/rakshasa/tmp/slackware-torrents/
```

Check the announce url of all 10,000 torrents:

```bash
$ time ./build/rbedit-darwin-amd64 announce get --input <(find /tmp/torrents -type f) --batch | tail -n3
http://trackers.transamrit.net:8082/announce
http://trackers.transamrit.net:8082/announce
http://trackers.transamrit.net:8082/announce

real    0m0.523s
user    0m0.375s
sys     0m0.374s
```

Replace the announce urls and overwrite the source files for all 10,000 torrents:

```bash
$ time ./build/rbedit-darwin-amd64 announce put --input <(find /tmp/torrents -type f) --batch --inplace http://new.example.com/announce

real    0m2.672s
user    0m0.841s
sys     0m1.858s
```

Verify the announce urls were changed:

```bash
$ ./build/rbedit-darwin-amd64 announce get --input <(find /tmp/torrents -type f) --batch | tail -n3
http://new.example.com/announce
http://new.example.com/announce
http://new.example.com/announce
```


Donate to rTorrent development
------------------------------

 * [Paypal](https://paypal.me/jarisundelljp)
 * [Patreon](https://www.patreon.com/rtorrent)
 * [SubscribeStar](https://www.subscribestar.com/rtorrent)
 * Bitcoin: 1MpmXm5AHtdBoDaLZstJw8nupJJaeKu8V8
 * Etherium: 0x9AB1e3C3d8a875e870f161b3e9287Db0E6DAfF78
 * Litecoin: LdyaVR67LBnTf6mAT4QJnjSG2Zk67qxmfQ
 * Cardano: addr1qytaslmqmk6dspltw06sp0zf83dh09u79j49ceh5y26zdcccgq4ph7nmx6kgmzeldauj43254ey97f3x4xw49d86aguqwfhlte
