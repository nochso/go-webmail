// Package model contains the types for schema ''.
package model

// GENERATED BY XO. DO NOT EDIT.

import "errors"

// FlagMail represents a row from 'flag_mail'.
type FlagMail struct {
	ID     int64 // id
	FlagID int64 // flag_id
	MailID int64 // mail_id

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the FlagMail exists in the database.
func (fm *FlagMail) Exists() bool {
	return fm._exists
}

// Deleted provides information if the FlagMail has been deleted from the database.
func (fm *FlagMail) Deleted() bool {
	return fm._deleted
}

// Insert inserts the FlagMail to the database.
func (fm *FlagMail) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if fm._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO flag_mail (` +
		`flag_id, mail_id` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, fm.FlagID, fm.MailID)
	res, err := db.Exec(sqlstr, fm.FlagID, fm.MailID)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	fm.ID = int64(id)
	fm._exists = true

	return nil
}

// Update updates the FlagMail in the database.
func (fm *FlagMail) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fm._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if fm._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE flag_mail SET ` +
		`flag_id = ?, mail_id = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, fm.FlagID, fm.MailID, fm.ID)
	_, err = db.Exec(sqlstr, fm.FlagID, fm.MailID, fm.ID)
	return err
}

// Save saves the FlagMail to the database.
func (fm *FlagMail) Save(db XODB) error {
	if fm.Exists() {
		return fm.Update(db)
	}

	return fm.Insert(db)
}

// Delete deletes the FlagMail from the database.
func (fm *FlagMail) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fm._exists {
		return nil
	}

	// if deleted, bail
	if fm._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM flag_mail WHERE id = ?`

	// run query
	XOLog(sqlstr, fm.ID)
	_, err = db.Exec(sqlstr, fm.ID)
	if err != nil {
		return err
	}

	// set deleted
	fm._deleted = true

	return nil
}

// Flag returns the Flag associated with the FlagMail's FlagID (flag_id).
//
// Generated from foreign key 'flag_mail_flag_id_fkey'.
func (fm *FlagMail) Flag(db XODB) (*Flag, error) {
	return FlagByID(db, fm.FlagID)
}

// MailByMailID returns the Mail associated with the FlagMail's MailID (mail_id).
//
// Generated from foreign key 'flag_mail_mail_id_fkey'.
func (fm *FlagMail) MailByMailID(db XODB) (*Mail, error) {
	return MailByID(db, fm.MailID)
}

// FlagMailByID retrieves a row from 'flag_mail' as a FlagMail.
//
// Generated from index 'flag_mail_id_pkey'.
func FlagMailByID(db XODB, id int64) (*FlagMail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, flag_id, mail_id ` +
		`FROM flag_mail ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	fm := FlagMail{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&fm.ID, &fm.FlagID, &fm.MailID)
	if err != nil {
		return nil, err
	}

	return &fm, nil
}

// FlagMailsByFlagIDMailID retrieves a row from 'flag_mail' as a FlagMail.
//
// Generated from index 'idx_flag_mail_flag_id_mail_id'.
func FlagMailsByFlagIDMailID(db XODB, flagID int64, mailID int64) ([]*FlagMail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, flag_id, mail_id ` +
		`FROM flag_mail ` +
		`WHERE flag_id = ? AND mail_id = ?`

	// run query
	XOLog(sqlstr, flagID, mailID)
	q, err := db.Query(sqlstr, flagID, mailID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*FlagMail{}
	for q.Next() {
		fm := FlagMail{
			_exists: true,
		}

		// scan
		err = q.Scan(&fm.ID, &fm.FlagID, &fm.MailID)
		if err != nil {
			return nil, err
		}

		res = append(res, &fm)
	}

	return res, nil
}

// FlagMailsByMailID retrieves a row from 'flag_mail' as a FlagMail.
//
// Generated from index 'idx_flag_mail_mail_id'.
func FlagMailsByMailID(db XODB, mailID int64) ([]*FlagMail, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, flag_id, mail_id ` +
		`FROM flag_mail ` +
		`WHERE mail_id = ?`

	// run query
	XOLog(sqlstr, mailID)
	q, err := db.Query(sqlstr, mailID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*FlagMail{}
	for q.Next() {
		fm := FlagMail{
			_exists: true,
		}

		// scan
		err = q.Scan(&fm.ID, &fm.FlagID, &fm.MailID)
		if err != nil {
			return nil, err
		}

		res = append(res, &fm)
	}

	return res, nil
}
