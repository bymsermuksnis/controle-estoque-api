package main

import (
	"controle-estoque/internal/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	database.InitDB()
	defer database.DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Conectado")
	})

	router.HandleFunc("/produto", createProdutoHandler).Methods("POST")

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createProdutoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
		return
	}

	var produto Produto

	err = json.Unmarshal(body, &produto)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO estoque(nome, quantidade, cor) VALUES(?, ?, ?)")
	if err != nil {
		http.Error(w, "Erro ao preparar a declaração SQL", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(produto.Nome, produto.Quantidade, produto.Cor)
	if err != nil {
		log.Printf("Erro ao inserir o produto no banco de dados: %v", err)
		http.Error(w, "Erro ao inserir o cliente no banco de dados", http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Erro ao obter o ID do novo cliente", http.StatusInternalServerError)
		return
	}
	produto.ID = int(lastID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produto)
}

type Produto struct {
	ID         int    `json:"id"`
	Nome       string `json:"nome"`
	Quantidade string `json:"quantidade"`
	Cor        string `json:"cor"`
}
