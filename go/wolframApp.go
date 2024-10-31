package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"io"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

// Struct to parse Wolfram Alpha API response
type WolframResponse struct {
    XMLName xml.Name `xml:"queryresult"`
    Success bool     `xml:"success,attr"`
    Pods    []Pod    `xml:"pod"`
}

type Pod struct {
    Title  string `xml:"title,attr"`
    SubPods []SubPod `xml:"subpod"`
}

type SubPod struct {
    Text string `xml:"plaintext"`
}

// Function to initialize Wolfram Alpha API client
func initializeClient() string {
	// Load API key from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	apiKey := os.Getenv("WOLFRAM_FULL_RESPONSE")
	if apiKey == "" {
		panic("WOLFRAM_FULL_RESPONSE is not set in .env file")
	}
	return apiKey
}

// Function to perform a query
func queryWolfram(query string) (string, error) {
    appID := initializeClient()
    encodedQuery := url.QueryEscape(query)
    url := fmt.Sprintf("http://api.wolframalpha.com/v2/query?input=%s&appid=%s", encodedQuery, appID)

    // Make the HTTP request
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    // Log status and body for debugging
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %w", err)
    }

    // Parse XML response
    var result WolframResponse
    if err := xml.Unmarshal(body, &result); err != nil {
        return "", fmt.Errorf("error parsing XML: %w", err)
    }

    // Check if the query was successful and look for the "Result" pod first
    if result.Success {
        // Try to find the "Result" pod with the calculated answer
        for _, pod := range result.Pods {
            if pod.Title == "Result" && len(pod.SubPods) > 0 {
                return pod.SubPods[0].Text, nil
            }
        }
        // If "Result" pod is not available, fallback to any available pod with plaintext
        for _, pod := range result.Pods {
            if len(pod.SubPods) > 0 && pod.SubPods[0].Text != "" {
                return pod.SubPods[0].Text, nil
            }
        }
    }
    return "No result found.", nil
}

func main() {
	// Create the main application window
	myApp := app.New()
	myWindow := myApp.NewWindow("Wolfram Alpha")
	myWindow.Resize(fyne.NewSize(600, 300))

	// Create input field and result label
	inputField := widget.NewEntry()
	inputField.SetPlaceHolder("Enter your query here...")

	resultBox := widget.NewLabel("")
	resultBox.Wrapping = fyne.TextWrapWord

	// Function to handle button click
	queryButton := widget.NewButton("Query", func() {
		query := inputField.Text
		if query == "" {
			resultBox.SetText("Please enter a query.")
			return
		}
		answer, err := queryWolfram(query)
		if err != nil {
			resultBox.SetText("Error: Could not retrieve answer.")
		} else {
			resultBox.SetText(fmt.Sprint(answer))
		}
	})

	// Set up the layout
	content := container.NewVBox(
		widget.NewLabelWithStyle("Wolfram Query App", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		inputField,
		queryButton,
		widget.NewSeparator(),
		resultBox,
	)

	// Set up the Enter key press to trigger the query
	inputField.OnSubmitted = func(text string) {
		queryButton.OnTapped()
	}

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
