package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Please enter 2 arguments: linked and original directory, respectively. Exiting now.")
	}
	orig := flag.Arg(1)
	loopbackfs := pathfs.NewLoopbackFileSystem(orig)
	finalFs := loopbackfs
	opts := &nodefs.Options{
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	}
	pathFsOpts := &pathfs.PathNodeFsOptions{ClientInodes: false}
	pathFs := pathfs.NewPathNodeFs(finalFs, pathFsOpts)
	conn := nodefs.NewFileSystemConnector(pathFs.Root(), opts)
	mountPoint := flag.Arg(0)
	origAbs, _ := filepath.Abs(orig)
	mOpts := &fuse.MountOptions{Name:"linkedfs",FsName:origAbs}
	state, err := fuse.NewServer(conn.RawFS(), mountPoint, mOpts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	fmt.Println("Mounted!")
	state.Serve()
}
