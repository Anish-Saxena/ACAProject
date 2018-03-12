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
