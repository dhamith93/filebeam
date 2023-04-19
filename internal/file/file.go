package file

import (
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type File struct {
	Id        int32
	Size      int64
	Name      string
	Type      string
	Extension string
	Path      string
	Dest      string
	Src       string
}

func CreateFile(path string) File {
	file, err := os.Open(path)
	if err != nil {
		return File{}
	}
	defer file.Close()
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return File{}
	}
	stat, _ := file.Stat()

	return File{
		Name:      filepath.Base(file.Name()),
		Size:      stat.Size(),
		Type:      mtype.String(),
		Extension: mtype.Extension(),
		Path:      path,
	}
}

func (f *File) IsFile() bool {
	_, err := os.Open(f.Path)
	return err == nil
}
