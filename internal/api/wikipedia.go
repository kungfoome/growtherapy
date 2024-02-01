package api

import (
    "encoding/json"
    "fmt"
    "context"
	"io"
    "net/http"
    "time"
    "github.com/kungfoome/growtherapy/internal/httpclient"

	"golang.org/x/time/rate"
)

type WikipediaAPI struct {
    Client    httpclient.Client
    RateLimiter *rate.Limiter
}

type wikiAPIResponse struct {
    Items []struct {
        Views int `json:"views"`
    } `json:"items"`
}

func NewWikipediaAPI(client httpclient.Client) *WikipediaAPI {
	// Wikipedia has a rate limit of less than 1 request/second when not authenticated
	limiter := rate.NewLimiter(1, 1)
    return &WikipediaAPI{
        Client:    client,
        RateLimiter: limiter,
    }
}

func (api *WikipediaAPI) GetArticleViews(ctx context.Context, article string, year int, month int) (int, error) {
    // Use the provided context with the rate limiter
    err := api.RateLimiter.Wait(ctx)
    if err != nil {
        return 0, err
    }

    var apiResp wikiAPIResponse
    retries := 3
    for i := 0; i < retries; i++ {
        if i > 0 {
            // Exponential backoff
            time.Sleep(time.Duration(100 * (1 << i)) * time.Millisecond)
        }

        startDate, endDate := getStartAndEndDate(year, month)
        url := fmt.Sprintf("https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia/all-access/all-agents/%s/monthly/%s/%s", article, startDate, endDate)
        
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            continue // Retry on creating request error
        }
        
        resp, err := api.Client.Do(req)
        if err != nil {
            continue // Retry on request error
        }
        defer resp.Body.Close()

        // Use io.ReadAll instead of ioutil.ReadAll
        if resp.StatusCode == http.StatusOK {
            if err := json.NewDecoder(resp.Body).Decode(&apiResp); err == nil {
                return apiResp.Items[0].Views, nil // Success
            }
        } else {
            body, _ := io.ReadAll(resp.Body)
            fmt.Printf("Wikipedia API error (%d): %s\n", resp.StatusCode, string(body))
        }
    }

    return 0, fmt.Errorf("failed to get article views after %d attempts", retries)
}

// Helper function to calculate start and end dates
func getStartAndEndDate(year int, month int) (string, string) {
    location, _ := time.LoadLocation("UTC")
    firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location)
    lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
    return firstOfMonth.Format("20060102"), lastOfMonth.Format("20060102")
}