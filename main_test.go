package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

// <<FUNÇÕES TEST>>

// TestMain : Função que será executada antes de todos os testes para executar alguns comandos necessários
func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "root", "rest_api_example")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

// TestEmptyTable : Este teste irá deletar os dados de user no banco e mandar um 'GET request'
// para o endpoint '/users', deve-se retornar um array vazio
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Array vazio esperado. Obteve: %s", body)
	}
}

// TestGetNonExistentUser : Este teste verifica se o status code 404 é recebido através de uma requisição
// de um usuário não existente
func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/user/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Chave 'error' da resposta, esperado para ser 'User not found'. Obteve: '%s'", m["error"])
	}
}

// TestCreateUser : Neste teste é adicionado um novo usuário no banco, e depois testamos se a resposta é
// correspondente ao que foi adicionado
func TestCreateUser(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test user","age":30}`)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test user" {
		t.Errorf("Esperava-se que nome de usuário fosse igual 'test user'. Obteve: '%v'", m["name"])
	}

	if m["age"] != 30.0 {
		t.Errorf("Esperava-se que idade do usuário fosse '30'. Obteve: '%v'", m["age"])
	}

	// O id é comparado com 1.0 porque o JSON unmarshaling converte numbers para floats,
	// quando o alvo é um map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Esperava-se que o ID de usuário fosse igual '1'. Obteve: '%v'", m["id"])
	}
}

// TestGetUser : Este teste irá adicionar um novo usuário no banco e irá checar se o endpoint resulta em uma
// resposta HTTP com status code 200
func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// TestUpdateUser : Este teste adiciona um novo usuário no banco e irá usar o endpoint de 'PUT' para atualizar
// o usuário do banco
func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	payload := []byte(`{"name":"test user - updated name","age":21}`)
	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Esperava-se que o ID se mantivesse o mesmo (%v). Obteve: %v", originalUser["id"], m["id"])
	}
	if m["name"] == originalUser["name"] {
		t.Errorf("Esperava-se que o nome mudasse de '%v' para '%v'. Obteve: '%v'", originalUser["name"], m["name"], m["name"])
	}
	if m["age"] == originalUser["age"] {
		t.Errorf("Esperava-se que a idade mudasse de '%v' para '%v'. Obteve: '%v'", originalUser["age"], m["age"], m["age"])
	}
}

// TestDeleteUser : Este teste irá criar um novo usuário, irá testar se ele foi adicionado usando 'GET', logo
// em seguida irá excluir o usuário criado anteriormente, e irá testar novamente se ele foi excluído, e para
// finalizar, ele irá fazer uma requisição 'GET' para verificar se existe um usuário no banco, e por final
// deve se retornar 'http.StatusNotFound'
func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)
	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// <<FUNÇÕES>>

// addUsers : Função para adicionar um user no banco de dados
func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)", ("User " + strconv.Itoa(i+1)), ((i + 1) * 10))
		a.DB.Exec(statement)
	}
}

// ensureTableExists : Verifica se a table que precisamos usar existe
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// clearTable :  Função para limpar uma table
func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

// executeRequest :  Função para executar uma requisição
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode : Função para testar se a 'HTTP response code' é o esperado
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Código resposta esperado: %d. Código resposta recebido: %d\n", expected, actual)
	}
}

// tableCreationQuery : Em caso a tabela não exista, cria uma
const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
	id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
)`
