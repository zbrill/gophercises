package main

import (
	"fmt"
	"link"
)

func main() {
	links := link.Parse("ex4.html")
	fmt.Println(links)
}
