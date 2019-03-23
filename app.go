package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App ...
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize : Responsável por criar a conexão com database e conectar as rotas
func (a *App) Initialize(user, password, dbname string) {}

// Run : Método para inicializar a aplicação
func (a *App) Run(addr string) {}
