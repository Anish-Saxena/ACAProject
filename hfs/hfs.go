package hfs

import (
	"flag"
	"fmt"
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type NewFS struct {
	pathfs.FileSystem
}

func (fs *NewFS) GetAttr(s string, c *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch s {
	case "file1.txt":
		return &fuse.Attr{Mode: fuse.S_IFREG | 0777, Size: uint64(len(s))}, fuse.OK
	case "file2.txt":
		return &fuse.Attr{Mode: fuse.S_IFREG | 0777, Size: uint64(len(s))}, fuse.OK
	case "file3.txt":
		return &fuse.Attr{Mode: fuse.S_IFREG | 0777, Size: uint64(len(s))}, fuse.OK
	case "":
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0777}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs *NewFS) OpenDir(s string, c *fuse.Context) (list []fuse.DirEntry, code fuse.Status) {
	if s == "" {
		list = []fuse.DirEntry{
			{Name: "file1.txt", Mode: fuse.S_IFREG},
			{Name: "file2.txt", Mode: fuse.S_IFREG},
			{Name: "file3.txt", Mode: fuse.S_IFREG}}
		return list, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs *NewFS) Open(s string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	switch s {
	case "file1.txt":
		return nodefs.NewDataFile([]byte("This is file 1, hello there!")), fuse.OK
	case "file2.txt":
		return nodefs.NewDataFile([]byte("This is file 2, hello again!")), fuse.OK
	case "file3.txt":
		return nodefs.NewDataFile([]byte("This is file 3, hello yet again!")), fuse.OK
	}
	return nil, fuse.ENOENT
}

func BeginServer(mp string) {
	fs := &NewFS{FileSystem: pathfs.NewDefaultFileSystem()}
	pnfs := pathfs.NewPathNodeFs(fs, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), pnfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mountfail: %v\n", err)
	}
	fmt.Println("Server is up and running")
	server.Serve()
}
