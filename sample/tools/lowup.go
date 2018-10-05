package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(byte('a'), byte('z'), byte('A'), byte('Z')) // 97 122 65 90
	fmt.Println(97 - 65, 122 - 90) // 25 25
	fmt.Println(string(byte('a') - 32))
	fmt.Println(toCamel("hello_word_Abc", true))
	fmt.Println(toCamel("hello_word_Abc", false))
	fmt.Println(toCamel("User", false))
	//mapType := "xx/time.Time"
	//if pos := strings.LastIndexByte(mapType, '.'); pos > 0 {
	//	fmt.Println(mapType[0:pos])
	//}
	//fmt.Println(toType("time.Time"))
	//fmt.Println(toType("xx/time.Time"))
	//fmt.Println(toType("int"))

}

func toCamel(src string, aheadUp bool) string {
	slices := strings.Split(src, "_")
	ss := ""
	for i, slice := range slices {
		p := []byte(slice)
		if i == 0 {
			if aheadUp && p[0] >= 97 && p[0] <= 122 {
				p[0] = p[0] - 32
			} else if !aheadUp && p[0] >= 65 && p[0] <= 90 {
				p[0] = p[0] + 32
			}
		} else {
			if p[0] >= 97 && p[0] <= 122 {
				p[0] = p[0] - 32
			}
		}
		ss += string(p)
	}
	return ss
}

func toType(src string) string {
	ss := src
	if pos := strings.LastIndexByte(src, '/'); pos > 0 {
		ss = src[pos+1:]
	}
	return ss
}