package database

import (
	"strconv"
	"time"

	"github.com/dhamith93/SyMon/pkg/memdb"
	"github.com/dhamith93/filebeam/internal/file"
)

type MemDatabase struct {
	Db memdb.Database
}

type Transfer struct {
	Ip             string
	Name           string
	Path           string
	Type           string
	Extension      string
	Status         string
	StartTime      string
	EndTime        string
	SizeBytes      int64
	CompletedBytes int64
	IsDownload     bool
}

func (d *MemDatabase) SetKey(key string) error {
	return d.Db.Tables["meta"].Insert("key", key)
}

func (d *MemDatabase) AddDevice(host string) error {
	return d.Db.Tables["device"].Insert("host", host)
}

func (d *MemDatabase) AddIncomingTransfer(src string, file string, fileType string, extension string, size int64) error {
	return d.Db.Tables["incoming_transfer"].Insert("src, file_name, type, extension, size_bytes, completed_bytes, status, start_time, end_time", src, file, fileType, extension, size, int64(0), "pending", "", "")
}

func (d *MemDatabase) UpdateIncomingTransferProgress(src string, file string, completed int64) error {
	res := d.Db.Tables["incoming_transfer"].Where("file_name", "==", file).And("src", "==", src)

	if res.Error != nil {
		return res.Error
	}

	res.Update("completed_bytes", completed)

	return res.Error
}

func (d *MemDatabase) UpdateIncomingTransferStatus(src string, file string, status string) error {
	res := d.Db.Tables["incoming_transfer"].Where("file_name", "==", file).And("src", "==", src)

	if res.Error != nil {
		return res.Error
	}

	res.Update("status", status)

	return res.Error
}

func (d *MemDatabase) UpdateIncomingTransferStartTime(src string, file string) error {
	res := d.Db.Tables["incoming_transfer"].Where("file_name", "==", file).And("src", "==", src)

	if res.Error != nil {
		return res.Error
	}

	t := strconv.FormatInt(time.Now().Unix(), 10)
	res.Update("start_time", t)

	return res.Error
}

func (d *MemDatabase) UpdateIncomingTransferEndTime(src string, file string) error {
	res := d.Db.Tables["incoming_transfer"].Where("file_name", "==", file).And("src", "==", src)

	if res.Error != nil {
		return res.Error
	}

	t := strconv.FormatInt(time.Now().Unix(), 10)
	res.Update("end_time", t)

	return res.Error
}

func (d *MemDatabase) AddTransfer(dest string, key string, file string, fileType string, extension string, path string, size int64) error {
	return d.Db.Tables["transfer"].Insert("dest, key, file_name, type, extension, file_path, size_bytes, completed_bytes, status, start_time, end_time", dest, key, file, fileType, extension, path, size, int64(0), "pending", "", "")
}

func (d *MemDatabase) UpdateTransferProgress(dest string, file string, completed int64, status string) error {
	res := d.Db.Tables["transfer"].Where("file_path", "==", file).And("dest", "==", dest)

	if res.Error != nil {
		return res.Error
	}

	res.Update("completed_bytes", completed)
	res.Update("status", status)

	return res.Error
}

func (d *MemDatabase) UpdateTransferStatus(dest string, file string, status string) error {
	res := d.Db.Tables["transfer"].Where("file_path", "==", file).And("dest", "==", dest)

	if res.Error != nil {
		return res.Error
	}

	res.Update("status", status)

	return res.Error
}

func (d *MemDatabase) UpdateTransferStartTime(dest string, file string) error {
	res := d.Db.Tables["transfer"].Where("file_path", "==", file).And("dest", "==", dest)

	if res.Error != nil {
		return res.Error
	}

	t := strconv.FormatInt(time.Now().Unix(), 10)
	res.Update("start_time", t)

	return res.Error
}

func (d *MemDatabase) UpdateTransferEndTime(dest string, file string) error {
	res := d.Db.Tables["transfer"].Where("file_path", "==", file).And("dest", "==", dest)

	if res.Error != nil {
		return res.Error
	}

	t := strconv.FormatInt(time.Now().Unix(), 10)
	res.Update("end_time", t)

	return res.Error
}

func (d *MemDatabase) GetDevices() ([]string, error) {
	output := []string{}
	res := d.Db.Tables["device"].Select("*")

	for _, row := range res.Rows {
		output = append(output, row.Columns["host"].StringVal)
	}

	return output, res.Error
}

