[![Build Status](https://github.com/mgechev/revive/actions/workflows/test.yaml/badge.svg)](https://github.com/mgechev/revive/actions/workflows/test.yaml)

# revive

Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint. **`Revive` provides a framework for development of custom rules, and lets you define a strict preset for enhancing your development & code review processes**.

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
  <br>
  Logo by <a href="https://github.com/hawkgs">Georgi Serev</a>
</p>

Here's how `revive` is different from `golint`:

- Allows to enable or disable rules using a configuration file.
- Allows to configure the linting rules with a TOML file.
- 2x faster running the same rules as golint.
- Provides functionality for disabling a specific rule or the entire linter for a file or a range of lines.
  - `golint` allows this only for generated files.
- Optional type checking. Most rules in golint do not require type checking. If you disable them in the config file, revive will run over 6x faster than golint.
- Provides multiple formatters which let us customize the output.
- Allows to customize the return code for the entire linter or based on the failure of only some rules.
- _Everyone can extend it easily with custom rules or formatters._
- `Revive` provides more rules compared to `golint`.

## Who uses Revive

- [`tidb`](https://github.com/pingcap/tidb) - TiDB is a distributed HTAP database compatible with the MySQL protocol
- [`grafana`](https://github.com/grafana/grafana) - The tool for beautiful monitoring and metric analytics & dashboards for Graphite, InfluxDB & Prometheus & More
- [`etcd`](https://github.com/etcd-io/etcd) - Distributed reliable key-value store for the most critical data of a distributed system
- [`cadence`](https://github.com/uber/cadence) - Cadence is a distributed, scalable, durable, and highly available orchestration engine by Uber to execute asynchronous long-running business logic in a scalable and resilient way
- [`ferret`](https://github.com/MontFerret/ferret) - Declarative web scraping
- [`gopass`](https://github.com/gopasspw/gopass) - The slightly more awesome standard unix password manager for teams
- [`gitea`](https://github.com/go-gitea/gitea) - Git with a cup of tea, painless self-hosted git service
- [`excelize`](https://github.com/360EntSecGroup-Skylar/excelize) - Go library for reading and writing Microsoft Excel™ (XLSX) files
- [`aurora`](https://github.com/xuri/aurora) - aurora is a web-based Beanstalk queue server console written in Go
- [`soar`](https://github.com/XiaoMi/soar) - SQL Optimizer And Rewriter
- [`pyroscope`](https://github.com/pyroscope-io/pyroscope) - Continuous profiling platform
- [`gorush`](https://github.com/appleboy/gorush) - A push notification server written in Go (Golang).
- [`dry`](https://github.com/moncho/dry) - dry - A Docker manager for the terminal.
- [`go-echarts`](https://github.com/chenjiandongx/go-echarts) - The adorable charts library for Golang
- [`reviewdog`](https://github.com/reviewdog/reviewdog) - Automated code review tool integrated with any code analysis tools regardless of programming language
- [`rudder-server`](https://github.com/rudderlabs/rudder-server) - Privacy and Security focused Segment-alternative, in Golang and React.
- [`sklearn`](https://github.com/pa-m/sklearn) - A partial port of scikit-learn written in Go.
- [`protoc-gen-doc`](https://github.com/pseudomuto/protoc-gen-doc) - Documentation generator plugin for Google Protocol Buffers.
- [`llvm`](https://github.com/llir/llvm) - Library for interacting with LLVM IR in pure Go.
- [`jenkins-library`](https://github.com/SAP/jenkins-library) - Jenkins shared library for Continuous Delivery pipelines by SAP.
- [`pd`](https://github.com/tikv/pd) - Placement driver for TiKV.
- [`shellhub`](https://github.com/shellhub-io/shellhub) - ShellHub enables teams to easily access any Linux device behind firewall and NAT.
- [`lorawan-stack`](https://github.com/TheThingsNetwork/lorawan-stack) - The Things Network Stack for LoRaWAN V3
- [`gin-jwt`](https://github.com/appleboy/gin-jwt) - This is a JWT middleware for Gin framework.
- [`gofight`](https://github.com/appleboy/gofight) - Testing API Handler written in Golang.
- [`Beaver`](https://github.com/Clivern/Beaver) - A Real Time Messaging Server.
- [`ggz`](https://github.com/go-ggz/ggz) - An URL shortener service written in Golang
- [`Codeac.io`](https://www.codeac.io?ref=revive) - Automated code review service integrates with GitHub, Bitbucket and GitLab (even self-hosted) and helps you fight technical debt.
- [`DevLake`](https://github.com/apache/incubator-devlake) - Apache DevLake is an open-source dev data platform to ingest, analyze, and visualize the fragmented data from DevOps tools，which can distill insights to improve engineering productivity.
- [`checker`](https://github.com/cinar/checker) - Checker helps validating user input through rules defined in struct tags or directly through functions.
- [`milvus`](https://github.com/milvus-io/milvus) - A cloud-native vector database, storage for next generation AI applications.

*Open a PR to add your project*.

<p align="center">
  <img src="./assets/demo.svg" alt="" width="700">
</p>

<!-- TOC -->

- [revive](#revive)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Docker](#docker)
    - [Bazel](#bazel)
    - [Text Editors](#text-editors)
    - [GitHub Actions](#github-actions)
    - [Continuous Integration](#continuous-integration)
    - [Linter Aggregators](#linter-aggregators)
      - [golangci-lint](#golangci-lint)
    - [Command Line Flags](#command-line-flags)
    - [Sample Invocations](#sample-invocations)
    - [Comment Directives](#comment-directives)
    - [Configuration](#configuration)
    - [Default Configuration](#default-configuration)
    - [Custom Configuration](#custom-configuration)
    - [Recommended Configuration](#recommended-configuration)
  - [Available Rules](#available-rules)
  - [Configurable rules](#configurable-rules)
    - [`var-naming`](#var-naming)
  - [Available Formatters](#available-formatters)
    - [Friendly](#friendly)
    - [Stylish](#stylish)
    - [Default](#default)
    - [Plain](#plain)
    - [Unix](#unix)
    - [SARIF](#sarif)
  - [Extensibility](#extensibility)
    - [Custom Rule](#writing-a-custom-rule)
      - [Example](#example)
    - [Custom Formatter](#custom-formatter)
  - [Speed Comparison](#speed-comparison)
    - [golint](#golint)
    - [revive](#revive)
  - [Overriding colorization detection](#overriding-colorization-detection)
  - [Contributors](#contributors)
  - [License](#license)

<!-- /TOC -->

## Installation

```bash
go install github.com/mgechev/revive@latest
```

or get a released executable from the [Releases](https://github.com/mgechev/revive/releases) page.

You can install the main branch (including the last commit) with:
```bash
go install github.com/mgechev/revive@master
```

## Usage

Since the default behavior of `revive` is compatible with `golint`, without providing any additional flags, the only difference you'd notice is faster execution.

`revive` supports a `-config` flag whose value should correspond to a TOML file describing which rules to use for `revive`'s linting. If not provided, `revive` will try to use a global config file (assumed to be located at `$HOME/revive.toml`). Otherwise, if no configuration TOML file is found then `revive` uses a built-in set of default linting rules.

### Docker
A volume must be mounted to share the current repository with the container.
Please refer to the [bind mounts Docker documentation](https://docs.docker.com/storage/bind-mounts/)

```bash
docker run -v "$(pwd)":/var/<repository> ghcr.io/mgechev/revive:v1.3.7 -config /var/<repository>/revive.toml -formatter stylish ./var/kidle/...
```

- `-v` is for the volume
- `ghcr.io/mgechev/revive:v1.3.7 ` is the image name and its version corresponds to `revive` command
- The provided flags are the same as the binary usage.

### Bazel

If you want to use revive with Bazel, look at the [rules](https://github.com/atlassian/bazel-tools/tree/master/gorevive) that Atlassian maintains.

### Text Editors

- Support for VSCode in [vscode-go](https://github.com/Microsoft/vscode-go/pull/1699).
- Support for GoLand via [File Watchers](https://dev.to/s0xzwasd/configure-revive-go-linter-in-goland-2ggl).
- Support for Atom via [linter-revive](https://github.com/morphy2k/linter-revive).
- Support for vim via [dense-analysis/ale](https://github.com/dense-analysis/ale).

  ```vim
  let g:ale_linters = {
  \   'go': ['revive'],
  \}
  ```

- Support for Neovim via [null-ls.nvim](https://github.com/jose-elias-alvarez/null-ls.nvim).

  ```lua
  require("null-ls").setup({
      sources = {
          require("null-ls").builtins.diagnostics.revive
      },
  })
  ```

### GitHub Actions

- [Revive Action](https://github.com/marketplace/actions/revive-action) with annotation support

### Continuous Integration

[Codeac.io](https://www.codeac.io?ref=revive) - Automated code review service integrates with GitHub, Bitbucket and GitLab (even self-hosted) and helps you fight technical debt. Check your [pull-requests](https://www.codeac.io/documentation/pull-requests.html?ref=revive) with [revive](https://www.codeac.io/documentation/revive-configuration.html?ref=revive) automatically. (free for open-source projects)

### Linter aggregators

#### golangci-lint
To enable `revive` in `golangci-lint` you need to add `revive` to the list of enabled linters:

```yaml
# golangci-lint configuration file
linters:
   enable:
     - revive
```
Then `revive` can be configured by adding an entry to the `linters-settings` section of the configuration, for example:

```yaml
# golangci-lint configuration file
linters-settings:
  revive:
    ignore-generated-header: true
    severity: warning
    rules:
      - name: atomic
      - name: line-length-limit
        severity: error
        arguments: [80]
      - name: unhandled-error
        arguments : ["fmt.Printf", "myFunction"]
```

The above configuration enables three rules of `revive`: _atomic_, _line-length-limit_ and _unhandled-error_ and pass some arguments to the last two.
The [Configuration](#configuration) section of this document provides details on how to configure `revive`. Note that while `revive` configuration is in TOML, that of `golangci-lint` is in YAML.

Please notice that if no particular configuration is provided, `revive` will behave as `go-lint` does, i.e. all `go-lint` rules are enabled (the [Available Rules table](#available-rules) details what are the `go-lint` rules). When a configuration is provided, only rules in the configuration are enabled.

### Command Line Flags

`revive` accepts the following command line parameters:

- `-config [PATH]` - path to the config file in TOML format, defaults to `$HOME/revive.toml` if present.
- `-exclude [PATTERN]` - pattern for files/directories/packages to be excluded for linting. You can specify the files you want to exclude for linting either as package name (i.e. `github.com/mgechev/revive`), list them as individual files (i.e. `file.go`), directories (i.e. `./foo/...`), or any combination of the three.
- `-formatter [NAME]` - formatter to be used for the output. The currently available formatters are:

  - `default` - will output the failures the same way that `golint` does.
  - `json` - outputs the failures in JSON format.
  - `ndjson` - outputs the failures as a stream in newline delimited JSON (NDJSON) format.
  - `friendly` - outputs the failures when found. Shows the summary of all the failures.
  - `stylish` - formats the failures in a table. Keep in mind that it doesn't stream the output so it might be perceived as slower compared to others.
  - `checkstyle` - outputs the failures in XML format compatible with that of Java's [Checkstyle](https://checkstyle.org/).
- `-max_open_files` -  maximum number of open files at the same time. Defaults to unlimited.
- `-set_exit_status` - set exit status to 1 if any issues are found, overwrites `errorCode` and `warningCode` in config.
- `-version` - get revive version.


### Sample Invocations

```shell
revive -config revive.toml -exclude file1.go -exclude file2.go -formatter friendly github.com/mgechev/revive package/...
```

- The command above will use the configuration from `revive.toml`
- `revive` will ignore `file1.go` and `file2.go`
- The output will be formatted with the `friendly` formatter
- The linter will analyze `github.com/mgechev/revive` and the files in `package`

### Comment Directives

Using comments, you can disable the linter for the entire file or only a range of lines:

```go
//revive:disable

func Public() {}
//revive:enable
```

The snippet above, will disable `revive` between the `revive:disable` and `revive:enable` comments. If you skip `revive:enable`, the linter will be disabled for the rest of the file.

With `revive:disable-next-line` and `revive:disable-line` you can disable `revive` on a particular code line.

You can do the same on a rule level. In case you want to disable only a particular rule, you can use:

```go
//revive:disable:unexported-return
func Public() private {
  return private
}
//revive:enable:unexported-return
```

This way, `revive` will not warn you that you're returning an object of an unexported type, from an exported function.

You can document why you disable the linter by adding a trailing text in the directive, for example

```go
//revive:disable Until the code is stable
```
```go
//revive:disable:cyclomatic High complexity score but easy to understand
```

You can also configure `revive` to enforce documenting linter disabling directives by adding

```toml
[directive.specify-disable-reason]
```

in the configuration. You can set the severity (defaults to _warning_) of the violation of this directive

```toml
[directive.specify-disable-reason]
    severity = "error"
```

### Configuration

`revive` can be configured with a TOML file. Here's a sample configuration with an explanation of the individual properties:

```toml
# When set to false, ignores files with "GENERATED" header, similar to golint
ignoreGeneratedHeader = true

# Sets the default severity to "warning"
severity = "warning"

# Sets the default failure confidence. This means that linting errors
# with less than 0.8 confidence will be ignored.
confidence = 0.8

# Sets the error code for failures with the "error" severity
errorCode = 0

# Sets the error code for failures with severity "warning"
warningCode = 0

# Configuration of the `cyclomatic` rule. Here we specify that
# the rule should fail if it detects code with higher complexity than 10.
[rule.cyclomatic]
  arguments = [10]

# Sets the severity of the `package-comments` rule to "error".
[rule.package-comments]
  severity = "error"
```

By default `revive` will enable only the linting rules that are named in the configuration file.
For example, the previous configuration file makes `revive` to enable only _cyclomatic_ and _package-comments_ linting rules.

To enable all available rules you need to add:

```toml
enableAllRules = true
```

This will enable all available rules no matter what rules are named in the configuration file.

To disable a rule, you simply mark it as disabled in the configuration.
For example:

```toml
[rule.line-length-limit]
    Disabled = true
```
When enabling all rules you still need/can provide specific configurations for rules.
The following file is an example configuration where all rules are enabled, except for those that are explicitly disabled, and some rules are configured with particular arguments:

```toml
severity = "warning"
confidence = 0.8
errorCode = 0
warningCode = 0

# Enable all available rules
enableAllRules = true

# Disabled rules
[rule.blank-imports]
    Disabled = true
[rule.file-header]
    Disabled = true
[rule.max-public-structs]
    Disabled = true
[rule.line-length-limit]
    Disabled = true
[rule.function-length]
    Disabled = true
[rule.banned-characters]
    Disabled = true

# Rule tuning
[rule.argument-limit]
    Arguments = [5]
[rule.cyclomatic]
    Arguments = [10]
[rule.cognitive-complexity]
    Arguments = [7]
[rule.function-result-limit]
    Arguments = [3]
[rule.error-strings]
    Arguments = ["mypackage.Error"]
```

### Default Configuration

The default configuration of `revive` can be found at `defaults.toml`. This will enable all rules available in `golint` and use their default configuration (i.e. the way they are hardcoded in `golint`).

```shell
revive -config defaults.toml github.com/mgechev/revive
```

This will use the configuration file `defaults.toml`, the `default` formatter, and will run linting over the `github.com/mgechev/revive` package.

### Custom Configuration

```shell
revive -config config.toml -formatter friendly github.com/mgechev/revive
```

This will use `config.toml`, the `friendly` formatter, and will run linting over the `github.com/mgechev/revive` package.

### Recommended Configuration

The following snippet contains the recommended `revive` configuration that you can use in your project:

```toml
ignoreGeneratedHeader = false
severity = "warning"
confidence = 0.8
errorCode = 0
warningCode = 0

[rule.blank-imports]
[rule.context-as-argument]
[rule.context-keys-type]
[rule.dot-imports]
[rule.error-return]
[rule.error-strings]
[rule.error-naming]
[rule.exported]
[rule.increment-decrement]
[rule.var-naming]
[rule.var-declaration]
[rule.package-comments]
[rule.range]
[rule.receiver-naming]
[rule.time-naming]
[rule.unexported-return]
[rule.indent-error-flow]
[rule.errorf]
[rule.empty-block]
[rule.superfluous-else]
[rule.unused-parameter]
[rule.unreachable-code]
[rule.redefines-builtin-id]
```

### Rule-level file excludes

You also can setup custom excludes for each rule.

It's an alternative for the global `-exclude` program arg.

```toml
ignoreGeneratedHeader = false
severity = "warning"
confidence = 0.8
errorCode = 0
warningCode = 0

[rule.blank-imports]
   Exclude=["**/*.pb.go"]
[rule.context-as-argument]
   Exclude=["src/somepkg/*.go", "TEST"]
```

You can use the following exclude patterns

1. full paths to files `src/pkg/mypkg/some.go`
2. globs `src/**/*.pb.go`
3. regexes (should have prefix ~) `~\.(pb|auto|generated)\.go$`
4. well-known `TEST` (same as `**/*_test.go`)
5. special cases:
  a. `*` and `~` patterns exclude all files (same effect as disabling the rule)
  b. `""` (empty) pattern excludes nothing

> NOTE: do not mess with `exclude` that can  be used at the top level of TOML file, that means "exclude package patterns", not "exclude file patterns"

## Available Rules

List of all available rules. The rules ported from `golint` are left unchanged and indicated in the `golint` column.
| Name                  | Config | Description                                                      | `golint` | Typed |
| --------------------- | :----: | :--------------------------------------------------------------- | :------: | :---: |
| [`context-keys-type`](./RULES_DESCRIPTIONS.md#context-key-types)   |  n/a   | Disallows the usage of basic types in `context.WithValue`.       |   yes    |  yes  |
| [`time-equal`](./RULES_DESCRIPTIONS.md#time-equal)         |  n/a   | Suggests to use `time.Time.Equal` instead of `==` and `!=` for equality check time.                 |   no    |  yes  |
| [`time-naming`](./RULES_DESCRIPTIONS.md#time-naming)         |  n/a   | Conventions around the naming of time variables.                 |   yes    |  yes  |
| [`unchecked-type-assertions`](./RULES_DESCRIPTIONS.md#unchecked-type-assertions)         |  n/a   | Disallows type assertions without checking the result.                 |   no    |  yes  |
| [`var-declaration`](./RULES_DESCRIPTIONS.md#var-declaration)     |  n/a   | Reduces redundancies around variable declaration.                |   yes    |  yes  |
| [`unexported-return`](./RULES_DESCRIPTIONS.md#unexported-return)   |  n/a   | Warns when a public return is from unexported type.              |   yes    |  yes  |
| [`errorf`](./RULES_DESCRIPTIONS.md#errorf)              |  n/a   | Should replace `errors.New(fmt.Sprintf())` with `fmt.Errorf()`   |   yes    |  yes  |
| [`blank-imports`](./RULES_DESCRIPTIONS.md#blank-imports)       |  n/a   | Disallows blank imports                                          |   yes    |  no   |
| [`context-as-argument`](./RULES_DESCRIPTIONS.md#context-as-argument) |  n/a   | `context.Context` should be the first argument of a function.    |   yes    |  no   |
| [`dot-imports`](./RULES_DESCRIPTIONS.md#dot-imports)         |  n/a   | Forbids `.` imports.                                             |   yes    |  no   |
| [`error-return`](./RULES_DESCRIPTIONS.md#error-return)        |  n/a   | The error return parameter should be last.                       |   yes    |  no   |
| [`error-strings`](./RULES_DESCRIPTIONS.md#error-strings)       |  []string   | Conventions around error strings.                                |   yes    |  no   |
| [`error-naming`](./RULES_DESCRIPTIONS.md#error-naming)        |  n/a   | Naming of error variables.                                       |   yes    |  no   |
| [`exported`](./RULES_DESCRIPTIONS.md#exported)            |  []string   | Naming and commenting conventions on exported symbols.           |   yes    |  no   |
| [`if-return`](./RULES_DESCRIPTIONS.md#if-return)           |  n/a   | Redundant if when returning an error.                            |   no    |  no   |
| [`increment-decrement`](./RULES_DESCRIPTIONS.md#increment-decrement) |  n/a   | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            |   yes    |  no   |
| [`var-naming`](./RULES_DESCRIPTIONS.md#var-naming)          |  allowlist & blocklist of initialisms   | Naming rules.                                                    |   yes    |  no   |
| [`package-comments`](./RULES_DESCRIPTIONS.md#package-comments)    |  n/a   | Package commenting conventions.                                  |   yes    |  no   |
| [`range`](./RULES_DESCRIPTIONS.md#range)               |  n/a   | Prevents redundant variables when iterating over a collection.   |   yes    |  no   |
| [`receiver-naming`](./RULES_DESCRIPTIONS.md#receiver-naming)     |  n/a   | Conventions around the naming of receivers.                      |   yes    |  no   |
| [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)   |  []string   | Prevents redundant else statements.                              |   yes    |  no   |
| [`argument-limit`](./RULES_DESCRIPTIONS.md#argument-limit)      |  int (defaults to 8)  | Specifies the maximum number of arguments a function can receive |    no    |  no   |
| [`cyclomatic`](./RULES_DESCRIPTIONS.md#cyclomatic)          |  int (defaults to 10)   | Sets restriction for maximum Cyclomatic complexity.              |    no    |  no   |
| [`max-public-structs`](./RULES_DESCRIPTIONS.md#max-public-structs)  |  int (defaults to 5)  | The maximum number of public structs in a file.                  |    no    |  no   |
| [`file-header`](./RULES_DESCRIPTIONS.md#file-header)         | string (defaults to none)| Header which each file should have.                              |    no    |  no   |
| [`empty-block`](./RULES_DESCRIPTIONS.md#empty-block)         |  n/a   | Warns on empty code blocks                                       |    no    |  yes   |
| [`superfluous-else`](./RULES_DESCRIPTIONS.md#superfluous-else)    |  []string   | Prevents redundant else statements (extends [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)) |    no    |  no   |
| [`confusing-naming`](./RULES_DESCRIPTIONS.md#confusing-naming)    |  n/a   | Warns on methods with names that differ only by capitalization   |    no    |  no   |
| [`get-return`](./RULES_DESCRIPTIONS.md#get-return)          |  n/a   | Warns on getters that do not yield any result                    |    no    |  no   |
| [`modifies-parameter`](./RULES_DESCRIPTIONS.md#modifies-parameter)  |  n/a   | Warns on assignments to function parameters                      |    no    |  no   |
| [`confusing-results`](./RULES_DESCRIPTIONS.md#confusing-results)   |  n/a   | Suggests to name potentially confusing function results          |    no    |  no   |
| [`deep-exit`](./RULES_DESCRIPTIONS.md#deep-exit)           |  n/a   | Looks for program exits in funcs other than `main()` or `init()` |    no    |  no   |
| [`unused-parameter`](./RULES_DESCRIPTIONS.md#unused-parameter)    |  n/a   | Suggests to rename or remove unused function parameters          |    no    |  no   |
| [`unreachable-code`](./RULES_DESCRIPTIONS.md#unreachable-code)    |  n/a   | Warns on unreachable code                                        |    no    |  no   |
| [`add-constant`](./RULES_DESCRIPTIONS.md#add-constant)        |  map   | Suggests using constant for magic numbers and string literals    |    no    |  no   |
| [`flag-parameter`](./RULES_DESCRIPTIONS.md#flag-parameter)      |  n/a   | Warns on boolean parameters that create a control coupling       |    no    |  no   |
| [`unnecessary-stmt`](./RULES_DESCRIPTIONS.md#unnecessary-stmt)    |  n/a   | Suggests removing or simplifying unnecessary statements          |    no    |  no   |
| [`struct-tag`](./RULES_DESCRIPTIONS.md#struct-tag)          |  []string   | Checks common struct tags like `json`, `xml`, `yaml`               |    no    |  no   |
| [`modifies-value-receiver`](./RULES_DESCRIPTIONS.md#modifies-value-receiver) |  n/a   | Warns on assignments to value-passed method receivers        |    no    |  yes  |
| [`constant-logical-expr`](./RULES_DESCRIPTIONS.md#constant-logical-expr)   |  n/a   | Warns on constant logical expressions                        |    no    |  no   |
| [`bool-literal-in-expr`](./RULES_DESCRIPTIONS.md#bool-literal-in-expr)|  n/a   | Suggests removing Boolean literals from logic expressions        |    no    |  no   |
| [`redefines-builtin-id`](./RULES_DESCRIPTIONS.md#redefines-builtin-id)|  n/a   | Warns on redefinitions of builtin identifiers                    |    no    |  no   |
| [`function-result-limit`](./RULES_DESCRIPTIONS.md#function-result-limit) |  int (defaults to 3)| Specifies the maximum number of results a function can return    |    no    |  no   |
| [`imports-blocklist`](./RULES_DESCRIPTIONS.md#imports-blocklist)   | []string | Disallows importing the specified packages                     |    no    |  no   |
| [`range-val-in-closure`](./RULES_DESCRIPTIONS.md#range-val-in-closure)|  n/a   | Warns if range value is used in a closure dispatched as goroutine|    no    |  no   |
| [`range-val-address`](./RULES_DESCRIPTIONS.md#range-val-address)|  n/a   | Warns if address of range value is used dangerously |    no    |  yes   |
| [`waitgroup-by-value`](./RULES_DESCRIPTIONS.md#waitgroup-by-value)  |  n/a   | Warns on functions taking sync.WaitGroup as a by-value parameter |    no    |  no   |
| [`atomic`](./RULES_DESCRIPTIONS.md#atomic)              |  n/a   | Check for common mistaken usages of the `sync/atomic` package    |    no    |  no   |
| [`empty-lines`](./RULES_DESCRIPTIONS.md#empty-lines)   | n/a | Warns when there are heading or trailing newlines in a block              |    no    |  no   |
| [`line-length-limit`](./RULES_DESCRIPTIONS.md#line-length-limit)   | int (defaults to 80) | Specifies the maximum number of characters in a line             |    no    |  no   |
| [`call-to-gc`](./RULES_DESCRIPTIONS.md#call-to-gc)   | n/a    | Warns on explicit call to the garbage collector    |    no    |  no   |
| [`duplicated-imports`](./RULES_DESCRIPTIONS.md#duplicated-imports) | n/a  | Looks for packages that are imported two or more times   |    no    |  no   |
| [`import-shadowing`](./RULES_DESCRIPTIONS.md#import-shadowing)   | n/a    | Spots identifiers that shadow an import    |    no    |  no   |
| [`bare-return`](./RULES_DESCRIPTIONS.md#bare-return) | n/a  | Warns on bare returns   |    no    |  no   |
| [`unused-receiver`](./RULES_DESCRIPTIONS.md#unused-receiver)   | n/a    | Suggests to rename or remove unused method receivers    |    no    |  no   |
| [`unhandled-error`](./RULES_DESCRIPTIONS.md#unhandled-error)   | []string   | Warns on unhandled errors returned by function calls    |    no    |  yes   |
| [`cognitive-complexity`](./RULES_DESCRIPTIONS.md#cognitive-complexity)          |  int (defaults to 7) | Sets restriction for maximum Cognitive complexity.              |    no    |  no   |
| [`string-of-int`](./RULES_DESCRIPTIONS.md#string-of-int)          |  n/a   | Warns on suspicious casts from int to string            |    no    |  yes   |
| [`string-format`](./RULES_DESCRIPTIONS.md#string-format)          |  map   | Warns on specific string literals that fail one or more user-configured regular expressions            |    no    |  no   |
| [`early-return`](./RULES_DESCRIPTIONS.md#early-return)          | []string   | Spots if-then-else statements where the predicate may be inverted to reduce nesting |    no    |  no   |
| [`unconditional-recursion`](./RULES_DESCRIPTIONS.md#unconditional-recursion)          |  n/a   | Warns on function calls that will lead to (direct) infinite recursion |    no    |  no   |
| [`identical-branches`](./RULES_DESCRIPTIONS.md#identical-branches)          |  n/a   | Spots if-then-else statements with identical `then` and `else` branches       |    no    |  no   |
| [`defer`](./RULES_DESCRIPTIONS.md#defer)          |  map   |  Warns on some [defer gotchas](https://blog.learngoprogramming.com/5-gotchas-of-defer-in-go-golang-part-iii-36a1ab3d6ef1)       |    no    |  no   |
| [`unexported-naming`](./RULES_DESCRIPTIONS.md#unexported-naming)          |  n/a   |  Warns on wrongly named un-exported symbols       |    no    |  no   |
| [`function-length`](./RULES_DESCRIPTIONS.md#function-length)          |  int, int (defaults to 50 statements, 75 lines)   |  Warns on functions exceeding the statements or lines max |    no    |  no   |
| [`nested-structs`](./RULES_DESCRIPTIONS.md#nested-structs)          |  n/a   |  Warns on structs within structs |    no    |  no   |
| [`useless-break`](./RULES_DESCRIPTIONS.md#useless-break)          |  n/a   |  Warns on useless `break` statements in case clauses |    no    |  no   |
| [`banned-characters`](./RULES_DESCRIPTIONS.md#banned-characters)          |  []string (defaults to []string{})   |  Checks banned characters in identifiers |    no    |  no   |
| [`optimize-operands-order`](./RULES_DESCRIPTIONS.md#optimize-operands-order)          |  n/a   |  Checks inefficient conditional expressions |    no    |  no   |
| [`use-any`](./RULES_DESCRIPTIONS.md#use-any)          |  n/a   |  Proposes to replace `interface{}` with its alias `any` |    no    |  no   |
| [`datarace`](./RULES_DESCRIPTIONS.md#datarace)          |  n/a   |  Spots potential dataraces |    no    |  no   |
| [`comment-spacings`](./RULES_DESCRIPTIONS.md#comment-spacings)          |  []string   |  Warns on malformed comments |    no    |  no   |
| [`redundant-import-alias`](./RULES_DESCRIPTIONS.md#redundant-import-alias)          |  n/a   |  Warns on import aliases matching the imported package name |    no    |  no   |
| [`import-alias-naming`](./RULES_DESCRIPTIONS.md#import-alias-naming)         | string or map[string]string (defaults to allow regex pattern ^[a-z][a-z0-9]{0,}$) | Conventions around the naming of import aliases.                              |    no    |  no   |
| [`enforce-map-style`](./RULES_DESCRIPTIONS.md#enforce-map-style) |  string (defaults to "any")  |  Enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization. Does not affect `make(map[type]type, size)` constructions. |    no    |  no   |
| [`enforce-slice-style`](./RULES_DESCRIPTIONS.md#enforce-slice-style) |  string (defaults to "any")  |  Enforces consistent usage of `make([]type, 0)` or `[]type{}` for slice initialization. Does not affect `make(map[type]type, non_zero_len, or_non_zero_cap)` constructions. |    no    |  no   |
| [`enforce-repeated-arg-type-style`](./RULES_DESCRIPTIONS.md#enforce-repeated-arg-type-style) |  string (defaults to "any")  |  Enforces consistent style for repeated argument and/or return value types. |    no    |  no   |
| [`max-control-nesting`](./RULES_DESCRIPTIONS.md#max-control-nesting) |  int (defaults to 5)  | Sets restriction for maximum nesting of control structures. |    no    |  no   |
| [`comments-density`](./RULES_DESCRIPTIONS.md#comments-density) |  int (defaults to 0)  | Enforces a minimum comment / code relation |    no    |  no   |


## Configurable rules

Here you can find how you can configure some existing rules:

### `var-naming`

This rule accepts two slices of strings, an allowlist and a blocklist of initialisms. By default, the rule behaves exactly as the alternative in `golint` but optionally, you can relax it (see [golint/lint/issues/89](https://github.com/golang/lint/issues/89))

```toml
[rule.var-naming]
  arguments = [["ID"], ["VM"]]
```

This way, revive will not warn for an identifier called `customId` but will warn that `customVm` should be called `customVM`.

## Available Formatters

This section lists all the available formatters and provides a screenshot for each one.

### Friendly

![Friendly formatter](/assets/formatter-friendly.png)

### Stylish

![Stylish formatter](/assets/formatter-stylish.png)

### Default

The default formatter produces the same output as `golint`.

![Default formatter](/assets/formatter-default.png)

### Plain

The plain formatter produces the same output as the default formatter and appends the URL to the rule description.

![Plain formatter](/assets/formatter-plain.png)

### Unix

The unix formatter produces the same output as the default formatter but surrounds the rules in `[]`.

![Unix formatter](/assets/formatter-unix.png)

### SARIF
The `sarif`  formatter produces outputs in SARIF, for _Static Analysis Results Interchange Format_, a standard JSON-based format for the output of static analysis tools defined and promoted by [OASIS](https://www.oasis-open.org/).

Current supported version of the standard is [SARIF-v2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/csprd01/sarif-v2.1.0-csprd01.html
).

## Extensibility

The tool can be extended with custom rules or formatters. This section contains additional information on how to implement such.

To extend the linter with a custom rule you can push it to this repository or use `revive` as a library (see below)

To add a custom formatter you'll have to push it to this repository or fork it. This is due to the limited `-buildmode=plugin` support which [works only on Linux (with known issues)](https://golang.org/pkg/plugin/).

### Writing a Custom Rule

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```

The `Arguments` type is an alias of the type `[]interface{}`. The arguments of the rule are passed from the configuration file.

#### Example

Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with a given identifier. We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
  arguments = ["Foo"]
```

With the snippet above we:

- Enable the rule with the name `ban-struct-name`. The `Name()` method of our rule should return a string that matches `ban-struct-name`.
- Configure the rule with the argument `Foo`. The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

A sample rule implementation can be found [here](/rule/argument-limit.go).

#### Using `revive` as a library
If a rule is specific to your use case
(i.e. it is not a good candidate to be added to `revive`'s rule set) you can add it to your linter using `revive` as a linting engine.

The following code shows how to use `revive` in your application.
In the example only one rule is added (`myRule`), of course, you can add as many as you need to.
Your rules can be configured programmatically or with the standard `revive` configuration file.
The full rule set of `revive` is also actionable by your application.

```go
package main

import (
	"github.com/mgechev/revive/cli"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

func main() {
	cli.RunRevive(revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}))
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(*lint.File, lint.Arguments) []lint.Failure { ... }
```

You can still go further and use `revive` without its CLI, as part of your library, or your CLI:

```go
package mylib

import (
	"github.com/mgechev/revive/cli"
	"github.com/mgechev/revive/revivelib"
	"github.com/mgechev/revive/lint"
)

// Error checking removed for clarity
func LintMyFile(file string) {
	conf, _:= config.GetConfig("../defaults.toml")

	revive, _ := revivelib.New(
		conf,  // Configuration file
		true,  // Set exit status
		2048,  // Max open files

		// Then add as many extra rules as you need
		revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}),
	)

	failuresChan, err := revive.Lint(
 		revivelib.Include(file),
 		revivelib.Exclude("./fixtures"),
 		// You can use as many revivelib.Include or revivelib.Exclude as required
 	)
  	if err != nil {
  	 	panic("Shouldn't have failed: " + err.Error)
  	}

  	// Now let's return the formatted errors
	failures, exitCode, _ := revive.Format("stylish", failuresChan)

  	// failures is the string with all formatted lint error messages
  	// exit code is 0 if no errors, 1 if errors (unless config options change it)
  	// ... do something with them
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(*lint.File, lint.Arguments) []lint.Failure { ... }
```

### Custom Formatter

Each formatter needs to implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, Config) (string, error)
	Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration of the enabled rules. The `Name()` method should return a string different from the names of the already existing rules. This string is used when specifying the formatter when invoking the `revive` CLI tool.

For a sample formatter, take a look at [this file](/formatter/json.go).

## Speed Comparison

Compared to `golint`, `revive` performs better because it lints the files for each individual rule into a separate goroutine. Here's a basic performance benchmark on MacBook Pro Early 2013 run on Kubernetes:

### golint

```shell
time golint kubernetes/... > /dev/null

real    0m54.837s
user    0m57.844s
sys     0m9.146s
```

### revive

```shell
# no type checking
time revive -config untyped.toml kubernetes/... > /dev/null

real    0m8.471s
user    0m40.721s
sys     0m3.262s
```

Keep in mind that if you use rules that require type checking, the performance may drop to 2x faster than `golint`:

```shell
# type checking enabled
time revive kubernetes/... > /dev/null

real    0m26.211s
user    2m6.708s
sys     0m17.192s
```

Currently, type-checking is enabled by default. If you want to run the linter without type-checking, remove all typed rules from the configuration file.

## Overriding colorization detection

By default, `revive` determines whether or not to colorize its output based on whether it's connected to a TTY or not.
This works for most use cases, but may not behave as expected if you use `revive` in a pipeline of commands,
where STDOUT is being piped to another command.

To force colorization, add `REVIVE_FORCE_COLOR=1` to the environment you're running in. For example:

```shell
REVIVE_FORCE_COLOR=1 revive -formatter friendly ./... | tee revive.log
```

## Contributors

[<img alt="renovate[bot]" src="https://avatars.githubusercontent.com/in/2740?v=4&s=117" width="117">](https://github.com/renovate[bot]) |[<img alt="mgechev" src="https://avatars.githubusercontent.com/u/455023?v=4&s=117" width="117">](https://github.com/mgechev) |[<img alt="chavacava" src="https://avatars.githubusercontent.com/u/25788468?v=4&s=117" width="117">](https://github.com/chavacava) |[<img alt="renovate-bot" src="https://avatars.githubusercontent.com/u/25180681?v=4&s=117" width="117">](https://github.com/renovate-bot) |[<img alt="xuri" src="https://avatars.githubusercontent.com/u/2809468?v=4&s=117" width="117">](https://github.com/xuri) |[<img alt="denisvmedia" src="https://avatars.githubusercontent.com/u/5462781?v=4&s=117" width="117">](https://github.com/denisvmedia) |
:---: |:---: |:---: |:---: |:---: |:---: |
[renovate[bot]](https://github.com/renovate[bot]) |[mgechev](https://github.com/mgechev) |[chavacava](https://github.com/chavacava) |[renovate-bot](https://github.com/renovate-bot) |[xuri](https://github.com/xuri) |[denisvmedia](https://github.com/denisvmedia) |

[<img alt="mfederowicz" src="https://avatars.githubusercontent.com/u/57678185?v=4&s=117" width="117">](https://github.com/mfederowicz) |[<img alt="doniacld" src="https://avatars.githubusercontent.com/u/19799268?v=4&s=117" width="117">](https://github.com/doniacld) |[<img alt="Clivern" src="https://avatars.githubusercontent.com/u/1634427?v=4&s=117" width="117">](https://github.com/Clivern) |[<img alt="morphy2k" src="https://avatars.githubusercontent.com/u/4280578?v=4&s=117" width="117">](https://github.com/morphy2k) |[<img alt="bernhardreisenberger" src="https://avatars.githubusercontent.com/u/5809300?v=4&s=117" width="117">](https://github.com/bernhardreisenberger) |[<img alt="alexandear" src="https://avatars.githubusercontent.com/u/3228886?v=4&s=117" width="117">](https://github.com/alexandear) |
:---: |:---: |:---: |:---: |:---: |:---: |
[mfederowicz](https://github.com/mfederowicz) |[doniacld](https://github.com/doniacld) |[Clivern](https://github.com/Clivern) |[morphy2k](https://github.com/morphy2k) |[bernhardreisenberger](https://github.com/bernhardreisenberger) |[alexandear](https://github.com/alexandear) |

[<img alt="dshemin" src="https://avatars.githubusercontent.com/u/11780307?v=4&s=117" width="117">](https://github.com/dshemin) |[<img alt="butuzov" src="https://avatars.githubusercontent.com/u/651824?v=4&s=117" width="117">](https://github.com/butuzov) |[<img alt="rawen17" src="https://avatars.githubusercontent.com/u/36483900?v=4&s=117" width="117">](https://github.com/rawen17) |[<img alt="sina-devel" src="https://avatars.githubusercontent.com/u/61763643?v=4&s=117" width="117">](https://github.com/sina-devel) |[<img alt="tymonx" src="https://avatars.githubusercontent.com/u/8367378?v=4&s=117" width="117">](https://github.com/tymonx) |[<img alt="mdelah" src="https://avatars.githubusercontent.com/u/4904544?v=4&s=117" width="117">](https://github.com/mdelah) |
:---: |:---: |:---: |:---: |:---: |:---: |
[dshemin](https://github.com/dshemin) |[butuzov](https://github.com/butuzov) |[rawen17](https://github.com/rawen17) |[sina-devel](https://github.com/sina-devel) |[tymonx](https://github.com/tymonx) |[mdelah](https://github.com/mdelah) |

[<img alt="ldez" src="https://avatars.githubusercontent.com/u/5674651?v=4&s=117" width="117">](https://github.com/ldez) |[<img alt="gsamokovarov" src="https://avatars.githubusercontent.com/u/604618?v=4&s=117" width="117">](https://github.com/gsamokovarov) |[<img alt="comdiv" src="https://avatars.githubusercontent.com/u/2387862?v=4&s=117" width="117">](https://github.com/comdiv) |[<img alt="heynemann" src="https://avatars.githubusercontent.com/u/60965?v=4&s=117" width="117">](https://github.com/heynemann) |[<img alt="cce" src="https://avatars.githubusercontent.com/u/51567?v=4&s=117" width="117">](https://github.com/cce) |[<img alt="mapreal19" src="https://avatars.githubusercontent.com/u/3055997?v=4&s=117" width="117">](https://github.com/mapreal19) |
:---: |:---: |:---: |:---: |:---: |:---: |
[ldez](https://github.com/ldez) |[gsamokovarov](https://github.com/gsamokovarov) |[comdiv](https://github.com/comdiv) |[heynemann](https://github.com/heynemann) |[cce](https://github.com/cce) |[mapreal19](https://github.com/mapreal19) |

[<img alt="zimmski" src="https://avatars.githubusercontent.com/u/1847950?v=4&s=117" width="117">](https://github.com/zimmski) |[<img alt="shmsr" src="https://avatars.githubusercontent.com/u/51480165?v=4&s=117" width="117">](https://github.com/shmsr) |[<img alt="git-hulk" src="https://avatars.githubusercontent.com/u/4987594?v=4&s=117" width="117">](https://github.com/git-hulk) |[<img alt="ccoVeille" src="https://avatars.githubusercontent.com/u/3875889?v=4&s=117" width="117">](https://github.com/ccoVeille) |[<img alt="tamird" src="https://avatars.githubusercontent.com/u/1535036?v=4&s=117" width="117">](https://github.com/tamird) |[<img alt="markelog" src="https://avatars.githubusercontent.com/u/945528?v=4&s=117" width="117">](https://github.com/markelog) |
:---: |:---: |:---: |:---: |:---: |:---: |
[zimmski](https://github.com/zimmski) |[shmsr](https://github.com/shmsr) |[git-hulk](https://github.com/git-hulk) |[ccoVeille](https://github.com/ccoVeille) |[tamird](https://github.com/tamird) |[markelog](https://github.com/markelog) |

[<img alt="mihaitodor" src="https://avatars.githubusercontent.com/u/788216?v=4&s=117" width="117">](https://github.com/mihaitodor) |[<img alt="dvejmz" src="https://avatars.githubusercontent.com/u/9487006?v=4&s=117" width="117">](https://github.com/dvejmz) |[<img alt="damif94" src="https://avatars.githubusercontent.com/u/29461526?v=4&s=117" width="117">](https://github.com/damif94) |[<img alt="abeltay" src="https://avatars.githubusercontent.com/u/15604207?v=4&s=117" width="117">](https://github.com/abeltay) |[<img alt="zeripath" src="https://avatars.githubusercontent.com/u/1824502?v=4&s=117" width="117">](https://github.com/zeripath) |[<img alt="Groxx" src="https://avatars.githubusercontent.com/u/77197?v=4&s=117" width="117">](https://github.com/Groxx) |
:---: |:---: |:---: |:---: |:---: |:---: |
[mihaitodor](https://github.com/mihaitodor) |[dvejmz](https://github.com/dvejmz) |[damif94](https://github.com/damif94) |[abeltay](https://github.com/abeltay) |[zeripath](https://github.com/zeripath) |[Groxx](https://github.com/Groxx) |

[<img alt="StephenBrown2" src="https://avatars.githubusercontent.com/u/1148665?v=4&s=117" width="117">](https://github.com/StephenBrown2) |[<img alt="qascade" src="https://avatars.githubusercontent.com/u/54154054?v=4&s=117" width="117">](https://github.com/qascade) |[<img alt="ridvansumset" src="https://avatars.githubusercontent.com/u/26631560?v=4&s=117" width="117">](https://github.com/ridvansumset) |[<img alt="rliebz" src="https://avatars.githubusercontent.com/u/5321575?v=4&s=117" width="117">](https://github.com/rliebz) |[<img alt="rdeusser" src="https://avatars.githubusercontent.com/u/5935071?v=4&s=117" width="117">](https://github.com/rdeusser) |[<img alt="rmarku" src="https://avatars.githubusercontent.com/u/1113370?v=4&s=117" width="117">](https://github.com/rmarku) |
:---: |:---: |:---: |:---: |:---: |:---: |
[StephenBrown2](https://github.com/StephenBrown2) |[qascade](https://github.com/qascade) |[ridvansumset](https://github.com/ridvansumset) |[rliebz](https://github.com/rliebz) |[rdeusser](https://github.com/rdeusser) |[rmarku](https://github.com/rmarku) |

[<img alt="rnikoopour" src="https://avatars.githubusercontent.com/u/7692789?v=4&s=117" width="117">](https://github.com/rnikoopour) |[<img alt="rafamadriz" src="https://avatars.githubusercontent.com/u/67771985?v=4&s=117" width="117">](https://github.com/rafamadriz) |[<img alt="paco0x" src="https://avatars.githubusercontent.com/u/6123425?v=4&s=117" width="117">](https://github.com/paco0x) |[<img alt="weastur" src="https://avatars.githubusercontent.com/u/10865586?v=4&s=117" width="117">](https://github.com/weastur) |[<img alt="pa-m" src="https://avatars.githubusercontent.com/u/5503106?v=4&s=117" width="117">](https://github.com/pa-m) |[<img alt="cinar" src="https://avatars.githubusercontent.com/u/1754092?v=4&s=117" width="117">](https://github.com/cinar) |
:---: |:---: |:---: |:---: |:---: |:---: |
[rnikoopour](https://github.com/rnikoopour) |[rafamadriz](https://github.com/rafamadriz) |[paco0x](https://github.com/paco0x) |[weastur](https://github.com/weastur) |[pa-m](https://github.com/pa-m) |[cinar](https://github.com/cinar) |

[<img alt="natefinch" src="https://avatars.githubusercontent.com/u/3185864?v=4&s=117" width="117">](https://github.com/natefinch) |[<img alt="flesser" src="https://avatars.githubusercontent.com/u/510681?v=4&s=117" width="117">](https://github.com/flesser) |[<img alt="y-yagi" src="https://avatars.githubusercontent.com/u/987638?v=4&s=117" width="117">](https://github.com/y-yagi) |[<img alt="techknowlogick" src="https://avatars.githubusercontent.com/u/164197?v=4&s=117" width="117">](https://github.com/techknowlogick) |[<img alt="okhowang" src="https://avatars.githubusercontent.com/u/3352585?v=4&s=117" width="117">](https://github.com/okhowang) |[<img alt="meanguy" src="https://avatars.githubusercontent.com/u/78570571?v=4&s=117" width="117">](https://github.com/meanguy) |
:---: |:---: |:---: |:---: |:---: |:---: |
[natefinch](https://github.com/natefinch) |[flesser](https://github.com/flesser) |[y-yagi](https://github.com/y-yagi) |[techknowlogick](https://github.com/techknowlogick) |[okhowang](https://github.com/okhowang) |[meanguy](https://github.com/meanguy) |

[<img alt="likyh" src="https://avatars.githubusercontent.com/u/3294100?v=4&s=117" width="117">](https://github.com/likyh) |[<img alt="kerneltravel" src="https://avatars.githubusercontent.com/u/437879?v=4&s=117" width="117">](https://github.com/kerneltravel) |[<img alt="jmckenzieark" src="https://avatars.githubusercontent.com/u/70923399?v=4&s=117" width="117">](https://github.com/jmckenzieark) |[<img alt="haya14busa" src="https://avatars.githubusercontent.com/u/3797062?v=4&s=117" width="117">](https://github.com/haya14busa) |[<img alt="fregin" src="https://avatars.githubusercontent.com/u/23256240?v=4&s=117" width="117">](https://github.com/fregin) |[<img alt="ydah" src="https://avatars.githubusercontent.com/u/13041216?v=4&s=117" width="117">](https://github.com/ydah) |
:---: |:---: |:---: |:---: |:---: |:---: |
[likyh](https://github.com/likyh) |[kerneltravel](https://github.com/kerneltravel) |[jmckenzieark](https://github.com/jmckenzieark) |[haya14busa](https://github.com/haya14busa) |[fregin](https://github.com/fregin) |[ydah](https://github.com/ydah) |

[<img alt="WillAbides" src="https://avatars.githubusercontent.com/u/233500?v=4&s=117" width="117">](https://github.com/WillAbides) |[<img alt="heyvito" src="https://avatars.githubusercontent.com/u/77198?v=4&s=117" width="117">](https://github.com/heyvito) |[<img alt="scop" src="https://avatars.githubusercontent.com/u/109152?v=4&s=117" width="117">](https://github.com/scop) |[<img alt="vkrol" src="https://avatars.githubusercontent.com/u/153412?v=4&s=117" width="117">](https://github.com/vkrol) |[<img alt="Jarema" src="https://avatars.githubusercontent.com/u/7369771?v=4&s=117" width="117">](https://github.com/Jarema) |[<img alt="tartale" src="https://avatars.githubusercontent.com/u/9323250?v=4&s=117" width="117">](https://github.com/tartale) |
:---: |:---: |:---: |:---: |:---: |:---: |
[WillAbides](https://github.com/WillAbides) |[heyvito](https://github.com/heyvito) |[scop](https://github.com/scop) |[vkrol](https://github.com/vkrol) |[Jarema](https://github.com/Jarema) |[tartale](https://github.com/tartale) |

[<img alt="tmzane" src="https://avatars.githubusercontent.com/u/73077675?v=4&s=117" width="117">](https://github.com/tmzane) |[<img alt="felipedavid" src="https://avatars.githubusercontent.com/u/75049173?v=4&s=117" width="117">](https://github.com/felipedavid) |[<img alt="euank" src="https://avatars.githubusercontent.com/u/2147649?v=4&s=117" width="117">](https://github.com/euank) |[<img alt="Juneezee" src="https://avatars.githubusercontent.com/u/20135478?v=4&s=117" width="117">](https://github.com/Juneezee) |[<img alt="echoix" src="https://avatars.githubusercontent.com/u/27212526?v=4&s=117" width="117">](https://github.com/echoix) |[<img alt="EXHades" src="https://avatars.githubusercontent.com/u/22260242?v=4&s=117" width="117">](https://github.com/EXHades) |
:---: |:---: |:---: |:---: |:---: |:---: |
[tmzane](https://github.com/tmzane) |[felipedavid](https://github.com/felipedavid) |[euank](https://github.com/euank) |[Juneezee](https://github.com/Juneezee) |[echoix](https://github.com/echoix) |[EXHades](https://github.com/EXHades) |

[<img alt="petethepig" src="https://avatars.githubusercontent.com/u/662636?v=4&s=117" width="117">](https://github.com/petethepig) |[<img alt="Dirk007" src="https://avatars.githubusercontent.com/u/17194484?v=4&s=117" width="117">](https://github.com/Dirk007) |[<img alt="yangdiangzb" src="https://avatars.githubusercontent.com/u/16643665?v=4&s=117" width="117">](https://github.com/yangdiangzb) |[<img alt="derekperkins" src="https://avatars.githubusercontent.com/u/3588778?v=4&s=117" width="117">](https://github.com/derekperkins) |[<img alt="bboreham" src="https://avatars.githubusercontent.com/u/8125524?v=4&s=117" width="117">](https://github.com/bboreham) |[<img alt="attiss" src="https://avatars.githubusercontent.com/u/23562566?v=4&s=117" width="117">](https://github.com/attiss) |
:---: |:---: |:---: |:---: |:---: |:---: |
[petethepig](https://github.com/petethepig) |[Dirk007](https://github.com/Dirk007) |[yangdiangzb](https://github.com/yangdiangzb) |[derekperkins](https://github.com/derekperkins) |[bboreham](https://github.com/bboreham) |[attiss](https://github.com/attiss) |

[<img alt="Aragur" src="https://avatars.githubusercontent.com/u/11004008?v=4&s=117" width="117">](https://github.com/Aragur) |[<img alt="kulti" src="https://avatars.githubusercontent.com/u/1286683?v=4&s=117" width="117">](https://github.com/kulti) |[<img alt="Abirdcfly" src="https://avatars.githubusercontent.com/u/5100555?v=4&s=117" width="117">](https://github.com/Abirdcfly) |[<img alt="abhinav" src="https://avatars.githubusercontent.com/u/41730?v=4&s=117" width="117">](https://github.com/abhinav) |[<img alt="r-ricci" src="https://avatars.githubusercontent.com/u/52817765?v=4&s=117" width="117">](https://github.com/r-ricci) |[<img alt="nunnatsa" src="https://avatars.githubusercontent.com/u/60659093?v=4&s=117" width="117">](https://github.com/nunnatsa) |
:---: |:---: |:---: |:---: |:---: |:---: |
[Aragur](https://github.com/Aragur) |[kulti](https://github.com/kulti) |[Abirdcfly](https://github.com/Abirdcfly) |[abhinav](https://github.com/abhinav) |[r-ricci](https://github.com/r-ricci) |[nunnatsa](https://github.com/nunnatsa) |

[<img alt="michalhisim" src="https://avatars.githubusercontent.com/u/764249?v=4&s=117" width="117">](https://github.com/michalhisim) |[<img alt="mathieu-aubin" src="https://avatars.githubusercontent.com/u/15820228?v=4&s=117" width="117">](https://github.com/mathieu-aubin) |[<img alt="martinsirbe" src="https://avatars.githubusercontent.com/u/13367583?v=4&s=117" width="117">](https://github.com/martinsirbe) |[<img alt="avorima" src="https://avatars.githubusercontent.com/u/15158349?v=4&s=117" width="117">](https://github.com/avorima) |[<img alt="very-amused" src="https://avatars.githubusercontent.com/u/44382255?v=4&s=117" width="117">](https://github.com/very-amused) |[<img alt="johnrichardrinehart" src="https://avatars.githubusercontent.com/u/6321578?v=4&s=117" width="117">](https://github.com/johnrichardrinehart) |
:---: |:---: |:---: |:---: |:---: |:---: |
[michalhisim](https://github.com/michalhisim) |[mathieu-aubin](https://github.com/mathieu-aubin) |[martinsirbe](https://github.com/martinsirbe) |[avorima](https://github.com/avorima) |[very-amused](https://github.com/very-amused) |[johnrichardrinehart](https://github.com/johnrichardrinehart) |

[<img alt="walles" src="https://avatars.githubusercontent.com/u/158201?v=4&s=117" width="117">](https://github.com/walles) |[<img alt="jefersonf" src="https://avatars.githubusercontent.com/u/3049540?v=4&s=117" width="117">](https://github.com/jefersonf) |[<img alt="jan-xyz" src="https://avatars.githubusercontent.com/u/5249233?v=4&s=117" width="117">](https://github.com/jan-xyz) |[<img alt="jamesmaidment" src="https://avatars.githubusercontent.com/u/2050324?v=4&s=117" width="117">](https://github.com/jamesmaidment) |[<img alt="grongor" src="https://avatars.githubusercontent.com/u/972493?v=4&s=117" width="117">](https://github.com/grongor) |[<img alt="tie" src="https://avatars.githubusercontent.com/u/14792994?v=4&s=117" width="117">](https://github.com/tie) |
:---: |:---: |:---: |:---: |:---: |:---: |
[walles](https://github.com/walles) |[jefersonf](https://github.com/jefersonf) |[jan-xyz](https://github.com/jan-xyz) |[jamesmaidment](https://github.com/jamesmaidment) |[grongor](https://github.com/grongor) |[tie](https://github.com/tie) |

[<img alt="quasilyte" src="https://avatars.githubusercontent.com/u/6286655?v=4&s=117" width="117">](https://github.com/quasilyte) |[<img alt="davidhsingyuchen" src="https://avatars.githubusercontent.com/u/17587061?v=4&s=117" width="117">](https://github.com/davidhsingyuchen) |[<img alt="gburanov" src="https://avatars.githubusercontent.com/u/2969603?v=4&s=117" width="117">](https://github.com/gburanov) |[<img alt="ginglis13" src="https://avatars.githubusercontent.com/u/43075615?v=4&s=117" width="117">](https://github.com/ginglis13) |
:---: |:---: |:---: |:---: |
[quasilyte](https://github.com/quasilyte) |[davidhsingyuchen](https://github.com/davidhsingyuchen) |[gburanov](https://github.com/gburanov) |[ginglis13](https://github.com/ginglis13) |

## License

MIT
