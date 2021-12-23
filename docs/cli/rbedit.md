## rbedit

A dependency-free bencode editor

### Synopsis

A statically compiled and dependency-free Bencode editor in Go, useful for command line use and scripts.

To enable bash completions add the following to your ~/.bash_profile file:

```
. <(./build/rbedit-darwin-amd64 completion bash --no-descriptions)
```


## Usage


### Examples

Print single-tracker announce URL:

```bash
$ rbedit announce get --input ./slackware-14.2-install-d1.torrent
http://trackers.transamrit.net:8082/announce
```

Change single-tracker announce URL:

```bash
$ rbedit announce put --input ./slackware-14.2-install-d1.torrent --inplace http://example.com/announce
$ rbedit announce put --input ./slackware-14.2-install-d1.torrent --output ./new.torrent http://example.com/announce
```

Add tracker URL to announce list:

```bash
$ rbedit announce-list get --input ./slackware-14.2-install-d1.torrent
0: http://tracker1.transamrit.net:8082/announce
0: http://tracker2.transamrit.net:8082/announce
0: http://tracker3.transamrit.net:8082/announce

$ rbedit announce-list append-tracker --input ~/Downloads/slackware-14.2-install-d1.torrent --output ./test.torrent 0 http://0.example.com/announce
$ rbedit announce-list get --input ./test.torrent
0: http://tracker1.transamrit.net:8082/announce
0: http://tracker2.transamrit.net:8082/announce
0: http://tracker3.transamrit.net:8082/announce
0: http://0.example.com/announce
```


### Available Commands

```
  announce      BitTorrent announce related commands
  announce-list BitTorrent announce-list related commands
  get           Get object
  hash          SHA256 hashing commands
  put           Put object
  remove        Remove object
```


### Options

```
  -h, --help   help for rbedit
```

### SEE ALSO

* [rbedit announce](rbedit_announce.md)	 - BitTorrent announce related commands
* [rbedit announce-list](rbedit_announce-list.md)	 - BitTorrent announce-list related commands
* [rbedit get](rbedit_get.md)	 - Get object
* [rbedit hash](rbedit_hash.md)	 - SHA256 hashing commands
* [rbedit put](rbedit_put.md)	 - Put object
* [rbedit remove](rbedit_remove.md)	 - Remove object

