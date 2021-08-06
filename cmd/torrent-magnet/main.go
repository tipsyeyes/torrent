package main

import (
	"fmt"
	"os"

	"github.com/anacrolix/tagflag"

	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	tagflag.Parse(nil, tagflag.Description("reads a torrent file from stdin and writes out its magnet link to stdout"))
	// 创建种子磁力链接
	// torrent-magnet < /Users/elias/sre/deploy/github.com/demo/go/src/regal/test/ax.torrent
	mi, err := metainfo.Load(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading metainfo from stdin: %s", err)
		os.Exit(1)
	}
	info, err := mi.UnmarshalInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling info: %s", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", mi.Magnet(nil, &info).String())
}
