package fileservice

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dhamith93/share_core/internal/file"
)

type FileService struct {
	Port string
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
		for {
			n, err := c.Read(buf)
			if err != nil {
				if err != io.EOF {
					return err
				}
				return nil
			}
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
