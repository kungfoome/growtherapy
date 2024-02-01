package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/kungfoome/growtherapy/internal/api"
    "github.com/kungfoome/growtherapy/internal/httpclient"
)

// getViewsHandler processes requests to /getViews, fetching view counts for a given article and month/year.
func getViewsHandler(w http.ResponseWriter, r *http.Request) {
    // Extract query parameters
    article := r.URL.Query().Get("article")
    yearStr := r.URL.Query().Get("year")
    monthStr := r.URL.Query().Get("month")

    // Convert year and month to integers
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        http.Error(w, "Invalid year parameter", http.StatusBadRequest)
        return
    }

    month, err := strconv.Atoi(monthStr)
    if err != nil {
        http.Error(w, "Invalid month parameter", http.StatusBadRequest)
        return
    }

    // Ensure wikipediaAPI is initialized (see below)
    views, err := wikipediaAPI.GetArticleViews(context.Background(), article, year, month)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching article views: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with the view count in JSON format
    response := struct {
        Views int `json:"views"`
    }{
        Views: views,
    }
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// wikipediaAPI is the global instance of the WikipediaAPI client.
var wikipediaAPI *api.WikipediaAPI

func main() {
    // Initialize the WikipediaAPI client with a real HTTP client from the httpclient package.
    // This is a simplistic way to set it up; consider using dependency injection for more complex applications.
    wikipediaAPI = api.NewWikipediaAPI(httpclient.NewClient())

    // Register the handler function with the HTTP server
    http.HandleFunc("/getViews", getViewsHandler)

    // Start the HTTP server
    fmt.Println("Server is running on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
