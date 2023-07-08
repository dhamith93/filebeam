package queue

import (
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
	Key            string
	SizeBytes      int64
	CompletedBytes int64
	StartTime      int64
	EndTime        int64
	IsDownload     bool
}

type Queue struct {
	Items []Transfer
}

func CreateQueue() Queue {
	return Queue{
		Items: []Transfer{},
	}
}

func (q *Queue) Index(hostArr []string, file file.File) int {
	if len(file.Path) == 0 {
		return q.IndexByName(hostArr, file)
	}
	for i, e := range q.Items {
		if e.Ip == hostArr[0] && e.FilePort == hostArr[1] && e.File.Path == file.Path {
			return i
		}
	}
	return -1
}

func (q *Queue) IndexByName(hostArr []string, file file.File) int {
	for i, e := range q.Items {
		if e.Ip == hostArr[0] && e.FilePort == hostArr[1] && e.File.Name == file.Name {
			return i
		}
	}
	return -1
}

func (q *Queue) AddToQueue(host string, key string, file file.File, isDownload bool) {
	hostArr := strings.Split(host, ":")
	q.Items = append(q.Items, Transfer{
		Ip:             hostArr[0],
		FilePort:       hostArr[1],
		File:           file,
		Status:         "pending",
		Key:            key,
		SizeBytes:      file.Size,
		CompletedBytes: 0,
		StartTime:      0,
		EndTime:        0,
		IsDownload:     isDownload,
	})
}

func (q *Queue) Remove(t Transfer) {
	idx := q.Index([]string{t.Ip, t.FilePort}, t.File)
	q.Items[idx] = q.Items[len(q.Items)-1]
	q.Items[len(q.Items)-1] = Transfer{}
	q.Items = q.Items[:len(q.Items)-1]
}

func (q *Queue) UpdateTransferProgress(host string, file file.File, completed int64, status string) {
	hostArr := strings.Split(host, ":")
	idx := q.Index(hostArr, file)
	q.Items[idx].Status = status
	q.Items[idx].CompletedBytes = completed
}

func (q *Queue) UpdateTransferStatus(host string, file file.File, status string) {
	hostArr := strings.Split(host, ":")
	idx := q.Index(hostArr, file)
	q.Items[idx].Status = status
}

func (q *Queue) UpdateTransferStartTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	idx := q.Index(hostArr, file)
	q.Items[idx].StartTime = time.Now().Unix()
}

func (q *Queue) UpdateTransferEndTime(host string, file file.File) {
	hostArr := strings.Split(host, ":")
	idx := q.Index(hostArr, file)
	q.Items[idx].EndTime = time.Now().Unix()
}

func (q *Queue) UpdateFilePortOfTransfer(ip string, oldPort string, newPort string, file file.File) {
	idx := q.Index([]string{ip, oldPort}, file)
	q.Items[idx].FilePort = newPort
}

func (q *Queue) IsTransferStopped(host string, file file.File) bool {
	hostArr := strings.Split(host, ":")
	idx := q.Index(hostArr, file)
	return q.Items[idx].Status != "processing" && q.Items[idx].Status != "pending"
}

func (q *Queue) GetFilePath(host string, file file.File) string {
	hostArr := strings.Split(host, ":")
	idx := -1
	for i, e := range q.Items {
		if e.Ip == hostArr[0] && e.FilePort == hostArr[1] && e.File.Name == file.Name {
			idx = i
			break
		}
	}
	return q.Items[idx].File.Path
}

func (q *Queue) Get(host string, f file.File) Transfer {
	hostArr := strings.Split(host, ":")
	idx := -1
	for i, e := range q.Items {
		if e.Ip == hostArr[0] && e.FilePort == hostArr[1] && e.File.Name == f.Name {
			idx = i
			break
		}
	}
	return q.Items[idx]
}

func (q *Queue) GetPendingTransfers() []Transfer {
	out := []Transfer{}
	for _, e := range q.Items {
		if e.Status == "pending" {
			out = append(out, e)
		}
	}
	return out
}

func (q *Queue) FileTransfersInProgress(count int) bool {
	counter := 0
	for _, item := range q.Items {
		if item.Status == "processing" {
			counter += 1
		}
	}
	return counter > count
}
