DROP TABLE IF EXISTS "address";
CREATE TABLE "address" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "address" text NOT NULL,
  "name" text NOT NULL
);
CREATE INDEX "idx_address_address" ON "address" ("address");
CREATE UNIQUE INDEX "uidx_address_address" ON "address" ("address");
DROP TABLE IF EXISTS "mail";
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "sender_id" integer NOT NULL,
  "recipient_id" integer NOT NULL,
  "content" text NOT NULL,
  "ts_received" integer NOT NULL,
  "subject" text NOT NULL,
  "is_deleted" integer NOT NULL,
  "ts_deleted" integer
);
CREATE INDEX "idx_mail_is_deleted_ts_deleted" ON "mail" ("is_deleted", "ts_deleted");
CREATE INDEX "idx_mail_sender_id" ON "mail" ("sender_id");
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
DROP VIEW IF EXISTS "v_mail";
CREATE TABLE "v_mail" ("id" integer, "subject" text, "content" text, "received" , "deleted" , "sender" text, "recipient" text, "is_deleted" integer);
DROP TABLE IF EXISTS "v_mail";
CREATE VIEW "v_mail" AS
SELECT
  mail.id,
  mail.subject,
  mail.content,
  datetime(mail.ts_received, 'unixepoch', 'localtime')         AS received,
  CASE WHEN mail.ts_deleted IS NULL
    THEN ''
  ELSE datetime(mail.ts_deleted, 'unixepoch', 'localtime') END AS deleted,
  s_addr.address                                               AS sender,
  r_addr.address                                               AS recipient,
  mail.is_deleted
FROM mail
  INNER JOIN address s_addr ON mail.sender_id = s_addr.id
  INNER JOIN address r_addr ON mail.recipient_id = r_addr.id;
