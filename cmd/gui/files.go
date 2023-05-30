package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified int64
	IsDir        bool
	Selected     bool
}

type ByTypeAndName []File

func (f ByTypeAndName) Len() int {
	return len(f)
}

func (f ByTypeAndName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f ByTypeAndName) Less(i, j int) bool {
	if f[i].IsDir && !f[j].IsDir {
		return true
	}
	if f[i].IsDir && f[j].IsDir {
		return f[i].Name < f[j].Name
	}
	if !f[i].IsDir && f[j].IsDir {
		return false
	}
	if !f[i].IsDir && !f[j].IsDir {
		return f[i].Name < f[j].Name
	}
	return false
}

func getDirectoryContent(path string) ([]File, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := []File{}
	for _, v := range files {
		out = append(out, File{
			Name:         v.Name(),
			Path:         filepath.Join(path, v.Name()),
			Size:         v.Size(),
			IsDir:        v.IsDir(),
			LastModified: v.ModTime().Unix(),
			Selected:     false,
		})
	}

	sort.Sort(ByTypeAndName(out))

	return out, nil
}

func getHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}
