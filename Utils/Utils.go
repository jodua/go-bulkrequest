package utils

import (
	"log"
	"math/rand"
	"regexp"
)

func Flatten(listOfLists [][]any) []any {
	var flattened []any
	for _, list := range listOfLists {
		flattened = append(flattened, list...)
	}
	return flattened
}

func Take(list []string, offset int, length int) []string {
	if offset < 0 {
		offset = 0
	}
	if length < 0 {
		length = 0
	}
	if offset+length > len(list) {
		length = len(list) - offset
	}
	return list[offset : offset+length]
}

func ShuffleList(list []string) {
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
}

func IsValidUrl(url string) bool {
	re := `(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(:[a-zA-Z0-9]*)?\/?([a-zA-Z0-9\-\._\?\,\'/\\\+&amp;%\$#\=~])*`
	reg, err := regexp.Compile(re)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return false
	}
	return reg.MatchString(url)
}
