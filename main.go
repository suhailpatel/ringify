package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ringToken struct {
	address string
	rack    string
	token   int64
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: <ringfile> <key>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open ring file: %v", err)
		os.Exit(1)

	}

	scanner := bufio.NewScanner(f)
	ringTokens := make([]ringToken, 0, 50*256)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "10.") {
			continue
		}

		text := regexp.MustCompile("[ ][ ]+").ReplaceAllString(scanner.Text(), " ")
		items := strings.Split(text, " ")

		tokenField := items[len(items)-1]
		if tokenField == "" {
			tokenField = items[len(items)-2]
		}
		token, err := strconv.ParseInt(tokenField, 10, 0)
		if err != nil {
			continue
		}

		ringTokens = append(ringTokens, ringToken{
			address: items[0],
			rack:    items[1],
			token:   token,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to scan ring file: %v", err)
		os.Exit(1)
	}

	sort.Slice(ringTokens, func(i, j int) bool {
		return ringTokens[i].token < ringTokens[j].token
	})

	hashable := hashableBytes([]byte(os.Args[2]))
	wantedToken := Murmur3H1(hashable)
	fmt.Printf("Wanted token: %v\n", wantedToken)

	const maxRacks = 3
	seenRacks := map[string]struct{}{}
	for _, rt := range ringTokens {
		if len(seenRacks) >= maxRacks {
			break
		}

		if _, ok := seenRacks[rt.rack]; ok {
			continue
		}

		if rt.token < wantedToken {
			continue
		}

		seenRacks[rt.rack] = struct{}{}
		fmt.Printf("%v (rack: %v, token: %v)\n", rt.address, rt.rack, rt.token)
	}
}

func hashableBytes(data []byte) []byte {
	split := bytes.Split(data, []byte(","))
	if len(split) == 1 {
		return data
	}

	// If you have a composite key, Cassandra does something funky with the
	// token marhsaling, it expects two bytes of length, followed up by the
	// bytes representing that split followed by a zero byte
	// Source: https://github.com/apache/cassandra/blob/cassandra-3.0.20/src/java/org/apache/cassandra/db/marshal/CompositeType.java#L363-L369
	//
	// So say you had a 2 level composite key hello:suhail
	//  	hello -> 68 65 6c 6c 6f
	// 		suhail -> 73 75 68 61 69 6c
	//
	//	This converts to
	//  	hello -> 00 05 68 65 6c 6c 6f 00
	// 		suhail -> 00 06 73 75 68 61 69 6c 00
	//
	//
	out := make([]byte, 0, len(split))
	for _, item := range split {
		out = append(out, byte((len(item)>>8)&0xFF))
		out = append(out, byte(len(item)&0xFF))
		out = append(out, item...)
		out = append(out, byte(0))
	}
	return out
}
