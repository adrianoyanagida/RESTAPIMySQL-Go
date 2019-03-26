package main

import (
	"database/sql"
	"errors"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *user) getUser(db *sql.DB) error {
	return errors.New("Não implementado")
}

func (u *user) updateUser(db *sql.DB) error {
	return errors.New("Não implementado")
}

func (u *user) deleteUser(db *sql.DB) error {
	return errors.New("Não implementado")
}

func (u *user) createUser(db *sql.DB) error {
	return errors.New("Não implementado")
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	return nil, errors.New("Não implementado")
}
