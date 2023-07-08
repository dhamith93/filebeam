package fileservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

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

func (f *FileService) ReceiveEncrypted(file file.File) error {
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
		newHost := strings.Split(c.RemoteAddr().String(), ":")
		f.DownloadQueue.UpdateFilePortOfTransfer(newHost[0], "xxxx", newHost[1], file)

		homeDir, _ := os.UserHomeDir()
		fo, err := os.Create(filepath.Join(homeDir, "Downloads", file.Name))
		if err != nil {
			f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
			return err
		}
		defer fo.Close()
		completed := 0

		block, err := aes.NewCipher([]byte(file.Key))
		if err != nil {
			return err
		}

		nonce := make([]byte, block.BlockSize())
		_, err = c.Read(nonce)
		if err != nil {
			return err
		}

		stream := cipher.NewCTR(block, nonce)

		buffer := make([]byte, 4096)

		f.DownloadQueue.UpdateTransferStartTime(c.RemoteAddr().String(), file)

		for {
			if f.DownloadQueue.IsTransferStopped(c.RemoteAddr().String(), file) {
				f.DownloadQueue.UpdateTransferEndTime(c.RemoteAddr().String(), file)
				return fmt.Errorf("download_canceled")
			}
			n, err := c.Read(buffer)
			if err != nil {
				f.DownloadQueue.UpdateTransferProgress(c.RemoteAddr().String(), file, int64(completed), "processing")
				f.DownloadQueue.UpdateTransferEndTime(c.RemoteAddr().String(), file)
				if err != io.EOF {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
					return err
				}
				if int64(completed) >= file.Size {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, "completed")
				} else {
					f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, "cancelled")
				}
				return nil
			}
			completed += n
			stream.XORKeyStream(buffer[:n], buffer[:n])
			if _, err := fo.Write(buffer[:n]); err != nil {
				f.DownloadQueue.UpdateTransferStatus(c.RemoteAddr().String(), file, err.Error())
				return err
			}
			f.DownloadQueue.UpdateTransferProgress(c.RemoteAddr().String(), file, int64(completed), "processing")
		}
	}
}

func (f *FileService) SendEncrypted(host string, file file.File) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return err
	}
	defer conn.Close()

	f.UploadQueue.UpdateTransferStatus(host, file, "processing")
	fileToSend, err := os.Open(file.Path)
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return err
	}
	defer fileToSend.Close()

	block, err := aes.NewCipher([]byte(file.Key))
	if err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return err
	}

	nonce := make([]byte, block.BlockSize())
	if _, err := rand.Read(nonce); err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return err
	}

	if _, err := conn.Write(nonce); err != nil {
		f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
		return err
	}

	stream := cipher.NewCTR(block, nonce)

	writer := &cipher.StreamWriter{
		S: stream,
		W: conn,
	}

	buffer := make([]byte, 4096)
	completed := 0

	f.UploadQueue.UpdateTransferStartTime(host, file)

	for {
		if f.UploadQueue.IsTransferStopped(host, file) {
			f.UploadQueue.UpdateTransferProgress(host, file, int64(completed), "upload_canceled")
			return err
		}
		n, err := fileToSend.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
			return err
		}

		if _, err := writer.Write(buffer[:n]); err != nil {
			f.UploadQueue.UpdateTransferProgress(host, file, 0, err.Error())
			return err
		}

		completed += n
		f.UploadQueue.UpdateTransferProgress(host, file, int64(completed), "processing")
	}

	if !f.UploadQueue.IsTransferStopped(host, file) {
		f.UploadQueue.UpdateTransferProgress(host, file, int64(completed), "completed")
	}
	f.UploadQueue.UpdateTransferEndTime(host, file)

	return nil
}
