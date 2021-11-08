package util

import (
	"fmt"
	"os"
	"regexp"
)

const (
	ProbURL int = iota
	ContestURL
	ProbID
	ContestID
)

func ProcessArg(arg string, ArgRegStr *map[string]int,
	argToURL func(string, int, []string) (string, int)) (string, int) {

	for regStr, ty := range *ArgRegStr {
		reg := regexp.MustCompile(regStr)
		match := reg.FindStringSubmatch(arg)
		if len(match) != 0 && match[0] != "" {
			return argToURL(arg, ty, match)
		}
	}
	fmt.Println("Invalid problem/contest")
	os.Exit(1)
	return "", 0
}
