[![Build Status](https://travis-ci.org/mgechev/revive.svg?branch=master)](https://travis-ci.org/mgechev/revive)

# revive

Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint. **`Revive` provides a framework for development of custom rules, and lets you define a strict preset for enhancing your development & code review processes**.

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
  <br>
  Logo by <a href="https://github.com/hawkgs">Georgi Serev</a>
</p>

Here's how `revive` is different from `golint`:

- Allows us to enable or disable rules using a configuration file.
- Allows us to configure the linting rules with a TOML file.
- 2x faster running the same rules as golint.
- Provides functionality for disabling a specific rule or the entire linter for a file or a range of lines.
  - `golint` allows this only for generated files.
- Optional type checking. Most rules in golint do not require type checking. If you disable them in the config file, revive will run over 6x faster than golint.
- Provides multiple formatters which let us customize the output.
- Allows us to customize the return code for the entire linter or based on the failure of only some rules.
- _Everyone can extend it easily with custom rules or formatters._
- `Revive` provides more rules compared to `golint`.

<p align="center">
  <img src="./assets/demo.svg" alt="" width="700">
</p>

<!-- TOC -->

- [revive](#revive)
  - [Usage](#usage)
    - [Text Editors](#text-editors)
    - [Installation](#installation)
    - [Command Line Flags](#command-line-flags)
    - [Sample Invocations](#sample-invocations)
    - [Comment Annotations](#comment-annotations)
    - [Configuration](#configuration)
    - [Default Configuration](#default-configuration)
    - [Recommended Configuration](#recommended-configuration)
  - [Available Rules](#available-rules)
  - [Available Formatters](#available-formatters)
    - [Friendly](#friendly)
    - [Stylish](#stylish)
    - [Default](#default)
  - [Extensibility](#extensibility)
    - [Custom Rule](#custom-rule)
      - [Example](#example)
    - [Custom Formatter](#custom-formatter)
  - [Speed Comparison](#speed-comparison)
    - [golint](#golint)
    - [revive](#revive-1)
  - [Contributors](#contributors)
  - [License](#license)

<!-- /TOC -->

## Usage

Since the default behavior of `revive` is compatible with `golint`, without providing any additional flags, the only difference you'd notice is faster execution.

### Text Editors

- Support for VSCode in [vscode-go](https://github.com/Microsoft/vscode-go/pull/1699).
- Support for vim via [w0rp/ale](https://github.com/w0rp/ale):

```vim
call ale#linter#Define('go', {
\   'name': 'revive',
\   'output_stream': 'both',
\   'executable': 'revive',
\   'read_buffer': 0,
\   'command': 'revive %t',
\   'callback': 'ale#handlers#unix#HandleAsWarning',
\})
```

- Support for Atom via [linter-revive](https://github.com/morphy2k/linter-revive).

### Installation

```bash
go get -u github.com/mgechev/revive
```

### Command Line Flags

`revive` accepts three command line parameters:

- `-config [PATH]` - path to config file in TOML format.
- `-exclude [PATTERN]` - pattern for files/directories/packages to be excluded for linting. You can specify the files you want to exclude for linting either as package name (i.e. `github.com/mgechev/revive`), list them as individual files (i.e. `file.go`), directories (i.e. `./foo/...`), or any combination of the three.
- `-formatter [NAME]` - formatter to be used for the output. The currently available formatters are:
  - `default` - will output the failures the same way that `golint` does.
  - `json` - outputs the failures in JSON format.
  - `ndjson` - outputs the failures as stream in newline delimited JSON (NDJSON) format.
  - `friendly` - outputs the failures when found. Shows summary of all the failures.
  - `stylish` - formats the failures in a table. Keep in mind that it doesn't stream the output so it might be perceived as slower compared to others.

### Sample Invocations

```shell
revive -config revive.toml -exclude file1.go -exclude file2.go -formatter friendly github.com/mgechev/revive package/...
```

- The command above will use the configuration from `revive.toml`
- `revive` will ignore `file1.go` and `file2.go`
- The output will be formatted with the `friendly` formatter
- The linter will analyze `github.com/mgechev/revive` and the files in `package`

### Comment Annotations

Using comments, you can disable the linter for the entire file or only range of lines:

```go
//revive:disable

func Public() {}
//revive:enable
```

The snippet above, will disable `revive` between the `revive:disable` and `revive:enable` comments. If you skip `revive:enable`, the linter will be disabled for the rest of the file.

You can do the same on a rule level. In case you want to disable only a particular rule, you can use:

```go
//revive:disable:unexported-return
func Public() private {
  return private
}
//revive:enable:unexported-return
```

This way, `revive` will not warn you for that you're returning an object of an unexported type, from an exported function.

### Configuration

`revive` can be configured with a TOML file. Here's a sample configuration with explanation for the individual properties:

```toml
# When set to false, ignores files with "GENERATED" header, similar to golint
ignoreGeneratedHeader = true

# Sets the default severity to "warning"
severity = "warning"

# Sets the default failure confidence. This means that linting errors
# with less than 0.8 confidence will be ignored.
confidence = 0.8

# Sets the error code for failures with severity "error"
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

### Default Configuration

The default configuration of `revive` can be found at `defaults.toml`. This will enable all rules available in `golint` and use their default configuration (i.e. the way they are hardcoded in `golint`).

```shell
revive -config defaults.toml github.com/mgechev/revive
```

This will use the configuration file `defaults.toml`, the `default` formatter, and will run linting over the `github.com/mgechev/revive` package.

### Recommended Configuration

```shell
revive -config config.toml -formatter friendly github.com/mgechev/revive
```

This will use `config.toml`, the `friendly` formatter, and will run linting over the `github.com/mgechev/revive` package.

## Available Rules

List of all available rules. The rules ported from `golint` are left unchanged and indicated in the `golint` column.

| Name                  | Config | Description                                                      | `golint` | Typed |
| --------------------- | :----: | :--------------------------------------------------------------- | :------: | :---: |
| `context-keys-type`   |  n/a   | Disallows the usage of basic types in `context.WithValue`.       |   yes    |  yes  |
| `time-naming`         |  n/a   | Conventions around the naming of time variables.                 |   yes    |  yes  |
| `var-declaration`     |  n/a   | Reduces redundancies around variable declaration.                |   yes    |  yes  |
| `unexported-return`   |  n/a   | Warns when a public return is from unexported type.              |   yes    |  yes  |
| `errorf`              |  n/a   | Should replace `error.New(fmt.Sprintf())` with `error.Errorf()`  |   yes    |  yes  |
| `blank-imports`       |  n/a   | Disallows blank imports                                          |   yes    |  no   |
| `context-as-argument` |  n/a   | `context.Context` should be the first argument of a function.    |   yes    |  no   |
| `dot-imports`         |  n/a   | Forbids `.` imports.                                             |   yes    |  no   |
| `error-return`        |  n/a   | The error return parameter should be last.                       |   yes    |  no   |
| `error-strings`       |  n/a   | Conventions around error strings.                                |   yes    |  no   |
| `error-naming`        |  n/a   | Naming of error variables.                                       |   yes    |  no   |
| `exported`            |  n/a   | Naming and commenting conventions on exported symbols.           |   yes    |  no   |
| `if-return`           |  n/a   | Redundant if when returning an error.                            |   yes    |  no   |
| `increment-decrement` |  n/a   | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            |   yes    |  no   |
| `var-naming`          |  n/a   | Naming rules.                                                    |   yes    |  no   |
| `package-comments`    |  n/a   | Package commenting conventions.                                  |   yes    |  no   |
| `range`               |  n/a   | Prevents redundant variables when iterating over a collection.   |   yes    |  no   |
| `receiver-naming`     |  n/a   | Conventions around the naming of receivers.                      |   yes    |  no   |
| `indent-error-flow`   |  n/a   | Prevents redundant else statements.                              |   yes    |  no   |
| `argument-limit`      |  int   | Specifies the maximum number of arguments a function can receive |    no    |  no   |
| `cyclomatic`          |  int   | Sets restriction for maximum Cyclomatic complexity.              |    no    |  no   |
| `max-public-structs`  |  int   | The maximum number of public structs in a file.                  |    no    |  no   |
| `file-header`         | string | Header which each file should have.                              |    no    |  no   |
| `empty-block`         |  n/a   | Warns on empty code blocks                                       |    no    |  no   |
| `superfluous-else`    |  n/a   | Prevents redundant else statements (extends `indent-error-flow`) |    no    |  no   |
| `confusing-naming`    |  n/a   | Warns on methods with names that differ only by capitalization   |    no    |  no   |
| `get-return`          |  n/a   | Warns on getters that do not yield any result                    |    no    |  no   |
| `modifies-parameter`  |  n/a   | Warns on assignments to function parameters                      |    no    |  no   |
| `confusing-results`   |  n/a   | Suggests to name potentially confusing function results          |    no    |  no   |
| `deep-exit`           |  n/a   | Looks for program exits in funcs other than `main()` or `init()` |    no    |  no   |
| `unused-parameter`    |  n/a   | Suggests to rename or remove unused function parameters          |    no    |  no   |

## Available Formatters

This section lists all the available formatters and provides a screenshot for each one.

### Friendly

![Friendly formatter](/assets/friendly-formatter.png)

### Stylish

![Stylish formatter](/assets/stylish-formatter.png)

### Default

![Default formatter](/assets/default-formatter.png)

## Extensibility

The tool can be extended with custom rules or formatters. This section contains additional information on how to implement such.

**To extend the linter with a custom rule or a formatter you'll have to push it to this repository or fork it**. This is due to the limited `-buildmode=plugin` support which [works only on Linux (with known issues)](https://golang.org/pkg/plugin/).

### Custom Rule

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```

The `Arguments` type is an alias of the type `[]interface{}`. The arguments of the rule are passed from the configuration file.

#### Example

Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with given identifier. We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
  arguments = ["Foo"]
```

With the snippet above we:

- Enable the rule with name `ban-struct-name`. The `Name()` method of our rule should return a string which matches `ban-struct-name`.
- Configure the rule with the argument `Foo`. The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

A sample rule implementation can be found [here](/rule/argument-limit.go).

### Custom Formatter

Each formatter needs to implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration of the enabled rules. The `Name()` method should return a string different from the names of the already existing rules. This string is used when specifying the formatter when invoking the `revive` CLI tool.

For a sample formatter, take a look at [this file](/formatter/json.go).

## Speed Comparison

Compared to `golint`, `revive` performs better because it lints the files for each individual rule into a separate goroutine. Here's a basic performance benchmark on MacBook Pro Early 2013 run on kubernetes:

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

Keep in mind that if you use rules which require type checking, the performance may drop to 2x faster than `golint`:

```shell
# type checking enabled
time revive kubernetes/... > /dev/null

real    0m26.211s
user    2m6.708s
sys     0m17.192s
```

Currently, type checking is enabled by default. If you want to run the linter without type checking, remove all typed rules from the configuration file.

## Contributors

| [<img alt="mgechev" src="https://avatars1.githubusercontent.com/u/455023?v=4&s=117" width="117">](https://github.com/mgechev) | [<img alt="chavacava" src="https://avatars2.githubusercontent.com/u/25788468?v=4&s=117" width="117">](https://github.com/chavacava) | [<img alt="morphy2k" src="https://avatars2.githubusercontent.com/u/4280578?v=4&s=117" width="117">](https://github.com/morphy2k) | [<img alt="tamird" src="https://avatars0.githubusercontent.com/u/1535036?v=4&s=117" width="117">](https://github.com/tamird) | [<img alt="paul-at-start" src="https://avatars2.githubusercontent.com/u/5486775?v=4&s=117" width="117">](https://github.com/paul-at-start) | [<img alt="vkrol" src="https://avatars3.githubusercontent.com/u/153412?v=4&s=117" width="117">](https://github.com/vkrol) |
| :---------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------: | :------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------: |
|                                             [mgechev](https://github.com/mgechev)                                             |                                              [chavacava](https://github.com/chavacava)                                              |                                             [morphy2k](https://github.com/morphy2k)                                              |                                             [tamird](https://github.com/tamird)                                              |                                             [paul-at-start](https://github.com/paul-at-start)                                              |                                             [vkrol](https://github.com/vkrol)                                             |

## License

MIT
