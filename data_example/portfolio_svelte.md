---
title: portfolio_svelte
description: An attempt of porting the porfolio project to Svelte in order to get first hand experience in porting a project to another tech stack and understand more about Svelte as a framework
label: building
url: https://github.com/zidariu-sabin/portfolio_svelte
updatedat: "2025-12-10T21:43:11Z"
languages:
- name: Svelte
  color: "#ff3e00"
  size: 72.6
- name: JavaScript
  color: "#f1e05a"
  size: 17.7
- name: CSS
  color: "#663399"
  size: 6.9
- name: HTML
  color: "#e34c26"
  size: 2.9
docfileoid: 75842c4
---
# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project in the current directory
npx sv create

# create a new project in my-app
npx sv create my-app
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
