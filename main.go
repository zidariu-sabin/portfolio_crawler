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
	"portfolio_crawler/globals"
	"portfolio_crawler/utils"
	"sort"
	"strconv"

	"github.com/joho/godotenv"
)

const githubGraphQLEndpoint = "https://api.github.com/graphql"

func extractDataFromResponse(user string, data *globals.Response) (*map[string]string, error) {
	//return the rawgithubcontent api for the repos that have readMeFiles to download
	var readMeFileLinks = make(map[string]string)
	var reposMetaData = make(map[string]globals.RepoMetaData)
	for _, repo := range data.Data.User.Repositories.Nodes {
		if repo.Object.AbreviatedOid == "" {
			continue
		}
		edges := repo.Languages.Edges
		totalSize := float32(0)
		for _, edge := range edges {
			totalSize += float32(edge.Size)
		}

		langs := make([]globals.LanguageData, 0)
		for _, edge := range edges {
			percentage, _ := strconv.ParseFloat(
				fmt.Sprintf("%.1f", float32(edge.Size)/totalSize*100),
				32,
			)

			langs = append(langs, globals.LanguageData{
				Name:  edge.Node.Name,
				Color: edge.Node.Color,
				Size:  float32(percentage),
			})
		}

		sort.Slice(langs, func(i, j int) bool {
			return langs[i].Size > langs[j].Size
		})

		reposMetaData[repo.Name] = globals.RepoMetaData{
			Title:       repo.Name,
			Description: repo.Description,
			Label:       "building", // temporary hardcoded label
			Url:         repo.Url,
			UpdatedAt:   repo.UpdatedAt,
			Languages:   langs,
			ReadMeOid:   repo.Object.AbreviatedOid,
		}

		readMeFileLinks[repo.Name] = fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/main/README.md", user, repo.Name)
	}

	globals.ReposData = data
	globals.ReposMetaData = reposMetaData

	return &readMeFileLinks, nil
}

func fetchRepos(user string, token string) (map[string]string, error) {
	//create gql query to gather data on repository and load data from the resourcePath with /blob/

	query := `query($login: String!){
		user(login: $login) {
			repositories(first: 100, orderBy: {field: UPDATED_AT, direction: DESC}) {
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

	reqBody := globals.GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"login": user,
		},
	}

	jsonBody, err := json.Marshal(reqBody)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", githubGraphQLEndpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data globals.Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err

	}

	readMeFileLinks, err := extractDataFromResponse(user, &data)

	if err != nil {
		return nil, err
	}

	return *readMeFileLinks, nil
}

//execute shell script to download files to folder

func main() {
	godotenv.Load(".env")
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	user := os.Getenv("GITHUB_USERNAME")

	if token == "" {
		log.Fatal("Error: please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}
	if user == "" {
		log.Fatal("Error: please provide a GitHub Username via env variable GITHUB_USERNAME")
	}

	desiredDir, err := utils.ValidateDest(os.Getenv("DESTINATION_FOLDER_PATH"))

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	globals.DestinationDir = desiredDir
	// globals.ReposMetaData = make(map[string]globals.RepoMetaData)

	readMeFileLinks, err := fetchRepos(user, token)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	reposJsonOutput, err := json.MarshalIndent(globals.ReposData.Data.User.Repositories.Nodes, "", "  ")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	utils.DownloadMany(readMeFileLinks, 3)

	utils.WriteJson(reposJsonOutput)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// utils.WriteYaml(globals.ReposMetaData)
	// fmt.Printf("List of metadata: %v", globals.ReposMetaData)

}
