package main

import (
	"fmt"
	"log"

	"github.com/anacrolix/tagflag"

	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	var args struct {
		tagflag.StartPos
		Files []string `arity:"+" type:"pos"`
	}
	tagflag.Parse(&args)
	// 计算种子hash
	// torrent-infohash /Users/elias/sre/deploy/github.com/demo/go/src/regal/test/torrents/ax.torrent
	for _, arg := range args.Files {
		mi, err := metainfo.LoadFromFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", mi.HashInfoBytes().HexString(), arg)
	}
}
