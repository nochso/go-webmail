package main

import (
	"database/sql"
	"fmt"
	"net/mail"
	"path"
	"strings"

	"github.com/nochso/go-webmail/model"
	"github.com/nochso/mlog"
)

func openDatabase() *sql.DB {
	mlog.Trace("Opening SQLite database")
	dbPath := path.Join(dataDir, "mail.sqlite")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		mlog.Fatalf("Unable to open or create SQLite database file '%s': %s", dbPath, err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	var fk int
	db.QueryRow("PRAGMA foreign_keys;").Scan(&fk)
	if err != nil || fk != 1 {
		mlog.Fatalf("Unable to enforce foreign key constraints: %s", err)
	}
	mlog.Trace("Enforcing SQLite foreign key constraints")
	model.XOLog = func(query string, data ...interface{}) {
		for _, value := range data {
			trimValue := fmt.Sprintf("%#v", value)
			if len(trimValue) > 40 {
				trimValue = trimValue[0:40] + ".."
				if trimValue[0] == '"' {
					trimValue += "\""
				}
			}
			query = strings.Replace(query, "?", trimValue, 1)
		}
		mlog.Trace("xo SQL: %s", query)
	}
	return db
}

func getAddressId(address *mail.Address) int64 {
	addrRow, err := model.AddressByAddress(db, address.Address)
	if err != nil {
		addrRow = &model.Address{
			Address: address.Address,
			Name:    address.Name,
		}
		addrRow.Insert(db)
		return addrRow.ID
	}
	return addrRow.ID
}

func prepareDatabase(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'mail'").Scan(&count)
	if err != nil {
		mlog.Fatal(err)
	}
	if count == 1 {
		return
	}
	mlog.Info("Setting up database schema")
	sql, err := Asset("model.sql")
	if err != nil {
		mlog.Fatalf("Unable to read embedded model.sql file: %s", err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		mlog.Fatalf("Unable to set up database schema: %s", err)
	}
}
