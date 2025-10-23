package globals

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

// check if frontend can support object like syntax in markdown metaData
type LanguageData struct {
	Name  string
	Color string
	Size  int
}

type RepoMetaData struct {
	Description string
	Url         string
	UpdatedAt   string
	Languages   []LanguageData
	ReadMeOid   string
}

var (
	ReposData      *Response //might end up deleting
	ReposMetaData  *map[string]RepoMetaData
	DestinationDir string
)
