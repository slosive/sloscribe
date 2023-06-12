---
sidebar_position: 1
---

# Tutorial Intro

Let's discover **Slotalk in less than 5 minutes**.

## Getting Started

Get started by **installing slotalk** .


## Add slotalk annotations to your codebase
1. Add comments to your source code. See [Declarative Comments](annotations/sloth/service).

## Generate Sloth Definitions from codebase
2. Run `slotalk` init in the project's root. This will parse your source code annotations and print the sloth definitions to standard out.
    ```shell
    ./slotalk init
    ```

   You can also specify the specific file to parse by using the `-f` flag.

    ```shell
    ./slotalk init -f metrics.go
    ```

   Another way would be to pass the input file through pipe.

    ```shell
    cat metrics.go | ./slotalk init -f -
    ```
   
## Generate a new site

Generate a new Docusaurus site using the **classic template**.

The classic template will automatically be added to your project after you run the command:

```bash
npm init docusaurus@latest my-website classic
```

You can type this command into Command Prompt, Powershell, Terminal, or any other integrated terminal of your code editor.

The command also installs all necessary dependencies you need to run Docusaurus.

## Start your site

Run the development server:

```bash
cd my-website
npm run start
```

The `cd` command changes the directory you're working with. In order to work with your newly created Docusaurus site, you'll need to navigate the terminal there.

The `npm run start` command builds your website locally and serves it through a development server, ready for you to view at http://localhost:3000/.

Open `docs/intro.md` (this page) and edit some lines: the site **reloads automatically** and displays your changes.
