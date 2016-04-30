package main

import (
	"log"
	"database/sql"
)

func prepareDatabase(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'mail'").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 1 {
		return
	}
	log.Println("Setting up database schema")
	_, err = db.Exec(`
CREATE TABLE address
(
    id integer PRIMARY KEY NOT NULL,
    address TEXT NOT NULL
);
CREATE UNIQUE INDEX "uidx_address_address" ON "address" ("address");
CREATE INDEX "idx_address_address" ON "address" ("address");
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "sender_id" integer NOT NULL,
  "content" text NOT NULL,
  "ts_received" integer NOT NULL,
  "subject" text NOT NULL,
  "is_deleted" integer NOT NULL,
  "ts_deleted" integer NOT NULL
);
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
CREATE INDEX "idx_mail_sender_id" ON "mail" ("sender_id");
CREATE INDEX "idx_mail_is_deleted_ts_deleted" ON "mail" ("is_deleted", "ts_deleted");
CREATE TABLE mail_recipient
(
    id integer PRIMARY KEY NOT NULL,
    mail_id integer NOT NULL,
    recipient_id integer NOT NULL
);
CREATE INDEX "idx_mail_recipient_recipient_id" ON "mail_recipient" ("recipient_id");
CREATE INDEX "idx_mail_recipient_mail_id" ON "mail_recipient" ("mail_id");
CREATE UNIQUE INDEX "uidx_mail_recipient_mail_id_recipient_id" ON "mail_recipient" ("mail_id", "recipient_id");
`)
	if err != nil {
		log.Fatalf("Unable to set up database schema: %s", err)
	}
}