func (d *MemDatabase) GetAllTransfers() ([]Transfer, error) {
	out := []Transfer{}
	res := d.Db.Tables["transfer"].Select("*")

	for _, row := range res.Rows {
		out = append(out, Transfer{
			IsDownload:     false,
			Name:           row.Columns["file_name"].StringVal,
			Path:           row.Columns["file_path"].StringVal,
			Type:           row.Columns["type"].StringVal,
			Extension:      row.Columns["extension"].StringVal,
			Ip:             row.Columns["dest"].StringVal,
			Status:         row.Columns["status"].StringVal,
			StartTime:      row.Columns["start_time"].StringVal,
			EndTime:        row.Columns["end_time"].StringVal,
			SizeBytes:      row.Columns["size_bytes"].Int64Val,
			CompletedBytes: row.Columns["completed_bytes"].Int64Val,
		})
	}

	res = d.Db.Tables["incoming_transfer"].Select("*")

	for _, row := range res.Rows {
		out = append(out, Transfer{
			IsDownload:     true,
			Name:           row.Columns["file_name"].StringVal,
			Type:           row.Columns["type"].StringVal,
			Extension:      row.Columns["extension"].StringVal,
			Ip:             row.Columns["src"].StringVal,
			Status:         row.Columns["status"].StringVal,
			StartTime:      row.Columns["start_time"].StringVal,
			EndTime:        row.Columns["end_time"].StringVal,
			SizeBytes:      row.Columns["size_bytes"].Int64Val,
			CompletedBytes: row.Columns["completed_bytes"].Int64Val,
		})
	}

	return out, nil
}

func (d *MemDatabase) GetFilePath(dest string, name string) string {
	output := ""
	res := d.Db.Tables["transfer"].Where("file_name", "==", name).And("dest", "==", dest).Select("file_path")
	for _, row := range res.Rows {
		output = row.Columns["file_path"].StringVal
	}
	return output
}

func (d *MemDatabase) GetPendingTransfers() ([]file.File, error) {
	output := []file.File{}
	res := d.Db.Tables["transfer"].Where("status", "==", "pending").Select("dest, key, file_name, file_path, type, extension, size_bytes")
	for _, row := range res.Rows {
		output = append(output, file.File{
			Dest:      row.Columns["dest"].StringVal,
			Key:       row.Columns["key"].StringVal,
			Name:      row.Columns["file_name"].StringVal,
			Path:      row.Columns["file_path"].StringVal,
			Type:      row.Columns["type"].StringVal,
			Extension: row.Columns["extension"].StringVal,
			Size:      row.Columns["size_bytes"].Int64Val,
		})
	}
	return output, res.Error
}

func (d *MemDatabase) FileTransfersInProgress(count int) bool {
	res := d.Db.Tables["transfer"].Where("status", "==", "processing").Select("*")
	return res.RowCount > 0
}

func (d *MemDatabase) IsIncomingTransferStopped(src string, filename string) bool {
	res := d.Db.Tables["incoming_transfer"].Where("src", "==", src).And("file_name", "==", filename).Select("status")
	for _, row := range res.Rows {
		return row.Columns["status"].StringVal == "cancelled"
	}
	return false
}

func (d *MemDatabase) IsTransferStopped(dest string, filepath string) bool {
	res := d.Db.Tables["transfer"].Where("dest", "==", dest).And("file_path", "==", filepath).Select("status")
	for _, row := range res.Rows {
		return row.Columns["status"].StringVal == "cancelled"
	}
	return false
}

func (d *MemDatabase) CreateDB() error {
	d.Db = memdb.CreateDatabase("filebeam")

	err := d.Db.Create(
		"meta",
		memdb.Col{Name: "key", Type: memdb.String},
	)
	if err != nil {
		return err
	}

	err = d.Db.Create(
		"device",
		memdb.Col{Name: "host", Type: memdb.String},
	)
	if err != nil {
		return err
	}

	err = d.Db.Create(
		"transfer",
		memdb.Col{Name: "dest", Type: memdb.String},
		memdb.Col{Name: "key", Type: memdb.String},
		memdb.Col{Name: "file_name", Type: memdb.String},
		memdb.Col{Name: "type", Type: memdb.String},
		memdb.Col{Name: "extension", Type: memdb.String},
		memdb.Col{Name: "file_path", Type: memdb.String},
		memdb.Col{Name: "size_bytes", Type: memdb.Int64},
		memdb.Col{Name: "completed_bytes", Type: memdb.Int64},
		memdb.Col{Name: "status", Type: memdb.String},
		memdb.Col{Name: "start_time", Type: memdb.String},
		memdb.Col{Name: "end_time", Type: memdb.String},
	)
	if err != nil {
		return err
	}

	err = d.Db.Create(
		"incoming_transfer",
		memdb.Col{Name: "src", Type: memdb.String},
		memdb.Col{Name: "file_name", Type: memdb.String},
		memdb.Col{Name: "type", Type: memdb.String},
		memdb.Col{Name: "extension", Type: memdb.String},
		memdb.Col{Name: "size_bytes", Type: memdb.Int64},
		memdb.Col{Name: "completed_bytes", Type: memdb.Int64},
		memdb.Col{Name: "status", Type: memdb.String},
		memdb.Col{Name: "start_time", Type: memdb.String},
		memdb.Col{Name: "end_time", Type: memdb.String},
	)

	return err
}
