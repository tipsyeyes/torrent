// Converts magnet URIs and info hashes into torrent metainfo files.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/tagflag"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

func main() {
	args := struct {
		tagflag.StartPos
		Magnet []string
	}{}
	tagflag.Parse(&args)
	cl, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "*** pid=%d, ppid=%d\n", os.Getpid(), os.Getppid())
	})
	http.HandleFunc("/torrent", func(w http.ResponseWriter, r *http.Request) {
		cl.WriteStatus(w)
	})
	http.HandleFunc("/dht", func(w http.ResponseWriter, r *http.Request) {
		for _, ds := range cl.DhtServers() {
			ds.WriteStatus(w)
		}
	})
	// add by elias
	go func() {
		log.Println("listen on localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
	wg := sync.WaitGroup{}
	for _, arg := range args.Magnet {
		t, err := cl.AddMagnet(arg)
		if err != nil {
			log.Fatalf("error adding magnet to client: %s", err)
		}
		wg.Add(1)
		go func() {
			log.Println("add magnet")
			defer wg.Done()
			<-t.GotInfo()
			log.Println("add magnet done")
			mi := t.Metainfo()
			t.Drop()
			f, err := os.Create(t.Info().Name + ".torrent")
			if err != nil {
				log.Fatalf("error creating torrent metainfo file: %s", err)
			}
			defer f.Close()
			err = bencode.NewEncoder(f).Encode(mi)
			if err != nil {
				log.Fatalf("error writing torrent metainfo file: %s", err)
			}
		}()
	}
	wg.Wait()
}
