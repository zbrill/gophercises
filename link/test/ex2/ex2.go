package main

import (
	"fmt"
	"link"
)

func main() {
	links := link.Parse("ex2.html")
	fmt.Println(links)
}
