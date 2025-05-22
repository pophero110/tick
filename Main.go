package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Task struct {
	Title      string    `json:"title"`
	Kind       string    `json:"kind"`
	Content    string    `json:"content,omitempty"`
	Desc       string    `json:"desc,omitempty"`
	IsAllDay   bool      `json:"isAllDay,omitempty"`
	StartDate  string    `json:"startDate,omitempty"`
	DueDate    string    `json:"dueDate,omitempty"`
	TimeZone   string    `json:"timeZone,omitempty"`
	Reminders  []string  `json:"reminders,omitempty"`
	RepeatFlag string    `json:"repeatFlag,omitempty"`
	Priority   int       `json:"priority,omitempty"`
	SortOrder  int       `json:"sortOrder,omitempty"`
	Items      []Subtask `json:"items,omitempty"`
	Tags       []string  `json:"tags,omitempty"`
}

type Subtask struct {
	Title         string `json:"title"`
	StartDate     string `json:"startDate,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	SortOrder     int    `json:"sortOrder,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
	Status        int    `json:"status,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
}

func createTask(kind string, title string) {
	accessToken := os.Getenv("TICKTOKEN")

	// task := Task{
	// 	Title:     "Write CLI in Go",
	// 	Content:   "Create a CLI using TickTick Open API",
	// 	Desc:      "Details about the Go implementation.",
	// 	IsAllDay:  false,
	// 	StartDate: time.Now().Format("2006-01-02T15:04:05-0700"),
	// 	DueDate:   time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05-0700"),
	// 	TimeZone:  "UTC",
	// 	Priority:  1,
	// 	Items: []Subtask{
	// 		{Title: "Define structure", SortOrder: 1},
	// 		{Title: "Write main logic", SortOrder: 2},
	// 	},
	// }

	task := Task{
		Title: title,
		Kind:  kind,
		Tags:   []string{"tick"},
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.ticktick.com/open/v1/task", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(body))
}

func auth() {
	clientID := os.Getenv("TICKCLIENTID")
	clientSecret := os.Getenv("TICKCLIENTSECRET")
	scope := "tasks:read tasks:write"
	redirectURI := "http://localhost"
	state := "xyz"

	authURL := fmt.Sprintf("https://ticktick.com/oauth/authorize?scope=%s&client_id=%s&state=%s&redirect_uri=%s&response_type=code",
		url.QueryEscape(scope), clientID, state, url.QueryEscape(redirectURI))

	fmt.Println("Please open the following URL in your browser and authorize the application:")
	fmt.Println(authURL)

	fmt.Print("Enter the authorization code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	code := scanner.Text()

	tokenURL := "https://ticktick.com/oauth/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", scope)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:")
	fmt.Println(string(body))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tick <command>")
		fmt.Println("tick task <Title>: Create a new task")
		fmt.Println("tick note <Title>: Create a new note")
		fmt.Println("tick auth: Authentication")
		return
	}

	switch os.Args[1] {
	case "task":
		createTask("TASK", os.Args[2])
	case "note":
		createTask("NOTE", os.Args[2])
	case "auth":
		auth()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
