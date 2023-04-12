package fileservice

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/dhamith93/share_core/internal/database"
	"github.com/dhamith93/share_core/internal/file"
)

type FileService struct {
	Port     string
	Database *database.Database
}

func (f *FileService) Receive(file file.File) error {
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
		fo, err := os.Create(file.Name)
		if err != nil {
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
					f.Database.UpdateIncomingTransfer(strings.Split(c.RemoteAddr().String(), ":")[0], file.Name, int64(completed))
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		for {
			n, err := c.Read(buf)
			if err != nil {
				// stop ticker
				close(quit)
				f.Database.UpdateIncomingTransfer(strings.Split(c.RemoteAddr().String(), ":")[0], file.Name, int64(completed))
				if err != io.EOF {
					return err
				}
				return nil
			}
			completed += n
			if _, err := fo.Write(buf[:n]); err != nil {
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
	log.Println("Connected to server.")

	fileToSend, err := os.Open(file.Path)
	if err != nil {
		log.Fatal(err)
	}

	pr, pw := io.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		_, err := io.Copy(pw, fileToSend)
		if err != nil {
			log.Fatal(err)
		}
		pw.Close()
	}()

	n, err := io.Copy(conn, pr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("copied to connection: %d", n)

	conn.Close()
}
