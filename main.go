package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	// "layout"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "strings"
	// "text/scanner"
)

// TemplateRenderer integrates html/template with Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render method implementation for Echo's renderer interface
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// setup server
	godotenv.Load()

	e := echo.New()
	// Little bit of middlewares for housekeeping
	// e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.Recover())
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// This will initiate our template renderer

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Set up the template renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFiles("view/layout/index.html", "view/lasTemplate.html")),
	}
	e.Renderer = renderer

	// Define the handler
	e.GET("/", func(c echo.Context) error {
		// Pass data if needed (nil here for simplicity)
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/las", func(c echo.Context) error {
		lasFile := "mockData/test2.las"
		data, err := scanLasFile(lasFile)
		if err != nil {
			fmt.Println(err)
		}
		// Extract the Well Name
		wellData := data.WellInformation.Fields
		WellName := wellData["WELL"]
		DepthData := data.WellData.Fields

		// fmt.Printf("Data: %s\n", data.WellInformation.Fields)

		serializedDepthData, err := json.Marshal(DepthData)
		if err != nil {
			log.Printf("Error serializing DepthData: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to process DepthData")
		}

		return c.Render(http.StatusOK, "lasTemplate.html", map[string]interface{}{
			"Wellname":  WellName.Value,
			"DepthData": template.JS(serializedDepthData), // Safely inject JSON
		})
	})

	// http.HandleFunc("/las", func(w http.ResponseWriter, r *http.Request) {
	// 	lasFile := "mockData/test2.las"
	// 	data := scanLasFile(lasFile)
	// 	RespondWithJson(w, data, r)
	// })

	// http.ListenAndServe(":"+port, nil)
	e.Logger.Fatal(e.Start(":" + port))

	// scanLasFile("mockData/test2.las")
	// fmt.Printf("LasData: %s\n", LasData)
}
