package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Produto struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func main() {
	db, err := sql.Open("postgres", "host=localhost user=postgres password=postgres dbname=produtos_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS produtos (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(255) NOT NULL
			)`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/produtos", getProdutos(db))
	router.POST("/produtos", createProduto(db))

	log.Fatal(router.Run(":8080"))
}

func getProdutos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM produtos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var produtos []Produto
		for rows.Next() {
			var produto Produto
			err := rows.Scan(&produto.ID, &produto.Nome)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			produtos = append(produtos, produto)
		}
		c.JSON(http.StatusOK, produtos)
	}
}

func createProduto(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var produto Produto
		err := c.ShouldBindJSON(&produto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO produtos (nome) values ($1)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(produto.Nome)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "produto criado com sucesso"})
	}
}
