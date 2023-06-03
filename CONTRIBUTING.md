# Contributing

Thanks for taking the time to contribute! The following is a set of guidelines for contributing to our project.
We encourage everyone to follow them with their best judgement.

## Prerequisites

- [Go 1.20+](https://go.dev/):
    - Install on macOS with `brew install go`.
    - Install on Ubuntu with `sudo apt install golang`.
    - Install on Windows with [this link](https://go.dev/doc/install) or `choco install go` 
- Make
- Nix:
  - Install [Nix](https://zero-to-nix.com/start/install).

Development can be done using Make and Go on the host machine or using the Nix development shell.

## Setting Up Your Environment

1. Fork the repository on GitHub.
2. Clone your forked repository to your local machine.

```shell
 git clone https://github.com/tfadeyi/slotalk.git
```
3. Change directory to the cloned repository.

```shell
cd slotalk
```

## Local Development

Similar to other Go project, this repo make use of Make targets to build, test and lint the program.

> Note: run `make generate` before committing.

## Local Development with Nix

Start the nix development shell.

```shell
source env-dev.sh && develop
```

The shell currently uses the unstable nixpkgs, so it will use the latest version of Go and gotools,
it already comes installed with:
* golangci-lint
* goreleaser
* gomarkdoc
* aloe-cli

## Making Changes

1. Create a new branch for your changes.

```shell
git checkout -b <issue number>-<branch name>
```

2. Make your changes and commit them.

```shell
git commit --signoff
```

3. Push your changes to your forked repository.

```shell
git push origin <issue number>-<branch name>
```

4. Open a pull request on GitHub from your forked repository to the original repository.

## Code Review Process

All contributions will be reviewed by the maintainers of the project. Here are a few things to keep in mind:
* Please fill the given Pull Request template to the best of your abilities.
* Opening an issue before starting a work pieces improves the chances of the work being approved.

### Naming Conventions

For pull requests and branches a standard naming convention will help with automatically linking the development work with the related issue(s).
For this reason, please follow the following naming conventions:

* Branches: When creating a new branch the issue number should be added as a prefix `<issue number>-<branch-name>`
* Commits: The commit body should reference the issue `Related <[#issue number](issue URL)>`
