PRAGMA foreign_keys = ON;
CREATE TABLE "address" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "address" text NOT NULL,
  "name" text NOT NULL
);
CREATE UNIQUE INDEX "uidx_address_address_name" ON "address" ("address", "name");
CREATE TABLE "address_cc" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "mail_id" integer NOT NULL,
  "address_id" integer NOT NULL,
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id"),
  FOREIGN KEY ("address_id") REFERENCES "address" ("id")
);
CREATE INDEX "idx_address_cc_address_id" ON "address_cc" ("address_id");
CREATE INDEX "idx_address_cc_mail_id" ON "address_cc" ("mail_id");
CREATE TABLE "address_from" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "mail_id" integer NOT NULL,
  "address_id" integer NOT NULL,
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id"),
  FOREIGN KEY ("address_id") REFERENCES "address" ("id")
);
CREATE INDEX "idx_address_from_address_id" ON "address_from" ("address_id");
CREATE INDEX "idx_address_from_mail_id" ON "address_from" ("mail_id");
CREATE TABLE "address_replyto" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "mail_id" integer NOT NULL,
  "address_id" integer NOT NULL,
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id"),
  FOREIGN KEY ("address_id") REFERENCES "address" ("id")
);
CREATE INDEX "idx_address_replyto_address_id" ON "address_replyto" ("address_id");
CREATE INDEX "idx_address_replyto_mail_id" ON "address_replyto" ("mail_id");
CREATE TABLE "address_to" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "mail_id" integer NOT NULL,
  "address_id" integer NOT NULL,
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id"),
  FOREIGN KEY ("address_id") REFERENCES "address" ("id")
);
CREATE INDEX "idx_address_to_mail_id" ON "address_to" ("mail_id");
CREATE INDEX "idx_address_to_address_id" ON "address_to" ("address_id");
CREATE TABLE "flag" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text NOT NULL
);
CREATE UNIQUE INDEX "uidx_flag_name" ON "flag" ("name");
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
CREATE INDEX "idx_flag_mail_mail_id" ON "flag_mail" ("mail_id");
CREATE INDEX "idx_flag_mail_flag_id_mail_id" ON "flag_mail" ("flag_id", "mail_id");
CREATE TABLE "header" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text NOT NULL
);
CREATE UNIQUE INDEX "uidx_header_name" ON "header" ("name");
CREATE TABLE "header_mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "header_id" integer NOT NULL,
  "mail_id" integer NOT NULL,
  "header_value" text NOT NULL,
  FOREIGN KEY ("header_id") REFERENCES "header" ("id"),
  FOREIGN KEY ("mail_id") REFERENCES "mail" ("id")
);
CREATE INDEX "idx_header_mail_header_id_header_value" ON "header_mail" ("header_id", "header_value");
CREATE INDEX "idx_header_mail_mail_id" ON "header_mail" ("mail_id");
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "content" text NOT NULL,
  "ts_received" integer NOT NULL
);
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
CREATE TABLE "v_mail" ("id" integer, "content" text, "ts_received" integer, "id:1" integer, "address" text, "name" text, "id:2" integer, "address:1" text, "name:1" text, "flags" );
DROP TABLE IF EXISTS "v_mail";
CREATE VIEW "v_mail" AS
SELECT mail.*, addfrom.*, addto.*, GROUP_CONCAT(f.name) flags
FROM mail
LEFT JOIN address_from adfr ON adfr.mail_id = mail.id
LEFT JOIN address addfrom ON addfrom.id = adfr.address_id
LEFT JOIN address_to adto ON adto.mail_id = mail.id
LEFT JOIN address addto ON addto.id = adto.address_id
LEFT JOIN flag_mail fm ON fm.mail_id = mail.id
LEFT JOIN flag f ON f.id = fm.flag_id
GROUP BY mail.id;
