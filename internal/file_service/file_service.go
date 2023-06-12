package fileservice

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhamith93/filebeam/internal/file"
	"github.com/dhamith93/filebeam/internal/queue"
)

type ReaderPassThru struct {
	io.Reader
	total int64
	host  string
	path  string
	queue *queue.Queue
}

func (r *ReaderPassThru) Read(p []byte) (int, error) {
	if r.queue.IsTransferStopped(r.host, file.File{Path: r.path}) {
		return 0, fmt.Errorf("upload_canceled")
	}

	n, err := r.Reader.Read(p)
	r.total += int64(n)

	if err != nil {
		return 0, err
	}

	r.queue.UpdateTransferProgress(r.host, file.File{Path: r.path}, int64(r.total), "processing")
	return n, err
}

type FileService struct {
	Port          string
	DownloadQueue *queue.Queue
	UploadQueue   *queue.Queue
}

func (f *FileService) Receive(file file.File) error {
	file.Name = strings.ReplaceAll(file.Name, ":", "_")
	listener, err := net.Listen("tcp", "0.0.0.0:")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	f.Port = strings.Split(listener.Addr().String(), ":")[3]
	log.Println("Server listening on: " + listener.Addr().String())
	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		defer c.Close()
		f.DownloadQueue.AddToQueue(c.RemoteAddr().String(), "", file)

		homeDir, _ := os.UserHomeDir()
		fo, err := os.Create(filepath.Join(homeDir, "Downloads", file.Name))
		if err != nil {
			f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
			return err
		}
		defer fo.Close()
		buf := make([]byte, 1024)
		completed := 0

		// Update progress every second
		ticker := time.NewTicker(time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					f.DownloadQueue.UpdateTransferProgress(c.RemoteAddr().String(), file, int64(completed), "processing")
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		f.DownloadQueue.UpdateTransferStartTime(c.RemoteAddr().String(), file)

		for {
			if f.DownloadQueue.IsTransferStopped(c.RemoteAddr().String(), file) {
				f.DownloadQueue.UpdateTransferEndTime(c.RemoteAddr().String(), file)
				return fmt.Errorf("download_canceled")
			}
			n, err := c.Read(buf)
			if err != nil {
				// stop ticker
				close(quit)
				f.DownloadQueue.UpdateTransferProgress(c.RemoteAddr().String(), file, int64(completed), "processing")
				f.DownloadQueue.UpdateTransferEndTime(c.RemoteAddr().String(), file)
				if err != io.EOF {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
					return err
				}
				if file.Size == int64(completed) {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, "completed")
				} else {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, "cancelled")
				}
				return nil
			}
			completed += n
			if _, err := fo.Write(buf[:n]); err != nil {
				f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
				return err
			}
		}
	}
}

func (f *FileService) Send(host string, file file.File) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		// log.Printf("%s => %s : %s\n", file.Path, host, err.Error())
		log.Fatal(err)
	}
	defer conn.Close()

	f.UploadQueue.UpdateTransferStatus(host, file, "processing")
	fileToSend, err := os.Open(file.Path)
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return
	}

	pr, pw := io.Pipe()
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return
	}

	f.UploadQueue.UpdateTransferStartTime(host, file)

	toSend := &ReaderPassThru{Reader: fileToSend, host: host, queue: f.UploadQueue, path: file.Path}

	go func() {
		_, err := io.Copy(pw, toSend)
		if err != nil {
			log.Println(err)
		}
		pw.Close()
	}()

	n, err := io.Copy(conn, pr)
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return
	}

	if !f.UploadQueue.IsTransferStopped(host, file) {
		f.UploadQueue.UpdateTransferProgress(host, file, int64(n), "completed")
	}
	f.UploadQueue.UpdateTransferEndTime(host, file)
}
