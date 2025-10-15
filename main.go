package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const githubGraphQLEndpoint = "https://api.github.com/graphql"

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type Response struct {
	Data struct {
		User struct {
			Repositories struct {
				Nodes []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Url         string `json:"url"`
					UpdatedAt   string `json:"updatedAt"`
					Languages   struct {
						Nodes []struct {
							Name  string `json:"name"`
							Color string `json:"color"`
						}
					} `json:"languages"`
					Object struct {
						AbreviatedOid string `json:"abbreviatedOId"`
					} `json:"object"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

func fetchRepos(user string, token string) ([]byte, []string, error) {
	//create gql query to gather data on repository and load data from the resourcePath with /blob/

	query := `query($login: String!){
		user(login: $login) {
			repositories(first: 4, orderBy: {field: UPDATED_AT, direction: DESC}) {
				nodes {
					name
					description
					url
					updatedAt
					languages(first:4){
						nodes{
							name
							color
						}
					}
					object(expression: "main:README.md") {
						... on Blob {
							abbreviatedOid
						}
					}
				}
			}
		}
	}`

	reqBody := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"login": user,
		},
	}

	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return nil, nil, err
	}

	// fmt.Println(string(jsonBody))

	req, _ := http.NewRequest("POST", githubGraphQLEndpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Raw Response Body:", string(body))
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	// fmt.Println(data)
	//return the rawgithubcontent api for the readMeFiles to download
	var readMeFiles []string
	for _, repo := range data.Data.User.Repositories.Nodes {
		readMeFiles = append(readMeFiles, fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/main/README.md", user, repo.Name))
	}
	jsonOutput, _ := json.MarshalIndent(data.Data.User.Repositories.Nodes, "", "  ")
	return jsonOutput, readMeFiles, nil
}

//execute shell script to download files to folder

func main() {
	godotenv.Load(".env")
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	if token == "" {
		log.Fatal("please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}
	if username == "" {
		log.Fatal("please provide a GitHub Username via env variable GITHUB_USERNAME")
	}

	repos, readMeFiles, err := fetchRepos(username, token)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Run the shell script and pass the readme URLs as individual arguments.
	// We use slice expansion (readMeFiles...) to pass []string as variadic args.
	// Note: this expects a POSIX shell (sh) to be available. On Windows you may
	// need to run this under Git Bash / WSL or change to a Windows-compatible script.

	// Build the args slice: the first arg to the shell is the script path,
	// followed by the list of README URLs.
	// Then expand the whole args slice once into exec.Command's variadic parameter.
	args := append([]string{"./echoStrings"}, readMeFiles...)
	for i := range args {
		fmt.Println(args[i])
	}

	// Try to find a POSIX shell in PATH (bash preferred, then sh).
	// On Windows `sh` may not exist even if curl.exe is available; curl is a native
	// binary in PATH but shells like bash are provided by Git for Windows or WSL.
	var cmd *exec.Cmd
	shellCandidates := []string{"bash", "sh"}
	var shellPath string
	for _, s := range shellCandidates {
		if p, err := exec.LookPath(s); err == nil {
			shellPath = p
			break
		}
	}
	fmt.Println(shellPath)
	if shellPath != "" {
		cmd = exec.Command(shellPath, args...)
	} else {
		log.Println("no POSIX shell (bash/sh) found in PATH. If you want to run shell scripts on Windows, install Git for Windows (add Git Bash to PATH) or use WSL. Attempting to execute the script directly...")
		cmd = exec.Command("./echoStrings", readMeFiles...)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("script execution failed: %v\noutput: %s", err, string(out))
	} else {
		fmt.Printf("script output: %s\n", string(out))
	}

	f, err := os.Create("repos.json")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer f.Close()

	_, err = f.Write(repos)

	if err != nil {
		defer f.Close()
		log.Fatalf("Error: %v", err)
	}
	// fmt.Printf("List of repositories: %v", repos)

}
