# revive

Fast, configurable, extensible, and beautiful linter for Go.

<p align="center">
  <img src="./assets/logo.png" alt="" width="200">
</p>

## Usage

Revive is **configurable** linter which you can fit your needs.

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

## License

MIT
