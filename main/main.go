package main

import (
  "context"
  "log"
  "net/http"
  "time"
  "os"
  "github.com/gin-gonic/gin"
  "github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func main() {
  gin.SetMode(gin.ReleaseMode)
  databaseUrl := "postgres://yrysmat:123456789@localhost:5432/yrysanime"
  var err error

  start := time.Now() // Временная метка для диагностики
  conn, err = pgx.Connect(context.Background(), databaseUrl)
  if err != nil {
    log.Fatalf("Unable to connect to database: %v\n", err)
    os.Exit(1)
  }
  defer conn.Close(context.Background())
  log.Printf("Connected to PostgreSQL database in %s", time.Since(start))

  r := gin.Default()
  r.SetTrustedProxies(nil) // Установите доверенные прокси, если это необходимо

  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "Hello, World!",
    })
  })

  r.GET("/items", getItems)

  r.Run(":5000") // запускаем сервер на порту 5000
}

func getItems(c *gin.Context) {
  rows, err := conn.Query(context.Background(), "SELECT id, name FROM items")
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  defer rows.Close()

  var items []map[string]interface{}
  for rows.Next() {
    var id int
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    items = append(items, map[string]interface{}{
      "id":   id,
      "name": name,
    })
  }

  if err = rows.Err(); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, items)
}