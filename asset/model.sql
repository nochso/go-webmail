PRAGMA foreign_keys = ON;
CREATE TABLE "address" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "address" text NOT NULL,
  "name" text NOT NULL
);
CREATE INDEX "idx_address_address" ON "address" ("address");
CREATE UNIQUE INDEX "uidx_address_address" ON "address" ("address");
CREATE TABLE "flag" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text NOT NULL
);
INSERT INTO "flag" ("id", "name") VALUES (1,	'\Seen');
INSERT INTO "flag" ("id", "name") VALUES (2,	'\Answered');
INSERT INTO "flag" ("id", "name") VALUES (3,	'\Flagged');
INSERT INTO "flag" ("id", "name") VALUES (4,	'\Deleted');
INSERT INTO "flag" ("id", "name") VALUES (5,	'\Draft');
INSERT INTO "flag" ("id", "name") VALUES (6,	'\Recent');
CREATE TABLE "flag_mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "flag_id" integer NOT NULL,
  "mail_id" integer NOT NULL,
  FOREIGN KEY ("flag_id") REFERENCES "flag" ("id"),
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id")
);
CREATE INDEX "idx_flag_mail_flag_id_mail_id" ON "flag_mail" ("flag_id", "mail_id");
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "sender_id" integer NOT NULL,
  "recipient_id" integer NOT NULL,
  "content" text NOT NULL,
  "subject" text NOT NULL,
  "ts_received" integer NOT NULL,
  FOREIGN KEY ("recipient_id") REFERENCES "address" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  FOREIGN KEY ("sender_id") REFERENCES "address" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
CREATE INDEX "idx_mail_sender_id" ON "mail" ("sender_id");
CREATE INDEX "idx_mail_recipient_id" ON "mail" ("recipient_id");