# revive

Fast, configurable, extensible, flexible, and beautiful linter for Go.

<p align="center">
  <img src="./assets/logo.jpg" alt="" width="300">
</p>

Here's how `revive` is different from `golint`:

* Allows you to enable or disable rules using a configuration file.
* Allows you to configure the linting rules with a TOML file.
* Provides functionality to disable a specific rule or the entire linter for a file or a range of lines.
  * `golint` allows this only for generated files.
* Provides multiple formatters which let you customize the output.
* Allows you to customize the return code for the entire linter or based on the failure of only some rules.
* Open for addition of new rules or formatters.
* Provides more rules compared to `golint`.
* Faster. It runs the rules over each file in a separate goroutine.

## Usage

Since the default behavior of `revive` is compatible with `golint`, without providing any additional flags, the only difference you'd notice is faster execution.

### Command Line Flags

`revive` accepts three command line parameters:

* `-config [PATH]` - path to config file in TOML format.
* `-exclude [PATTERN]` - pattern for files/directories/packages to be excluded for linting. You can specify the files you want to exclude for linting either as package name (i.e. `github.com/mgechev/revive`), list them as individual files (i.e. `file.go`), directories (i.e. `./foo/...`), or any combination of the three.
* `-formatter [NAME]` - formatter to be used for the output. The currently available formatters are:
  * `default` - will output the failures the same way that `golint` does.
  * `json` - outputs the failures in JSON format.
  * `friendly` - outputs the failures when found. Shows summary of all the failures.
  * `stylish` - formats the failures in a table. Keep in mind that it doesn't stream the output so it might be perceived as slower compared to others.

### Sample Invocations

```shell
revive -config revive.toml -exclude file1.go -exclude file2.go -formatter friendly github.com/mgechev/revive package/...
```

* The command above will use the configuration from `revive.toml`
* `revive` will ignore `file1.go` and `file2.go`
* The output will be formatted with the `friendly` formatter
* The linter will analyze `github.com/mgechev/revive` and the files in `package`

### Configuration

`revive` can be configured with a TOML file. Here's a sample configuration with explanation for the individual properties:

```toml
# Ignores files with "GENERATED" header, similar to golint
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

List of all available rules. The rules ported from `golint` are left unchanged and indicated in the `golit` column.

| Name                  | Config | Description                                                      | `golint` |
| --------------------- | :----: | :--------------------------------------------------------------- | :------: |
| `blank-imports`       |  N/A   | Disallows blank imports                                          |   yes    |
| `context-arguments`   |  N/A   | `context.Context` should be the first argument of a function.    |   yes    |
| `context-key-type`    |  N/A   | Disallows the usage of basic types in `context.WithValue`.       |   yes    |
| `dot-imports`         |  N/A   | Forbids `.` imports.                                             |   yes    |
| `error-return`        |  N/A   | The error return parameter should be last.                       |   yes    |
| `error-strings`       |  N/A   | Conventions around error strings.                                |   yes    |
| `errors`              |  N/A   | Naming of error variables.                                       |   yes    |
| `exported`            |  N/A   | Naming and commenting conventions on exported symbols.           |   yes    |
| `if-return`           |  N/A   | Redundant if when returning an error.                            |   yes    |
| `increment-decrement` |  N/A   | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            |   yes    |
| `names`               |  N/A   | Naming rules.                                                    |   yes    |
| `package-comments`    |  N/A   | Package commenting conventions.                                  |   yes    |
| `range`               |  N/A   | Prevents redundant variables when iterating over a collection.   |   yes    |
| `receiver-name`       |  N/A   | Conventions around the naming of receivers.                      |   yes    |
| `time-name`           |  N/A   | Conventions around the naming of time variables.                 |   yes    |
| `unexpected-return`   |  N/A   | Warns when a public return is from unexported type.              |   yes    |
| `unexpected-return`   |  N/A   | Conventions around variable declarations.                        |   yes    |
| `indent-error-flow`   |  N/A   | Prevents redundant else statements.                              |   yes    |
| `errorf`              |  N/A   | Should replace `error.New(fmt.Sprintf())` with `error.Errorf()`  |   yes    |
| `argument-limit`      |  int   | Specifies the maximum number of arguments a function can receive |    no    |
| `cyclomatic`          |  int   | Sets restriction for maximum Cyclomatic complexity.              |    no    |
| `max-public-structs`  |  int   | The maximum number of public structs in a file.                  |    no    |
| `file-header`         | string | Header which each file should have.                              |    no    |

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

* Enable the rule with name `ban-struct-name`. The `Name()` method of our rule should return a string which matches `ban-struct-name`.
* Configure the rule with the argument `Foo`. The list of arguments will be passed to `Apply(*File, Arguments)` together with the target file we're linting currently.

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

## License

MIT
