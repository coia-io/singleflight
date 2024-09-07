# singleflight

This is merely a wrapper for x/sync/singleflight, removing the annoying type assertions when in use.

## Install

```shell
go get github.com/coia-io/singleflight@v0.1.2
```

## Example usage

```go
package main

import (
	"fmt"
	"github.com/coia-io/singleflight"
)

type User struct {
	Name string
}

func (u *User) GetName() string { return u.Name }

type Cat struct {
	Color string
}

func (c *Cat) GetColor() string { return c.Color }

func main() {
	u, err, _ := singleflight.Do[User]("key", func() (any, error) { return User{"Zephyr"}, nil })
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u.GetName())
	
	c, err, _ := singleflight.Do[Cat]("key", func() (any, error) { return Cat{"Black"}, nil })
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.GetColor())
}
```
