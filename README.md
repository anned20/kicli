# Kicli

This tool is a CLI interface to the [Kimai time tracking project](https://www.kimai.org/).

## Installation

You can install this tool in various ways:

### Using `go install`

If you have Go on your system, you can install this tool by using `go install github.com/anned20/kicli@latest`.

### Docker

To download the image use `docker pull ghcr.io/anned20/kicli:main`.

To use it, use `docker run -it -v $HOME/:/root/ --rm ghcr.io/anned20/kicli:main [args]`.

### Binary download

You can download a static binary from the latest release on Github.
Place this binary in your path and you can use it like any normal command.

## Usage

Before using this tool, you need to set up the used config variables.
You can do this by using `kicli setup`.
This command will interactively ask you for the required config variables.

For information about what the tool can do and how to use it, use `kicli help`.

## Development

### Commit message format

Commitizen is used to commit changes in this repository.
It can be installed using this snippet:

```bash
sudo pip3 install -U Commitizen
```

### Pre-commit hook

A pre-commit hook can be installed by using [pre-commit](https://pre-commit.com/).
Install this in your development environment using `pip install pre-commit`.

This pre-commit hook makes use of [golangci-lint](https://golangci-lint.run/) which can be installed using [these instructions](https://golangci-lint.run/usage/install/#local-installation).

Enable the pre-commit hook by executing the following snippet in the root of the repository:

```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

### TODO

- [ ] Better test coverage
