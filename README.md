# Ringify

A small utility to take in a Cassandra Ring File and figure out which nodes are
responsible for a particular partition key.

## Usage

Install the tool using `go get`

```
$ go get github.com/suhailpatel/ringify
```

Generate a ring file from nodetool and query
```
$ nodetool ring > ringfile.txt
$ ringify ringfile.txt partition_abcd1234
```

You can also query for partition keys spanning multiple columns using `:` as
a separator
```
$ ringify ringfile.txt partition_abcd1234:pos_1234
```
