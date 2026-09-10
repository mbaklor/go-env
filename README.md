# ENV

A small package to read environment variables into a struct.

I wanted something that works like json unmarshalling, using struct tags to
indicate the name of the environment variable.


```shell
go get github.com/mbaklor/env
```


## How to use

```go
package main

import (
        "fmt"

        "github.com/mbaklor/env"
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
