package main

import (
	"fmt"
	"os"
	"gitee.com/li-jia666/serialize"
)

type ColorGroup struct {
	ID     int	`mytag:"color_id"`
	Name   string	`mytag:"color_name"`
	Colors []string	`mytag:"colors"`
}

func main() {
    group := ColorGroup{
        ID:     1,
        Name:   "Reds",
        Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
    }
    b, err := serialize.JsonMarshal(group)
    if err != nil {
        fmt.Println("error:", err)
    }
    os.Stdout.Write(b)
}