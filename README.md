rbEdit
======

A statically compiled and dependency-free Bencode editor in Go, useful for command line use and scripts.


Quick Start
-----------

```bash
# Compile for linux arch:
./scripts/build.sh

# Compile for darwin arch:
RBEDIT_ARCH=darwin ./scripts/build.sh
```

The only dependency is Docker.

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


Donate to rTorrent development
------------------------------

 * [Paypal](https://paypal.me/jarisundelljp)
 * [Patreon](https://www.patreon.com/rtorrent)
 * [SubscribeStar](https://www.subscribestar.com/rtorrent)
 * Bitcoin: 1MpmXm5AHtdBoDaLZstJw8nupJJaeKu8V8
 * Etherium: 0x9AB1e3C3d8a875e870f161b3e9287Db0E6DAfF78
 * Litecoin: LdyaVR67LBnTf6mAT4QJnjSG2Zk67qxmfQ
 * Cardano: addr1qytaslmqmk6dspltw06sp0zf83dh09u79j49ceh5y26zdcccgq4ph7nmx6kgmzeldauj43254ey97f3x4xw49d86aguqwfhlte
