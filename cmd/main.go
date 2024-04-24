package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Templates is a struct holding a reference to a collection of templates
type Templates struct {
	templates *template.Template
}

// Render executes a template by name and writes the output to the provided io.Writer.
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// newTemplate creates a new instance of the Templates struct
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Data struct {
	ID   int
	Name string
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Create a new instance of echo
	e := echo.New()

	// Middleware to log
	e.Use(middleware.Logger())

	// Initialize the templates
	e.Renderer = newTemplate()

	// Routes
	e.GET("/", func(c echo.Context) error {
		// Connect to the PostgreSQL database
		connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		db, err := sql.Open("postgres", connectionString)

		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Query the database
		rows, err := db.Query("SELECT id, name FROM public.companies")

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		// Create a slice of Data struct to hold the fetched data
		var data []Data

		// Iterate through the result set
		for rows.Next() {
			var d Data

			err := rows.Scan(&d.ID, &d.Name)

			if err != nil {
				log.Fatal(err)
			}
			data = append(data, d)
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		return c.Render(http.StatusOK, "index", data)
	})

	// Serve the /favicon.ico
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
