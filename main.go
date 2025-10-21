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
	"portfolio_crawler/utils"

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
						Edges []struct {
							Node struct {
								Name  string `json:"name"`
								Color string `json:"color"`
							} `json:"node"`
							Size int `json:"size"`
						} `json:"edges"`
					} `json:"languages"`
					Object struct {
						AbreviatedOid string `json:"abbreviatedOId"`
					} `json:"object"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

func fetchRepos(user string, token string) ([]byte, map[string]string, error) {
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
						edges{
             			node{
							name
							color
						}
						size
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
	var readMeFileLinks = make(map[string]string)
	for _, repo := range data.Data.User.Repositories.Nodes {
		if repo.Object.AbreviatedOid == "" {
			continue
		}
		// readMeFiles = append(readMeFiles, fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/main/README.md", user, repo.Name))
		readMeFileLinks[repo.Name] = fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/main/README.md", user, repo.Name)
	}
	jsonOutput, _ := json.MarshalIndent(data.Data.User.Repositories.Nodes, "", "  ")
	return jsonOutput, readMeFileLinks, nil
}

//execute shell script to download files to folder

func main() {
	godotenv.Load(".env")
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	destFile := os.Getenv("DESTINATION_FOLDER_PATH")
	if token == "" {
		log.Fatal("please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}
	if username == "" {
		log.Fatal("please provide a GitHub Username via env variable GITHUB_USERNAME")
	}

	repos, readMeFileLinks, err := fetchRepos(username, token)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	utils.DownloadMany(readMeFileLinks, 3, destFile)

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
