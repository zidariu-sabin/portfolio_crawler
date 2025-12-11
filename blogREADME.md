# PortfolioCrawler

## Description
This project represents a Golang script that allows for pulling data about the documented repositories present on a Github Profile in order to display in a portfolio.

## How it's done

### Filtering the Repositories
 An http post request is sent to the GitHub Graphql api using the personal access token for authentication in order to gather some metadata like the name, description and languages used, together with the OId of a specifically mentioned markdown file. If the OId is present, it means the repository is documented and so will be included in the gathered data.
 
### Gathering the Data
 Next, with the filtered repositories in mind, it is going to download the content of the specified doc files asyncroniously using channels, through the raw.githubusercontent domain and recreate each markdown file, named by Repository, with the injected metadata as yaml front matter and the current content.
 
## How I use it

So far I have used this script to create a portfolio section in my website.

Using a markdown parser and a file indexer, the generated files can be used to create components and display the markdown elements after a custom designed style as both separate pages and inside components for a content display

So far I have reacreated this behaviour in Nuxt using the Nuxt Content Module and Svelte using Mdsvex.

I have created a blogREADME.md file in each repository that i want to be displayed in my portfolio, and used this script to pull the data and create the markdown files with the metadata.

In this case, I set the env variable DOCUMENTATION_FILE_NAME to "blogREADME" in order to pull those files.

Feel Free to check out in detail how i used the files for my portfolio project:
[Nuxt](https://github.com/zidariu-sabin/nuxt_portfolio), [Svelte](https://github.com/zidariu-sabin/portfolio_svelte)
 