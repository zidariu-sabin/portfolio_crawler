# PortfolioCrawler

## Description
This project represents a Golang script that allows for pulling data about the documented repositories present on a Github Profile in order to display in a portfolio.

## How it's done

### Filtering the Repositories
 An http post request is sent to the GitHub Graphql api using the personal access token for authentication in order to gather some metadata like the name, description and languages used, together with the OId of a specifically created markdown file(the blogREADME.md). If the OId is present, it means the repository is documented and so will be included in the portfolio.
 
### Gathering the Data
 Next, with the filtered repositories in mind, we are going to download the content of the blogREADME.md files asyncroniously using channels, through the raw.githubusercontent domain and recreate each file markdown file, named by Repository, with the injected metadata as yaml front matter and paste the content of the blogReadmMe.
 
## How I use it

So far I have used this script to create a portfolio section in my website.

Using a markdown parser, and a file indexer, the files can be used to create components and display the markdown elements after a custom designed style as both separate documented pages and inside components for a content display

So far i have reacreated this behaviour in Nuxt using the Nuxt Content Module and Svelte using Mdsvex.

Feel Free to check out how i did it :
[Nuxt](https://github.com/zidariu-sabin/nuxt_portfolio), [Svelte](https://github.com/zidariu-sabin/portfolio_svelte)
 