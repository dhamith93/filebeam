package database

import (
	"database/sql"
	"log"
	"os"

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

func (d *Database) AddTransfer(dest string, path string, size int64) error {
	return d.execute(
		"INSERT INTO transfer (dest, file_path, size_bytes) VALUES (?, ?, ?);",
		dest, path, size, 0,
	)
}

func (d *Database) AddIncomingTransfer(src string, file string, size int64) error {
	return d.execute(
		"INSERT INTO incoming_transfer (src, file, size_bytes, completed_bytes) VALUES (?, ?, ?, ?);",
		src, file, size, 0,
	)
}

func (d *Database) UpdateIncomingTransfer(src string, file string, completed int64) error {
	return d.execute(
		"UPDATE incoming_transfer SET completed_bytes = ? WHERE file = ? AND src = ?",
		completed, file, src,
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

func (d *Database) initDB() {
	query := `CREATE TABLE device (host TEXT);
	CREATE TABLE transfer (dest TEXT, file_path TEXT, size_bytes INTEGER);
	CREATE TABLE incoming_transfer (src TEXT, file TEXT, size_bytes INTEGER, completed_bytes INTEGER);`
	_, err := d.Db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
