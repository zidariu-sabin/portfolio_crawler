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

func extractDataFromResponse(user string, docFile string, data *globals.Response) (*map[string]string, error) {
	//return the rawgithubcontent api for the repos that have docFiles to download
	var docFileLinks = make(map[string]string)
	var reposMetaData = make(map[string]globals.RepoMetaData)
	var rmts = make([]globals.RepoMetaData, 0)
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

		rmt := globals.RepoMetaData{
			Title:       repo.Name,
			Description: repo.Description,
			Label:       "building", // temporary hardcoded label
			Url:         repo.Url,
			UpdatedAt:   repo.UpdatedAt,
			Languages:   langs,
			DocFileOid:  repo.Object.AbreviatedOid,
		}

		rmts = append(rmts, rmt)

		reposMetaData[repo.Name] = rmt

		docFileLinks[repo.Name] = fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/main/%s.md", user, repo.Name, docFile)
	}

	_, err := utils.GenerateJsonOfAllMetaData(globals.JsonFileDesiredDir+"/reposMetaData.json", rmts)

	if err != nil {
		return nil, err
	}

	globals.ReposData = data
	globals.ReposMetaData = reposMetaData

	return &docFileLinks, nil
}

func fetchRepos(user string, token string, docFile string) (map[string]string, error) {
	//create gql query to gather data on repository and load data from the resourcePath with /blob/

	query := fmt.Sprintf(`query($login: String!){
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
					object(expression: "main:%s.md") {
						... on Blob {
							abbreviatedOid
						}
					}
				}
			}
		}
	}`, docFile)

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

	docFileLinks, err := extractDataFromResponse(user, docFile, &data)

	if err != nil {
		return nil, err
	}

	return *docFileLinks, nil
}

//execute shell script to download files to folder

func main() {
	godotenv.Load(".env")
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	user := os.Getenv("GITHUB_USERNAME")
	docFile := os.Getenv("DOCUMENTATION_FILE_NAME")

	if token == "" {
		log.Fatal("Error: please provide a GitHub API token via env variable GITHUB_AUTH_TOKEN")
	}
	if user == "" {
		log.Fatal("Error: please provide a GitHub Username via env variable GITHUB_USERNAME")
	}

	if docFile == "" {
		log.Fatal("Error: please provide a docFile name to pull via env variable DOCUMENTATION_FILE_NAME")
	}

	mdFilesDesiredDir, err := utils.ValidateDest(os.Getenv("MARKDOWN_FILES_DESTINATION_FOLDER_PATH"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	jsonFileDesiredDir, err := utils.ValidateDest(os.Getenv("JSON_DESTINATION_FOLDER_PATH"))

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	globals.DestinationDir = mdFilesDesiredDir
	globals.JsonFileDesiredDir = jsonFileDesiredDir

	docFileLinks, err := fetchRepos(user, token, docFile)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	reposJsonOutput, err := json.MarshalIndent(globals.ReposData.Data.User.Repositories.Nodes, "", "  ")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	utils.DownloadMany(docFileLinks, 3)

	utils.WriteJson("/repos.json", reposJsonOutput)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// utils.WriteYaml(globals.ReposMetaData)
	// fmt.Printf("List of metadata: %v", globals.ReposMetaData)

}
