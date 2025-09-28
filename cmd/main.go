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

	router.HandleFunc("/produtos", getProdutosHandler).Methods("GET")

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getProdutosHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT idestoque, nome, quantidade, cor FROM estoque")
	if err != nil {
		log.Printf("Erro ao buscar produtos no DB: %v", err)
		http.Error(w, "Erro interno ao buscar produtos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	produtos := []Produto{}

	for rows.Next() {
		var p Produto

		if err := rows.Scan(&p.ID, &p.Nome, &p.Quantidade, &p.Cor); err != nil {
			log.Printf("Erro ao mapear linha do DB: %v", err)
			continue
		}
		produtos = append(produtos, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erro durante a iteração das linhas: %v", err)
		http.Error(w, "Erro interno ao processar dados", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produtos)
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
