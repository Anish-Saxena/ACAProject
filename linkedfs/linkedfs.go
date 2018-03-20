package linkedfs

import (
	"fmt"
	"io"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"syscall"
	"time"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type loopbackFileSystem struct {
	pathfs.FileSystem
	Root string
}

func NewLoopbackFileSystem(root string) pathfs.FileSystem {
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	return &loopbackFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Root:       root,
	}
}


func (fs *loopbackFileSystem) Readlink(name string, context *fuse.Context) (out string, code fuse.Status) {
	f, err := os.Readlink(fs.GetPath(name))
	return f, fuse.ToStatus(err)
}

func (fs *loopbackFileSystem) GetPath(relPath string) string {
	return filepath.Join(fs.Root, relPath)
}

func (fs *loopbackFileSystem) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	fullPath := fs.GetPath(name)
	var err error = nil
	st := syscall.Stat_t{}
	if name == "" {
		err = syscall.Stat(fullPath, &st)
	} else {
		err = syscall.Lstat(fullPath, &st)
	}

	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	a = &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}

func (fs *loopbackFileSystem) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	f, err := os.Open(fs.GetPath(name))
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	want := 500
	output := make([]fuse.DirEntry, 0, want)
	for {
		infos, err := f.Readdir(want)
		for i := range infos {
			if infos[i] == nil {
				continue
			}
			n := infos[i].Name()
			d := fuse.DirEntry{
				Name: n,
			}
			if s := fuse.ToStatT(infos[i]); s != nil {
				d.Mode = uint32(s.Mode)
				d.Ino = s.Ino
			} else {
				log.Printf("ReadDir entry %q for %q has no stat info", n, name)
			}
			output = append(output, d)
		}
		if len(infos) < want || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Readdir() returned err:", err)
			break
		}
	}
	f.Close()

	return output, fuse.OK
}

func (fs *loopbackFileSystem) Open(name string, flags uint32, context *fuse.Context) (fuseFile nodefs.File, status fuse.Status) {

	oc, err:= ioutil.ReadFile(fs.GetPath(name))
	if err!=nil {
		log.Fatal("Error in Reading original file: %v",err)
	}
	s:=[]byte("\nAutomatically Added!\n")
	c := append(oc,s...)
	f := nodefs.NewDataFile(c)
	_, err1:=os.Open(fs.GetPath(name))
	if err1 != nil {
		return nil, fuse.ToStatus(err)
	}
	return f, fuse.OK
}  

func Begin(orig string, link string) {

	loopbackfs := NewLoopbackFileSystem(orig)
	finalFs := loopbackfs
	opts := &nodefs.Options{
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	}
	pathFsOpts := &pathfs.PathNodeFsOptions{ClientInodes: false}
	pathFs := pathfs.NewPathNodeFs(finalFs, pathFsOpts)
	conn := nodefs.NewFileSystemConnector(pathFs.Root(), opts)
	mountPoint := link
	origAbs, _ := filepath.Abs(orig)
	mOpts := &fuse.MountOptions{Name: "linkedfs", FsName: origAbs}
	state, err := fuse.NewServer(conn.RawFS(), mountPoint, mOpts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	fmt.Println("Mounted!")
	state.Serve()
}
