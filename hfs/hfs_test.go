package hfs

import (
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func TestGetAttr(t *testing.T) {
	fs := &NewFS{FileSystem: pathfs.NewDefaultFileSystem()}
	attr := []fuse.Attr{
		{Mode: fuse.S_IFREG | 0777, Size: uint64(len("file1.txt"))},
		{Mode: fuse.S_IFREG | 0777, Size: uint64(len("file2.txt"))},
		{Mode: fuse.S_IFREG | 0777, Size: uint64(len("file3.txt"))},
		{Mode: fuse.S_IFDIR | 0777, Size: uint64(len(""))}}
	status := fuse.OK
	statuserr := fuse.ENOENT
	str := []string{
		"file1.txt", "file2.txt", "file3.txt", ""}
	c := &fuse.Context{}

	falsestr := []string{
		"testing", "what", "file1.tx", "file2"}
	for i := 0; i < 4; i++ {
		a, b := fs.GetAttr(falsestr[i], c)
		if nil != a || statuserr != b {
			t.Error("Get Attr Test for False case ", i+1, " failed!")
		}
	}
	for i := 0; i < 4; i++ {
		a, b := fs.GetAttr(str[i], c)
		if attr[i].Mode != a.Mode || attr[i].Size != a.Size || status != b {
			t.Error("GetAttr Test ", i+1, " failed!")
		}
	}
}

func TestOpenDir(t *testing.T) {
	fs := &NewFS{FileSystem: pathfs.NewDefaultFileSystem()}
	expectl := []fuse.DirEntry{
		{Name: "file1.txt", Mode: fuse.S_IFREG},
		{Name: "file2.txt", Mode: fuse.S_IFREG},
		{Name: "file3.txt", Mode: fuse.S_IFREG}}
	str := []string{
		"", "abc", "xyz", "a1b2c3", ".", "cd"}
	statusok := fuse.OK
	statuserr := fuse.ENOENT
	c := &fuse.Context{}
	for i := 0; i < 6; i++ {
		a, b := fs.OpenDir(str[i], c)
		if i == 0 {
			for j := 0; j < 3; j++ {
				if a[j].Name != expectl[j].Name || a[j].Mode != expectl[j].Mode || b != statusok {
					t.Error("Open Dir Test for true case, entry number ", j+1, " failed!")
				}
			}
		} else {
			if a != nil || b != statuserr {
				t.Error("Open Dir Test for false case number ", i+1, " failed!")
			}
		}
	}
}

func TestOpen(t *testing.T) {
	fs := &NewFS{FileSystem: pathfs.NewDefaultFileSystem()}
	statuserr := fuse.ENOENT
	flag := uint32(0)
	c := &fuse.Context{}
	fwdata := []string{
		"f1.txt", "file.txt", "file1", "*.txt", "ls"}
	for j := 0; j < 5; j++ {
		a, b := fs.Open(fwdata[j], flag, c)
		if a != nil || b != statuserr {
			t.Error("Open Test for false case number ", j+1, " failed!")
		}
	}
}
