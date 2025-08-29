# Developer's Guide

This document explains how to build, test, and develop features for revive.

## Installation

Clone the project:

```bash
git clone git@github.com:mgechev/revive.git
cd revive
```

## Build

In order to build the project run:

```bash
make build
```

The command will produce the `revive` binary in the root of the project.

## Debug

To enable debug logging, set the `DEBUG` environment variable:

```sh
DEBUG=1 go run main.go
```

This will output debug information to `stderr` and to the log file `revive.log` created in the current working directory.

## Coding standards

Follow [the instructions](.github/instructions/) which contain Go coding standards and conventions used by both humans and GitHub Copilot.

## Development of rules

If you want to develop a new rule, follow as an example the already existing rules in the [rule package](https://github.com/mgechev/revive/tree/master/rule).

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```

All rules with a configuration must implement `lint.ConfigurableRule` interface:

```go
type ConfigurableRule interface {
	Configure(Arguments) error
}
```

The `Arguments` type is an alias of the type `[]any`. The arguments of the rule are passed from the configuration file.

### Example

Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with a given identifier.
We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
arguments = ["Foo"]
```

With the snippet above we:

- Enable the rule with the name `ban-struct-name`. The `Name()` method of our rule should return a string that matches `ban-struct-name`.
- Configure the rule with the argument `Foo`.
The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

A sample rule implementation can be [found here](/rule/argument_limit.go).

## Development of formatters

If you want to develop a new formatter, follow as an example the already existing formatters in the [formatter package](https://github.com/mgechev/revive/tree/master/formatter).

All formatters should implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```

## Lint

### Lint Markdown files

We use [markdownlint](https://github.com/DavidAnson/markdownlint),
[markdown-toc](https://github.com/jonschlinkert/markdown-toc),
and [mdsf](https://github.com/hougesen/mdsf) to check Markdown files.
`markdownlint` verifies document formatting, such as line length and empty lines.
`markdown-toc` checks the entries in the table of contents.
`mdsf` is responsible for formatting code snippets.

1. Install [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2#install).
2. Install [markdown-toc](https://github.com/jonschlinkert/markdown-toc#quick-start).
3. Install [mdsf](https://mdsf.mhouge.dk/#installation) and formatters:
    - [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) for `go`: `go install golang.org/x/tools/cmd/goimports@latest`
    - [shfmt](https://github.com/mvdan/sh#shfmt) for `sh, shell, bash`: `go install mvdan.cc/sh/v3/cmd/shfmt@latest`
    - [taplo](https://taplo.tamasfe.dev/cli/installation/binary.html) for `toml`
4. Run the following command to check formatting:

```shellsession
$ markdownlint-cli2 .
Finding: *.{md,markdown} *.md
Found:
 CODE_OF_CONDUCT.md
 CONTRIBUTING.md
 DEVELOPING.md
 README.md
 RULES_DESCRIPTIONS.md
Linting: 5 file(s)
Summary: 0 error(s)
```

_The `markdownlint-cli2` tool automatically uses the config file [.markdownlint-cli2.yaml](./.markdownlint-cli2.yaml)._
\
4. Run the following command to check TOC:

```sh
markdown-toc --maxdepth 4 --no-first1 --bullets "-" -i README.md && git diff --exit-code README.md
markdown-toc --maxdepth 2 --no-first1 --bullets "-" -i RULES_DESCRIPTIONS.md && git diff --exit-code RULES_DESCRIPTIONS.md
```

\
5. Run the following commands to verify and format code snippets:

```sh
mdsf verify .
```

```sh
mdsf format .
```

_Note: Use `golang` for Go code snippets that are intentionally non-compilable.
However, it is recommended to avoid this and use `go` whenever possible._
