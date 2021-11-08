/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/oriser/regroup"
	"github.com/spf13/cobra"
	"github.com/thallium/cp_parse/util"
)

var probInfo = &util.ProblemInfo{
	NameRg:   regexp.MustCompile(`<div class="title">([[:print:]]+?)<`),
	InputRg:  regexp.MustCompile(`class="input"[\s\S]*?<pre>([\s\S]*?)</pre>`),
	OutputRg: regexp.MustCompile(`class="output"[\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var contestInfo = &util.ContestInfo{
	ProbsRg:  regroup.MustCompile(`<td class="id">\s*?<a href="(?P<link>[[:print:]]+?)">\s*(?P<index>\w+?)\s*<`),
	NameRg:   regexp.MustCompile(`<table class="rtable ">[\s\S]*?<a.*?href.*?>([[:print:]]+?)</a>`),
	ProbInfo: probInfo,
	BaseURL:  `https://codeforces.com`,
}

var argRegStr = map[string]int{
	`^https://codeforces.com/problemset/problem/\d+/[[:alpha:]]\d?$`: util.ProbURL,
	`^https://codeforces.com/contest/\d+/problem/[[:alpha:]]\d?$`:    util.ProbURL,
	`^https://codeforces.com/contest/\d+$`:                           util.ContestURL,
	`^(\d+)([[:alpha:]]\d?)$`:                                        util.ProbID,
	`^\d+$`:                                                          util.ContestID,
}

func argToURL(arg string, ty int, match []string) (string, int) {
	if ty == 3 {
		return `https://codeforces.com/contest/` + arg, 1
	} else if ty == 2 {
		return fmt.Sprintf(`https://codeforces.com/problemset/problem/%v/%v`, match[1], match[2]), 0
	} else {
		return arg, ty
	}
}

// cfCmd represents the cf command
var cfCmd = &cobra.Command{
	Use:   "cf",
	Short: "Parse problems/contests from codeforces.com",
	Long: `Usage: 
    cp_parse cf [contest/problem]
Contest can be:
    URL             e.g. https://codeforces.com/contest/1490
    Contest id      e.g. 1490

Problem can be:
    URL             e.g. https://codeforces.com/contest/1490/problem/A
                         https://codeforces.com/problemset/problem/1490/A
    Problem id      e.g. 1490A, 1490a

Example:
    cp_parse cf https://codeforces.com/contest/1490/problem/A
    cp_parse cf 1490`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Exact one argument should be provided, but get %v arguments\n", len(args))
			os.Exit(1)
		}
		dir, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		URL, ty := util.ProcessArg(args[0], &argRegStr, argToURL)
		if ty == 0 {
			err = util.ParseProblem(URL, dir, probInfo)
		} else if ty == 1 {
			err = util.ParseContest(URL, dir, contestInfo)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cfCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
