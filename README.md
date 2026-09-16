# go-ENV

A small package to read environment variables into a struct.

I wanted something that works like json unmarshalling, using struct tags to
indicate the name of the environment variable.

```shell
go get github.com/mbaklor/go-env
```

## How to use

```go
package main

import (
        "fmt"

        "github.com/mbaklor/go-env"
)

type Config struct {
        DatabaseHost string `env:"DATABASE_HOST"`
        DatabasePort int    `env:"DATABASE_PORT"`
}

func main() {
        var c Config

        err := env.Load(&c)
        if err != nil {
                fmt.Printf("Failed to load env vars: %s\n", err)
                return
        }
        fmt.Printf("Database connection: %s:%d\n", c.DatabaseHost, c.DatabasePort)
}
```

Set the environment variables and run the program

```shell
DATABASE_HOST=127.0.0.1 DATABASE_PORT=5432 go run main.go
```

You should see this output

```shell
Database connection: 127.0.0.1:5432
```

### Auto load .env file

I added an `autofile` package which should auto load a `.env` file from your
current directory.

To use this functionality, in the imports in your project add

```go
import _ "github.com/mbaklor/go-env/autofile"
```

and the env file should load to the program's environment variables.

Note that any error encountered during the autoload are saved in
`autofile.AutoLoadFileErr` so you might want to check that if something isn't
working as expected.

## Secondary Use: Load .env Files

This package can also be used to load `.env` files.

For exmple if you have

```env
# .env

DATABASE_HOST=127.0.0.1
DATABASE_PORT=5432
```

You can add to the `main` function above

```go

func main() {
        var c Config

        if err := env.LoadFiles(); err != nil {
                fmt.Printf("Failed to set env vars from file: %s\n", err)
        }

        err := env.Load(&c)
        // continue as example above
```

and after a

```shell
go run .
```

the shell output should still be

```shell
Database connection: 127.0.0.1:5432
```

## Roadmap

### Currently supported

- use `env:"-"` to skip a field
- string
- int
- bool
- slice of any of the above, including use of a custom delimiter with `delim=`
- `time.Duration` specifically because it isn't a struct(??)
- nested struct of any supported types
- pointer to any supported type
- anything implementing the `env.Unmarshaler` interface
- anything implementing the `encoding.TextUnmarshaler` interface

- loading env files
- autoloading a .env file

### Planned support

- we'll see what else I feel is missing over time

## Other (frankly better) projects

I will admit this project is mostly a vanity "let's see how I do this" kind of
deal. I mostly just hyperfixated on it one night and am now running this
fixation until I ultimately burn through the obsession.

I didn't look at the code for or try to use any of the other myriad env struct
tag loaders, but after I finished my first working version I was curious and
found about 100 `goenv` `go-env` or simply `env` packages.

My assumption is any one of those will be more mature and fleshed out than this
project, so please go check them out.

The other similar project is [joho/godotenv](https://github.com/joho/godotenv)
which again I haven't personally used but from a cursory glance at the readme
it looks to be doing what I do but better, please go check that project out.
