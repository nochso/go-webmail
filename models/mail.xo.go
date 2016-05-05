// Package models contains the types for schema ''.
package models

// GENERATED BY XO. DO NOT EDIT.

import "errors"

// Mail represents a row from 'mail'.
type Mail struct {
	ID          int64  // id
	SenderID    int64  // sender_id
	RecipientID int64  // recipient_id
	Content     string // content
	Subject     string // subject
	TsReceived  int64  // ts_received

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Mail exists in the database.
func (m *Mail) Exists() bool {
	return m._exists
}

// Deleted provides information if the Mail has been deleted from the database.
func (m *Mail) Deleted() bool {
	return m._deleted
}

// Insert inserts the Mail to the database.
func (m *Mail) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if m._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO mail (` +
		`sender_id, recipient_id, content, subject, ts_received` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, m.SenderID, m.RecipientID, m.Content, m.Subject, m.TsReceived)
	res, err := db.Exec(sqlstr, m.SenderID, m.RecipientID, m.Content, m.Subject, m.TsReceived)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	m.ID = int64(id)
	m._exists = true

	return nil
}

// Update updates the Mail in the database.
func (m *Mail) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !m._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if m._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE mail SET ` +
		`sender_id = ?, recipient_id = ?, content = ?, subject = ?, ts_received = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, m.SenderID, m.RecipientID, m.Content, m.Subject, m.TsReceived, m.ID)
	_, err = db.Exec(sqlstr, m.SenderID, m.RecipientID, m.Content, m.Subject, m.TsReceived, m.ID)
	return err
}

// Save saves the Mail to the database.
func (m *Mail) Save(db XODB) error {
	if m.Exists() {
		return m.Update(db)
	}

	return m.Insert(db)
}

// Delete deletes the Mail from the database.
func (m *Mail) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !m._exists {
		return nil
	}

	// if deleted, bail
	if m._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM mail WHERE id = ?`

	// run query
	XOLog(sqlstr, m.ID)
	_, err = db.Exec(sqlstr, m.ID)
	if err != nil {
		return err
	}

	// set deleted
	m._deleted = true

	return nil
}

// AddressByRecipientID returns the Address associated with the Mail's RecipientID (recipient_id).
//
// Generated from foreign key 'mail_recipient_id_fkey'.
func (m *Mail) AddressByRecipientID(db XODB) (*Address, error) {
	return AddressByID(db, m.RecipientID)
}

// AddressBySenderID returns the Address associated with the Mail's SenderID (sender_id).
//
// Generated from foreign key 'mail_sender_id_fkey'.
func (m *Mail) AddressBySenderID(db XODB) (*Address, error) {
	return AddressByID(db, m.SenderID)
}

// MailsByRecipientID retrieves a row from 'mail' as a Mail.
//
// Generated from index 'idx_mail_recipient_id'.
func MailsByRecipientID(db XODB, recipientID int64) ([]*Mail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, sender_id, recipient_id, content, subject, ts_received ` +
		`FROM mail ` +
		`WHERE recipient_id = ?`

	// run query
	XOLog(sqlstr, recipientID)
	q, err := db.Query(sqlstr, recipientID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Mail{}
	for q.Next() {
		m := Mail{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.Subject, &m.TsReceived)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

// MailsBySenderID retrieves a row from 'mail' as a Mail.
//
// Generated from index 'idx_mail_sender_id'.
func MailsBySenderID(db XODB, senderID int64) ([]*Mail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, sender_id, recipient_id, content, subject, ts_received ` +
		`FROM mail ` +
		`WHERE sender_id = ?`

	// run query
	XOLog(sqlstr, senderID)
	q, err := db.Query(sqlstr, senderID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Mail{}
	for q.Next() {
		m := Mail{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.Subject, &m.TsReceived)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

// MailsByTsReceived retrieves a row from 'mail' as a Mail.
//
// Generated from index 'idx_mail_ts_received'.
func MailsByTsReceived(db XODB, tsReceived int64) ([]*Mail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, sender_id, recipient_id, content, subject, ts_received ` +
		`FROM mail ` +
		`WHERE ts_received = ?`

	// run query
	XOLog(sqlstr, tsReceived)
	q, err := db.Query(sqlstr, tsReceived)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Mail{}
	for q.Next() {
		m := Mail{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.Subject, &m.TsReceived)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

// MailByID retrieves a row from 'mail' as a Mail.
//
// Generated from index 'mail_id_pkey'.
func MailByID(db XODB, id int64) (*Mail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, sender_id, recipient_id, content, subject, ts_received ` +
		`FROM mail ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	m := Mail{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.Subject, &m.TsReceived)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
