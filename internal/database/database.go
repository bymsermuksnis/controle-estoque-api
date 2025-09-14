package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	dsn := "root:@tcp(127.0.0.1:3306)/controle-estoque"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão com o banco de dados: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao pingar o banco de dados: %v", err)
	}
	fmt.Println("Conexão com o banco de dados MySQL estabelecida!")

	if err != nil {
		log.Fatalf("Erro ao criar a tabela clientes: %v", err)
	}
	fmt.Println("Tabela 'clientes' verificada/criada com sucesso.")
}
