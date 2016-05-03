
CREATE TABLE "address" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "address" text NOT NULL,
  "name" text NOT NULL
);
CREATE UNIQUE INDEX "uidx_address_address" ON "address" ("address");
CREATE INDEX "idx_address_address" ON "address" ("address");
CREATE TABLE "mail" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "sender_id" integer NOT NULL,
  "recipient_id" integer NOT NULL,
  "content" text NOT NULL,
  "ts_received" integer NOT NULL,
  "subject" text NOT NULL,
  "is_deleted" integer NOT NULL,
  "ts_deleted" integer NULL,
  FOREIGN KEY ("recipient_id") REFERENCES "address" ("id"),
  FOREIGN KEY ("sender_id") REFERENCES "address" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);
CREATE INDEX "idx_mail_recipient_id" ON "mail" ("recipient_id");
CREATE INDEX "idx_mail_is_deleted_ts_deleted" ON "mail" ("is_deleted", "ts_deleted");
CREATE INDEX "idx_mail_sender_id" ON "mail" ("sender_id");
CREATE INDEX "idx_mail_ts_received" ON "mail" ("ts_received");
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
