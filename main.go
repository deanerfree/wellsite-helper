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

	// Serve static files
	e.Static("/static", "static")
	// Set up the template renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFiles("view/layout/index.html", "view/wellpathTemplate.html", "view/uploadTemplate.html")),
	}
	e.Renderer = renderer

	// Define the handler
	e.GET("/", func(c echo.Context) error {
		// Pass data if needed (nil here for simplicity)
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/wellpath", func(c echo.Context) error {
		lasFile := "mockData/test3.las"
		// surveyFile := "mockData/survey3.txt"

		data, err := ScanLasFile(lasFile)
		if err != nil {
			fmt.Println(err)
		}
		// Extract the Well Name
		wellData := data.WellInformation.Fields
		WellName := wellData["WELL"]
		// if the Fields contains an key with the name "gas" then we can assume that the data is a gas data and change the object key to "total_gas"

		DepthData := data.WellData.Fields

		// fmt.Printf("Data: %s\n", data.WellInformation.Fields)

		serializedDepthData, err := json.Marshal(DepthData)
		if err != nil {
			log.Printf("Error serializing DepthData: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to process DepthData")
		}

		// create an struct to hold the data of the first row of the data
		firstRowData := DepthData[0]
		// create an array to hold the keys of the first row data
		var keys []string
		// loop through the first row data and append the keys to the keys array

		for key := range firstRowData {
			keys = append(keys, key)
		}

		if keys[0][0] != 'D' && keys[0][0] != 'M' { // Check if the first key doesn't start with 'D' or 'M'
			for i, key := range keys {
				if key[0] == 'D' || key[0] == 'M' { // Find a key starting with 'D' or 'M'
					// Remove the key from the array
					keys = append(keys[:i], keys[i+1:]...)
					// Insert the key at the front of the array
					keys = append([]string{key}, keys...)
					break // Exit the loop after moving the key
				}
			}
		}

		// Move specific items, like "rop", to the end
		for i, key := range keys {
			if strings.ToLower(key) == "rop" {
				keys = append(keys[:i], keys[i+1:]...) // Remove "rop"
				keys = append(keys, "rop")             // Add "rop" to the end
				break
			}
		}

		// if err != nil {
		// 	log.Printf("Error serializing Keys: %v", err)
		// 	return c.String(http.StatusInternalServerError, "Failed to process Keys")
		// }

		return c.Render(http.StatusOK, "wellpathTemplate.html", map[string]interface{}{
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
