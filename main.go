package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type DataItem struct {
    ID       string `json:"id"`
    Metadata struct {
        Name string `json:"name"`
    } `json:"metadata"`
}


func main() {

    // Database URL format: postgres://user:password@host:port/database_name?sslmode=disable
    db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5439/admin-management?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to the PostgreSQL database!")
    
    r := gin.Default()

    // Define a simple route
    r.GET("/acknowledge", func(c *gin.Context) {
        rows, err := db.Query("SELECT id, metadata FROM acknowledge")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
            return
        }
        defer rows.Close()

        var data []DataItem

        for rows.Next() {
            var id string
            var metadataJSON string
            if err := rows.Scan(&id, &metadataJSON); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error"})
                return
            }

            var metadata struct {
                Name string `json:"name"`
            }

            if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error"})
                return
            }

            item := DataItem{
                ID:       id,
                Metadata: metadata,
            }
            data = append(data, item)
        }

        // Return the data array as JSON
        c.JSON(http.StatusOK, data)
    })

    // Start the server on port 8080
    r.Run(":8080")
}

func loadEnv() {
	panic("unimplemented")
}
