package main

import "fmt"

func main() {
	str := "query($login: String!){\n\t\tuser(login: $login) {\n\t\t\trepositories(first: 4, orderBy: {field: UPDATED_AT, direction: DESC}) {\n\t\t\t\tnodes {\n\t\t\t\t\tname\n\t\t\t\t\tdescription\n\t\t\t\t\turl\n\t\t\t\t\tupdatedAt\n\t\t\t\t\tlanguages(first:4){\n\t\t\t\t\t\tnodes{\n\t\t\t\t\t\t\tname\n\t\t\t\t\t\t\tcolor\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t\tobject(expression: \"main:README.md\") {\n\t\t\t\t\t\t... on Blob {\n\t\t\t\t\t\t\tabbreviatedOId\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}"

	fmt.Println(str)
}
