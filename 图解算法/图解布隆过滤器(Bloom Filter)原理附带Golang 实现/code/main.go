package main

import "fmt"

func main() {
	filter, err := NewBloomFilter(3, 10000, nil)
	if err != nil {
		return
	}
	err = filter.Add([]byte("hello"))
	filter.Add([]byte("world"))
	filter.Add([]byte("baby"))

	ok, _ := filter.IsExists([]byte("hello"))
	if ok {
		fmt.Println("the hello is exists")
	}

	ok, _ = filter.IsExists([]byte("bob"))
	if ok {
		fmt.Println("the bob is exists")
	}

}
