// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "RAJAN?"
	fmt.Println(val(str))
	fmt.Println(val("WHAT ARE YOU ?DOING"))
	fmt.Println(val("YELL AT HIM"))
	fmt.Println(val("Hello"))

	fmt.Println(val(" "))

}
func val(s string) interface{} {
	str := strings.TrimSpace(s)
	if str == "" {
		return "Fine. Be that way!"
	} else if str == strings.ToUpper(str) && strings.HasSuffix(str, "?") {
		return "Calm down, I know what I'm doing!"
	} else if strings.HasSuffix(str, "?") {
		return "Sure"
	} else if str == strings.ToUpper(str) {
		return "Whoa, chill out!"
	} else {
		return "Whatever"
	}

}
