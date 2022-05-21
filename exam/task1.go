package main

import "fmt"

type request struct {
	first  int
	second int
}

func product(req request) int {
	return req.first * req.second
}

func main() {
	req := request{first: 2, second: 3}
	fmt.Println(product(req))
}
