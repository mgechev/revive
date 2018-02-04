# revive

Fast, configurable, extensible, flexible, and beautiful linter for Go.

<p align="center">
  <img src="./assets/logo.png" alt="" width="200">
</p>

Here's how `revive` is different from `golint`:

* Allows you to enable or disable rules using a configuration file.
* Allows you to configure the linting rules with a configuration file.
* Provides functionality to disable a specific rule or the entire linter for a file or a range of lines.
* Provides more rules compared to `golint`.
* Provides multiple formatters which let you customize the output.
* Allows you to customize the return code for the entire linter or based on the failure of only some rules.
* Open for addition of new rules or formatters.
* Faster since it runs the rules over each file in a separate goroutine.

## Usage

`Revive` is **configurable** linter which you can fit your needs. By default you can use `revive` with the default configuration options. This way the linter will work the same way `golint` does.

### Command Line Flags

Revive accepts only three command line parameters:

* `config` - path to config file in TOML format.
* `exclude` - pattern for files/directories/packages to be excluded for linting. You can specify the files you want to exclude for linting either as package name (i.e. `github.com/mgechev/revive`), list them as individual files (i.e. `file.go file2.go`), directories (i.e. `./foo/...`), or any combination of the three.
* `formatter` - formatter to be used for the output. The currently available formatters are:
  * `default` - will output the warnings the same way that `golint` does.
  * `json` - outputs the warnings in JSON format.
  * `cli` - formats the warnings in a table.

### Configuration

Revive can be configured with a TOML file

### Default Configuration

The default configuration of `revive` can be found at `defaults.toml`. This will enable all rules available in `golint` and use their default configuration (i.e. the config which is hardcoded in `golint`).

```shell
revive -config defaults.toml github.com/mgechev/revive
```

This will use `defaults.toml`, the `default` formatter, and will run linting over the `github.com/mgechev/revive` package.

### Recommended Configuration

```shell
revive -config config.toml -formatter cli github.com/mgechev/revive
```

This will use `config.toml`, the `cli` formatter, and will run linting over the `github.com/mgechev/revive` package.

## Extension

The tool can be extended with custom rules or formatters. This section contains additional information on how to implement such.

**To extend the linter with a custom rule or a formatter you'll have to push it to this repository**. This is due to the limited `-buildmode=plugin` support which [works only on Linux](https://golang.org/pkg/plugin/).

### Custom Rule

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}
```

The `Arguments` type is an alias of the type `[]interface{}` which means that you can pass arguments from any type to your rule. Let's suppose we have developed a rule called `BanStructNameRule` which disallow us to name a structure with given identifier. We can set the banned identifier by using the TOML configuration file:

```toml
[rule.ban-struct-name]
  arguments = ["Foo"]
```

With the snippet above we:

* Enable the rule `ban-struct-name` which is supposed to be the value returned by the `Name()` method of our rule.
* Pass an argument with value `"Foo"` to the `Apply` method of the rule once invoked with a file.

A sample rule implementation can be found [here](/rule/argument-limit.go).

### Custom Formatter

Each formatter needs to implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, RulesConfig) (string, error)
	Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration the enabled rules. The `Name()` method should return an string different from the names of the already existing rules.

For a sample formatter, take a look at [this file](/formatter/json.go).

## License

MIT
