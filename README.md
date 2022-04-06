# ECO - Enviroment Configurations to Go Structs

[![Test Coverage](https://api.codeclimate.com/v1/badges/7f1d9250d3e0ce38cda5/test_coverage)](https://codeclimate.com/github/orkungursel/go-eco/test_coverage)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/orkungursel/go-eco)](https://pkg.go.dev/mod/github.com/orkungursel/go-eco)

This is a tiny package to map the environment variables to Go structs w/o the need of a config file. It supports `prefixing`, `nesting` and `default values` in your config structs. Please check out the examples below.

## Supported Types

- [x] `string`
- [x] `int`, `Uint`, `int8`, `Uint8`, `int16`, `Uint16`, `int32`, `Uint32`, `int64`, `Uint64`
- [x] `float32`, `float64`
- [x] `bool`
- [x] `[]string`
- [x] `[]int`, `[]int64`
- [x] `[]float32`, `[]float64`

## Installation

```bash
go get github.com/orkungursel/go-eco
```

## Examples

### Using Global API

```go
package main

import (
	"fmt"

	"github.com/orkungursel/go-eco"
)

type Config struct {
	Port int    `env:"PORT" default:"8080"`
	Host string `env:"HOST" default:"localhost"`
}

func main() {
	config := Config{}

	if err := eco.Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```

```bash
$ PORT=8081 go run main.go

{Port:8081 Host:localhost}
```

### Using Local Instance

```go
...
	config := Config{}

	e := eco.New()
	e.SetPrefix("APP")

	if err := e.Unmarshal(&config); err != nil {
		panic(err)
	}
...
```

```bash
$ APP_PORT=8081 go run main.go

{Port:8081 Host:localhost}
```

### Nested Structs

```go
package main

import (
	"fmt"

	"github.com/orkungursel/go-eco"
)

type Config struct {
	Port int    `env:"PORT" default:"8080"`
	Host string `env:"HOST" default:"localhost"`
	Logger struct {
		Level string `env:"LEVEL" default:"info"`
	}
}

func main() {
	config := Config{}

	if err := eco.Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```

```bash
$ PORT=8081 LOGGER_LEVEL=debug go run main.go

{Port:8081 Host:localhost Logger: {Level:debug}}
```

## API

### SetPrefix

```go
func SetPrefix(prefix string)
```
    SetPrefix sets the prefix for the environment variables.

    By default, the prefix is empty.

<details>
  <summary>Example</summary>

```go
...

func main() {
	config := Config{}

	if err := eco.SetPrefix("custom_prefix").Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```

```bash
$ CUSTOM_PREFIX_PORT=8081 go run main.go

{Port:8081 Host:localhost}
```
</details>

### SetArraySeparator

```go
func SetArraySeparator(sep string)
```
    SetArraySeparator sets the separator for array values.

    By default, the separator is `,`.

<details>
  <summary>Example</summary>

```go
package main

import (
	"fmt"

	"github.com/orkungursel/go-eco"
)

type Config struct {
    Foo []string `default:"foo,bar,baz"`
}

func main() {
	config := Config{}

	if err := eco.SetArraySeparator(",").Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```
</details>

### SetEnvNameSeparator

```go
func SetEnvNameSeparator(envNameSeparator string) *eco {
```
    SetEnvNameSeparator sets the separator for the environment variable names.

    By default, the separator is `_`.

<details>
  <summary>Example</summary>

```go
package main

import (
	"fmt"

	"github.com/orkungursel/go-eco"
)

type Config struct {
	Sub struct {
		Foo string
	}
}

func main() {
	config := Config{}

	if err := eco.SetEnvNameSeparator("__").Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```

```bash
$ SUB__FOO=bar go run main.go

{Sub: {Foo:bar}}
```
</details>



### Unmarshal

```go
func Unmarshal(v interface{}) error
```

    Unmarshal takes a pointer to a struct and unmarshals the environment variables to the struct.

## License

This project is licensed under the [MIT](LICENSE) License.