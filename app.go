package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App ...
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize : Responsável por criar a conexão com database e conectar as rotas
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

// Run : Método para inicializar a aplicação
func (a *App) Run(addr string) {}
