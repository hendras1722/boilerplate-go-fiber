package main

import (
	"fmt"
	"strconv"
)

func main() {
	page, _ := strconv.Atoi("invalid")
	page = max(1, page)
	
	limit, _ := strconv.Atoi("-5")
	limit = max(1, limit)

	fmt.Printf("page=%d limit=%d\n", page, limit)
}
