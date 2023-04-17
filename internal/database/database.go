package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/dhamith93/share_core/internal/file"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

type IncomingTransfer struct {
	Src            string
	File           string
	SizeBytes      int64
	CompletedBytes int64
}

func (d *Database) CreateDB(dbName string) {
	file, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Println(err)
	}
	d.Db = db
	d.initDB()
}

func (d *Database) AddDevice(host string) error {
	return d.execute(
		"INSERT INTO device (host) VALUES (?);",
		host,
	)
}

func (d *Database) AddIncomingTransfer(src string, file string, size int64) error {
	return d.execute(
		"INSERT INTO incoming_transfer (src, file_name, size_bytes, completed_bytes) VALUES (?, ?, ?, ?);",
		src, file, size, 0,
	)
}

func (d *Database) UpdateIncomingTransferProgress(src string, file string, completed int64) error {
	return d.execute(
		"UPDATE incoming_transfer SET completed_bytes = ? WHERE file_name = ? AND src = ?",
		completed, file, src,
	)
}

func (d *Database) UpdateTransferProgress(dest string, file string, completed int64, status string) error {
	dest = dest + ":%"
	return d.execute(
		"UPDATE transfer SET completed_bytes = ?, status = ? WHERE file_path = ? AND dest LIKE '"+dest+"';",
		completed, status, file,
	)
}

func (d *Database) UpdateTransferStatus(dest string, file string, status string) error {
	return d.execute(
		"UPDATE transfer SET status = ? WHERE file_path = ? AND dest = ?",
		status, file, dest,
	)
}

func (d *Database) execute(query string, args ...interface{}) error {
	stmt, err := d.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetDevices() ([]string, error) {
	output := []string{}
	rows, err := d.Db.Query("SELECT * FROM device;")
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var device string

		err = rows.Scan(&device)

		if err != nil {
			return output, err
		}

		output = append(output, device)
	}

	return output, nil
}

func (d *Database) GetIncomingTransfers() ([]IncomingTransfer, error) {
	output := []IncomingTransfer{}
	rows, err := d.Db.Query("SELECT * FROM incoming_transfer WHERE size_bytes > completed_bytes;")
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var transfer IncomingTransfer

		err = rows.Scan(&transfer.Src, transfer.File, transfer.SizeBytes, transfer.CompletedBytes)

		if err != nil {
			return output, err
		}

		output = append(output, transfer)
	}

	return output, nil
}

func (d *Database) GetFilePath(dest string, name string) string {
	output := ""
	query := "SELECT file_path FROM transfer WHERE file_name = '" + name + "' AND dest LIKE '" + dest + ":%';"
	rows, err := d.Db.Query(query)
	if err != nil {
		return output
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&output)
		if err != nil {
			return output
		}
	}

	return output
}

func (d *Database) GetPendingTransfers() ([]file.File, error) {
	output := []file.File{}
	rows, err := d.Db.Query("SELECT dest, file_name, file_path, type, extension, size_bytes FROM transfer WHERE status = 'pending';")
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var f file.File

		err = rows.Scan(&f.Dest, &f.Name, &f.Path, &f.Type, &f.Extension, &f.Size)

		if err != nil {
			return output, err
		}

		output = append(output, f)
	}

	return output, nil
}

func (d *Database) FileTransferInProgress() bool {
	rows, err := d.Db.Query("SELECT ROWID FROM transfer WHERE status = 'processing';")
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (d *Database) initDB() {
	query := `CREATE TABLE device (host TEXT);
	CREATE TABLE transfer (dest TEXT, file_name TEXT, type TEXT, extension TEXT, file_path TEXT, size_bytes INTEGER, completed_bytes INTEGER, status TEXT);
	CREATE TABLE incoming_transfer (src TEXT, file_name TEXT, type TEXT, extension TEXT, size_bytes INTEGER, completed_bytes INTEGER);`
	_, err := d.Db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
