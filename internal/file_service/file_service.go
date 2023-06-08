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

	"github.com/dhamith93/filebeam/internal/database"
	"github.com/dhamith93/filebeam/internal/file"
	"github.com/dhamith93/filebeam/internal/queue"
)

type ReaderPassThru struct {
	io.Reader
	total    int64
	ip       string
	host     string
	path     string
	queue    *queue.Queue
	database *database.MemDatabase
}

func (r *ReaderPassThru) Read(p []byte) (int, error) {
	if r.queue.IsTransferStopped(r.host, file.File{Path: r.path}) {
		return 0, fmt.Errorf("q: upload_canceled")
	}

	if r.database.IsTransferStopped(r.ip, r.path) {
		return 0, fmt.Errorf("upload_canceled")
	}

	n, err := r.Reader.Read(p)
	r.total += int64(n)

	if err != nil {
		return 0, err
	}

	r.database.UpdateTransferProgress(r.ip, r.path, int64(r.total), "processing")
	r.queue.UpdateTransferProgress(r.host, file.File{Path: r.path}, int64(r.total), "processing")
	return n, err
}

type FileService struct {
	Port     string
	Queue    *queue.Queue
	Database *database.MemDatabase
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
		ip := strings.Split(c.RemoteAddr().String(), ":")[0]
		f.Queue.AddToQueue(c.RemoteAddr().String(), file, true)

		homeDir, _ := os.UserHomeDir()
		fo, err := os.Create(filepath.Join(homeDir, "Downloads", file.Name))
		if err != nil {
			f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_create_file")
			f.Queue.UpdateIncomingTransferStatus(c.RemoteAddr().String(), file, err.Error())
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
					f.Database.UpdateIncomingTransferProgress(ip, file.Name, int64(completed))
					f.Queue.UpdateIncomingTransferProgress(c.RemoteAddr().String(), file, int64(completed))
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		f.Database.UpdateIncomingTransferStartTime(ip, file.Name)
		f.Queue.UpdateIncomingTransferStartTime(c.RemoteAddr().String(), file)

		for {
			if f.Database.IsIncomingTransferStopped(ip, file.Name) {
				f.Database.UpdateIncomingTransferEndTime(ip, file.Name)
				f.Queue.UpdateIncomingTransferEndTime(c.RemoteAddr().String(), file)
				return fmt.Errorf("download_canceled")
			}
			n, err := c.Read(buf)
			if err != nil {
				// stop ticker
				close(quit)
				f.Database.UpdateIncomingTransferProgress(ip, file.Name, int64(completed))
				f.Database.UpdateIncomingTransferEndTime(ip, file.Name)
				f.Queue.UpdateIncomingTransferProgress(c.RemoteAddr().String(), file, int64(completed))
				f.Queue.UpdateIncomingTransferEndTime(c.RemoteAddr().String(), file)
				if err != io.EOF {
					f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_read_incoming_file")
					f.Queue.UpdateIncomingTransferStatus(c.RemoteAddr().String(), file, err.Error())
					return err
				}
				if file.Size == int64(completed) {
					f.Database.UpdateIncomingTransferStatus(ip, file.Name, "completed")
					f.Queue.UpdateIncomingTransferStatus(c.RemoteAddr().String(), file, "completed")
				} else {
					f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cancelled")
					f.Queue.UpdateIncomingTransferStatus(c.RemoteAddr().String(), file, "cancelled")
				}
				return nil
			}
			completed += n
			if _, err := fo.Write(buf[:n]); err != nil {
				f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_write_to_file")
				f.Queue.UpdateIncomingTransferStatus(c.RemoteAddr().String(), file, err.Error())
				return err
			}
		}
	}
}

func (f *FileService) Send(host string, file file.File) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ip := strings.Split(host, ":")[0]
	fileToSend, err := os.Open(file.Path)
	if err != nil {
		f.Database.UpdateTransferProgress(ip, file.Path, 0, "cannot_read_file")
		f.Queue.UpdateTransferProgress(host, file, 0, err.Error())
		return
	}

	pr, pw := io.Pipe()
	if err != nil {
		f.Database.UpdateTransferProgress(ip, file.Path, 0, "cannot_read_file")
		f.Queue.UpdateTransferProgress(host, file, 0, err.Error())
		return
	}

	f.Database.UpdateTransferStartTime(ip, file.Path)
	f.Queue.UpdateTransferStartTime(host, file)

	toSend := &ReaderPassThru{Reader: fileToSend, database: f.Database, ip: ip, host: host, path: file.Path}

	go func() {
		_, err := io.Copy(pw, toSend)
		if err != nil {
			log.Println(err)
		}
		pw.Close()
	}()

	n, err := io.Copy(conn, pr)
	if err != nil {
		f.Database.UpdateTransferProgress(ip, file.Path, 0, "cannot_read_file")
		f.Queue.UpdateTransferProgress(ip, file, 0, err.Error())
		return
	}

	if !f.Database.IsTransferStopped(ip, file.Path) {
		f.Database.UpdateTransferProgress(ip, file.Path, int64(n), "completed")
		f.Queue.UpdateTransferProgress(host, file, int64(n), "completed")
	}
	f.Database.UpdateTransferEndTime(ip, file.Path)
	f.Queue.UpdateTransferEndTime(host, file)
}
