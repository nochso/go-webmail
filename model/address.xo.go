// Package model contains the types for schema ''.
package model

// GENERATED BY XO. DO NOT EDIT.

import "errors"

// Address represents a row from 'address'.
type Address struct {
	ID      int64  // id
	Address string // address
	Name    string // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Address exists in the database.
func (a *Address) Exists() bool {
	return a._exists
}

// Deleted provides information if the Address has been deleted from the database.
func (a *Address) Deleted() bool {
	return a._deleted
}

// Insert inserts the Address to the database.
func (a *Address) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO address (` +
		`address, name` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, a.Address, a.Name)
	res, err := db.Exec(sqlstr, a.Address, a.Name)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	a.ID = int64(id)
	a._exists = true

	return nil
}

// Update updates the Address in the database.
func (a *Address) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE address SET ` +
		`address = ?, name = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, a.Address, a.Name, a.ID)
	_, err = db.Exec(sqlstr, a.Address, a.Name, a.ID)
	return err
}

// Save saves the Address to the database.
func (a *Address) Save(db XODB) error {
	if a.Exists() {
		return a.Update(db)
	}

	return a.Insert(db)
}

// Delete deletes the Address from the database.
func (a *Address) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM address WHERE id = ?`

	// run query
	XOLog(sqlstr, a.ID)
	_, err = db.Exec(sqlstr, a.ID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// AddressByID retrieves a row from 'address' as a Address.
//
// Generated from index 'address_id_pkey'.
func AddressByID(db XODB, id int64) (*Address, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, address, name ` +
		`FROM address ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	a := Address{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&a.ID, &a.Address, &a.Name)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AddressByAddressName retrieves a row from 'address' as a Address.
//
// Generated from index 'uidx_address_address_name'.
func AddressByAddressName(db XODB, address string, name string) (*Address, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, address, name ` +
		`FROM address ` +
		`WHERE address = ? AND name = ?`

	// run query
	XOLog(sqlstr, address, name)
	a := Address{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, address, name).Scan(&a.ID, &a.Address, &a.Name)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
