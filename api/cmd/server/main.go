package main

import (
    "calculator/api/internal/handlers"
    "calculator/api/internal/repository"
    "calculator/api/internal/service"
    "database/sql"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/lib/pq"
)

func main() {
    dbDSN := os.Getenv("DB_DSN")
    if dbDSN == "" {
        dbDSN = "postgres://calculator:calcpass@localhost:5432/calc_db?sslmode=disable"
    }
    var db *sql.DB
    var err error
    for i := 0; i < 10; i++ {
        db, err = sql.Open("postgres", dbDSN)
        if err == nil {
            err = db.Ping()
            if err == nil {
                break
            }
        }
        log.Printf("Waiting for database... (%d/10)", i+1)
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        log.Fatal("Could not connect to database: ", err)
    }
    defer db.Close()
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)

    createTableSQL := `
    CREATE TABLE IF NOT EXISTS calculations (
        id SERIAL PRIMARY KEY,
        operand1 DOUBLE PRECISION NOT NULL,
        operand2 DOUBLE PRECISION NOT NULL,
        operation TEXT NOT NULL,
        result DOUBLE PRECISION NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );`
    if _, err := db.Exec(createTableSQL); err != nil {
        log.Fatal("Migration failed: ", err)
    }
    log.Println("Database ready")

    repo := repository.NewPostgresCalculationRepository(db)
    calcService := service.NewCalculationService(repo)
    calcHandler := handlers.NewCalculationHandler(calcService)

    mux := http.NewServeMux()
    mux.HandleFunc("/api/calculate", calcHandler.Calculate)
    mux.HandleFunc("/api/calculations", calcHandler.GetAllCalculations)
    mux.HandleFunc("/api/calculations/", calcHandler.DeleteCalculation)

    port := ":8080"
    log.Printf("API server listening on %s", port)
    log.Fatal(http.ListenAndServe(port, mux))
}