package main

import (
	"database/sql"
	"github.com/jbrodriguez/mlog"
	"fmt"
	"github.com/nochso/smtpd/models"
	"path"
	"strings"
)

func openDatabase() *sql.DB {
	mlog.Info("Opening SQLite database")
	dbPath := path.Join(dataDir, "mail.sqlite")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		mlog.Fatalf("Unable to open or create SQLite database file '%s': %s", dbPath, err)
	}
	models.XOLog = func(query string, data ...interface{}) {
		trimData := ""
		for _, value := range data {
			trimValue := fmt.Sprintf("%#v", value)
			if len(trimValue) > 20 {
				trimValue = trimValue[0:20] + ".."
			}
			trimData = fmt.Sprintf("%s, %s", trimData, trimValue)
		}
		mlog.Trace("SQL: %s (%s)", query, strings.TrimLeft(trimData, ", "))
	}
	return db
}

func getAddressId(address string) int {
	addr, err := models.AddressByAddress(db, address)
	if err != nil {
		addr = &models.Address{Address: address}
		addr.Insert(db)
		return addr.ID
	}
	return addr.ID
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
	_, err = db.Exec(`
CREATE TABLE "address" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "address" text NOT NULL,
  "name" text NOT NULL
);
CREATE INDEX "idx_address_address" ON "address" ("address");
CREATE UNIQUE INDEX "uidx_address_address" ON "address" ("address");
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "sender_id" integer NOT NULL,
  "content" text NOT NULL,
  "ts_received" integer NOT NULL,
  "subject" text NOT NULL,
  "is_deleted" integer NOT NULL,
  "ts_deleted" integer NOT NULL
);
CREATE INDEX "idx_mail_is_deleted_ts_deleted" ON "mail" ("is_deleted", "ts_deleted");
CREATE INDEX "idx_mail_sender_id" ON "mail" ("sender_id");
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
CREATE TABLE mail_recipient
(
    id integer PRIMARY KEY NOT NULL,
    mail_id integer NOT NULL,
    recipient_id integer NOT NULL
);
CREATE UNIQUE INDEX "uidx_mail_recipient_mail_id_recipient_id" ON "mail_recipient" ("mail_id", "recipient_id");
CREATE INDEX "idx_mail_recipient_mail_id" ON "mail_recipient" ("mail_id");
CREATE INDEX "idx_mail_recipient_recipient_id" ON "mail_recipient" ("recipient_id");
`)
	if err != nil {
		mlog.Fatalf("Unable to set up database schema: %s", err)
	}
}
