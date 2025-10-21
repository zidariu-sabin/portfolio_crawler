package main

import "fmt"

func main() {
	str := "testing error"

	fmt.Println(fmt.Errorf("creating directory %s", str))
}
