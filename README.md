# Ringify

A small utility to take in a Cassandra Ring File and figure out which nodes are
responsible for a particular partition key (a local version of 
`nodetool getendpoints`)

## Usage

Install the tool using `go get`

```
$ go get github.com/suhailpatel/ringify
```

Generate a ring file from nodetool and query
```
$ nodetool ring > ringfile.txt
$ ringify ringfile.txt partition_abcd1234
Wanted token: 1589072005239795781
10.0.0.1 (rack: eu-west-1a, token: 1589297164579232499)
10.0.0.3 (rack: eu-west-1c, token: 1589797902066161209)
10.0.0.2 (rack: eu-west-1b, token: 1592399308380267011)
```

You can also query for partition keys spanning multiple columns using `,` as
a separator
```
$ ringify ringfile.txt partition_abcd1234,pos_1234
Wanted token: 3942213685374363656
10.0.0.3 (rack: eu-west-1c, token: 3943089618890493298)
10.0.0.2 (rack: eu-west-1b, token: 3943666626842003849)
10.0.0.1 (rack: eu-west-1a, token: 3948204454273716332)
```

## TODO

- This script does not handle token range wrap-around
- It would be nice if the number of racks was configurable
