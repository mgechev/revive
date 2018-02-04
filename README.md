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

## License

MIT
