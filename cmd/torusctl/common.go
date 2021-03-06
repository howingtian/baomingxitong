package main

import (
	"fmt"
	"os"

	"github.com/alternative-storage/torus"
	"github.com/alternative-storage/torus/distributor"
	"github.com/alternative-storage/torus/internal/flagconfig"

	// Register all the drivers.
	_ "github.com/alternative-storage/torus/metadata/etcd"
	_ "github.com/alternative-storage/torus/storage"

	"github.com/dustin/go-humanize"
)

func die(why string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, why+"\n", args...)
	os.Exit(1)
}

func mustConnectToMDS() torus.MetadataService {
	cfg := flagconfig.BuildConfigFromFlags()
	mds, err := torus.CreateMetadataService("etcd", cfg)
	if err != nil {
		die("couldn't connect to etcd: %v", err)
	}
	return mds
}

func createServer() *torus.Server {
	cfg := flagconfig.BuildConfigFromFlags()
	srv, err := torus.NewServer(cfg, "etcd", "temp")
	if err != nil {
		die("couldn't start: %s", err)
	}
	err = distributor.OpenReplication(srv)
	if err != nil {
		die("couldn't start: %s", err)
	}
	return srv
}

func bytesOrIbytes(s uint64, si bool) string {
	if si {
		return humanize.Bytes(s)
	}
	return humanize.IBytes(s)
}
