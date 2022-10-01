package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func formatDSN() string {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env.development file")
	}
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := "localhost"
	dbPort := "3306"
	dbName := os.Getenv("MYSQL_DATABASE")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
}

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", formatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type albumDB struct {
	db *sql.DB
}

func (albumdb albumDB) getAlbums(c *gin.Context) {
	var albums []album
	rows, err := albumdb.db.Query("SELECT * FROM albums")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		newAlbum := album{}
		rows.Scan(&newAlbum.ID, &newAlbum.Title, &newAlbum.Artist, &newAlbum.Price)
		albums = append(albums, newAlbum)
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func (albumdb albumDB) getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	row := albumdb.db.QueryRow("SELECT * FROM albums WHERE id = ?", id)
	if row == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	newAlbum := album{}
	row.Scan(&newAlbum.ID, &newAlbum.Title, &newAlbum.Artist, &newAlbum.Price)
	c.IndentedJSON(http.StatusOK, newAlbum)
}

func (albumdb albumDB) postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	row := albumdb.db.QueryRow("INSERT INTO albums (id, title, artist, price) VALUES (?, ?, ?, ?)", newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if row == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	row.Scan(&newAlbum.ID, &newAlbum.Title, &newAlbum.Artist, &newAlbum.Price)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// albums slice to seed record album data.
func main() {
	db := NewDB()
	defer db.Close()
	if err := db.PingContext(context.Background()); err != nil {
		log.Printf("failed to ping err = %s", err.Error())
		return
	}

	albumdb := albumDB{db: db}
	router := gin.Default()
	router.GET("/albums", albumdb.getAlbums)
	router.GET("/albums/:id", albumdb.getAlbumByID)
	router.POST("/albums", albumdb.postAlbums)
	router.Run("localhost:8080")
}
