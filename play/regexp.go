package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println("Regexp test")
	validID := regexp.MustCompile(`\/test(?:/([0-9]+))(?:/([0-9]+))*`)
	fmt.Println(validID.MatchString("/test/zip"))
	fmt.Println(validID.FindAllString("/test/675", -1))
	fmt.Println(validID.FindAllStringSubmatch("/test/675", -1))
	fmt.Println(validID.FindAllStringSubmatch("/test/675/745", -1))
	fmt.Println(validID.MatchString("/atest"))
}
