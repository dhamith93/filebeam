package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/dhamith93/filebeam/internal/file"
)

type Transfer struct {
	Id             int
	Ip             string
	FilePort       string
	File           file.File
	Status         string
	SizeBytes      int64
	CompletedBytes int64
	StartTime      int64
	EndTime        int64
}

type Queue struct {
	Upload   []Transfer
	Download []Transfer
}

func CreateQueue() Queue {
	return Queue{
		Download: []Transfer{},
		Upload:   []Transfer{},
	}
}

func (q *Queue) Init() {
	q.Download = []Transfer{}
	q.Upload = []Transfer{}
}

func (q *Queue) AddToQueue(host string, file file.File, isDownload bool) {
	hostArr := strings.Split(host, ":")
	if isDownload {
		q.Download = append(q.Download, Transfer{
			Ip:             hostArr[0],
			FilePort:       hostArr[1],
			Status:         "pending",
			SizeBytes:      file.Size,
			CompletedBytes: 0,
			StartTime:      0,
			EndTime:        0,
		})
		fmt.Println(q.Download)
	} else {
		q.Upload = append(q.Upload, Transfer{
			Ip:             hostArr[0],
			FilePort:       hostArr[1],
			File:           file,
			Status:         "pending",
			SizeBytes:      file.Size,
			CompletedBytes: 0,
			StartTime:      0,
			EndTime:        0,
		})
		fmt.Println(q.Upload)
	}
}

func (q *Queue) UpdateTransferProgress(host string, file file.File, completed int64, status string) {
	hostArr := strings.Split(host, ":")
	for _, upload := range q.Upload {
		if upload.Ip == hostArr[0] && upload.FilePort == hostArr[1] && upload.File.Path == file.Path {
			upload.Status = status
			upload.CompletedBytes = completed
			fmt.Println(upload.CompletedBytes)
		}
	}
}

func (q *Queue) UpdateTransferStartTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	for _, upload := range q.Upload {
		if upload.Ip == hostArr[0] && upload.FilePort == hostArr[1] && upload.File.Path == file.Path {
			upload.StartTime = time.Now().Unix()
			fmt.Println(upload.StartTime)
		}
	}
}

func (q *Queue) UpdateTransferEndTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	for _, upload := range q.Upload {
		if upload.Ip == hostArr[0] && upload.FilePort == hostArr[1] && upload.File.Path == file.Path {
			upload.EndTime = time.Now().Unix()
			fmt.Println(upload.StartTime)
		}
	}
}

func (q *Queue) IsTransferStopped(host string, file file.File) bool {
	fmt.Println(host)
	hostArr := strings.Split(host, ":")
	for _, upload := range q.Upload {
		if upload.Ip == hostArr[0] && upload.FilePort == hostArr[1] && upload.File.Path == file.Path {
			return upload.Status != "processing"
		}
	}
	return true
}

func (q *Queue) UpdateIncomingTransferProgress(host string, file file.File, completed int64) {
	hostArr := strings.Split(host, ":")
	for _, download := range q.Download {
		if download.Ip == hostArr[0] && download.FilePort == hostArr[1] && download.File.Path == file.Path {
			download.CompletedBytes = completed
			fmt.Println(download.CompletedBytes)
		}
	}
}

func (q *Queue) UpdateIncomingTransferStatus(host string, file file.File, status string) {
	fmt.Println(host)
	hostArr := strings.Split(host, ":")
	for _, download := range q.Download {
		if download.Ip == hostArr[0] && download.FilePort == hostArr[1] && download.File.Path == file.Path {
			download.Status = status
			fmt.Println(download.Status)
		}
	}
}

func (q *Queue) UpdateIncomingTransferStartTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	for _, download := range q.Download {
		if download.Ip == hostArr[0] && download.FilePort == hostArr[1] && download.File.Path == file.Path {
			download.StartTime = time.Now().Unix()
			fmt.Println(download.StartTime)
		}
	}
}

func (q *Queue) UpdateIncomingTransferEndTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	for _, download := range q.Download {
		if download.Ip == hostArr[0] && download.FilePort == hostArr[1] && download.File.Path == file.Path {
			download.EndTime = time.Now().Unix()
			fmt.Println(download.StartTime)
		}
	}
}

func (q *Queue) IsIncomingTransferStopped(host string, file file.File) bool {
	fmt.Println(host)
	hostArr := strings.Split(host, ":")
	for _, download := range q.Download {
		if download.Ip == hostArr[0] && download.FilePort == hostArr[1] && download.File.Path == file.Path {
			return download.Status != "processing"
		}
	}
	return true
}

func (q *Queue) GetAllTransfers() []Transfer {
	return append(q.Download, q.Upload...)
}

func (q *Queue) GetPendingTransfers() []Transfer {
	out := []Transfer{}
	for _, upload := range q.Download {
		if upload.Status == "pending" {
			out = append(out, upload)
		}
	}
	return out
}

func (q *Queue) FileTransfersInProgress(count int) bool {
	out := []Transfer{}
	for _, upload := range q.Upload {
		if upload.Status == "processing" {
			out = append(out, upload)
		}
	}
	return len(out) > count
}
