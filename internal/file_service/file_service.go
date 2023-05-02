package fileservice

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/dhamith93/share_core/internal/database"
	"github.com/dhamith93/share_core/internal/file"
)

type ReaderPassThru struct {
	io.Reader
	total    int64
	ip       string
	path     string
	database *database.Database
}

func (r *ReaderPassThru) Read(p []byte) (int, error) {
	if r.database.IsTransferStopped(r.ip, r.path) {
		return 0, fmt.Errorf("upload_canceled")
	}

	n, err := r.Reader.Read(p)
	r.total += int64(n)

	if err != nil {
		return 0, err
	}

	r.database.UpdateTransferProgress(r.ip, r.path, int64(r.total), "processing")
	return n, err
}

type FileService struct {
	Port     string
	Database *database.Database
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

		fo, err := os.Create(file.Name)
		if err != nil {
			f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_create_file")
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
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		f.Database.UpdateIncomingTransferStartTime(ip, file.Name)

		for {
			if f.Database.IsIncomingTransferStopped(ip, file.Name) {
				f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cancelled")
				return fmt.Errorf("download_canceled")
			}
			n, err := c.Read(buf)
			if err != nil {
				// stop ticker
				close(quit)
				f.Database.UpdateIncomingTransferProgress(ip, file.Name, int64(completed))
				f.Database.UpdateIncomingTransferEndTime(ip, file.Name)
				if err != io.EOF {
					f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_read_incoming_file")
					return err
				}
				f.Database.UpdateIncomingTransferStatus(ip, file.Name, "completed")
				return nil
			}
			completed += n
			if _, err := fo.Write(buf[:n]); err != nil {
				f.Database.UpdateIncomingTransferStatus(ip, file.Name, "cannot_write_to_file")
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
		return
	}

	pr, pw := io.Pipe()
	if err != nil {
		f.Database.UpdateTransferProgress(ip, file.Path, 0, "cannot_read_file")
		return
	}

	f.Database.UpdateTransferStartTime(ip, file.Path)

	toSend := &ReaderPassThru{Reader: fileToSend, database: f.Database, ip: ip, path: file.Path}

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
		return
	}

	if !f.Database.IsTransferStopped(ip, file.Path) {
		f.Database.UpdateTransferProgress(ip, file.Path, int64(n), "completed")
	}
	f.Database.UpdateTransferEndTime(ip, file.Path)
}
