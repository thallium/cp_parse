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

	"github.com/spf13/cobra"
	"github.com/thallium/cp_parser/util"
)

// cfCmd represents the cf command
var cfCmd = &cobra.Command{
	Use:   "cf",
	Short: "Parse problems/contests from codeforces.com",
    Long:`Usage: 
    cp_parser cf [contest/problem]
Contest can be:
    URL             e.g. https://codeforces.com/contest/1490
    Contest id      e.g. 1490

Problem can be:
    URL             e.g. https://codeforces.com/contest/1490/problem/A
                         https://codeforces.com/problemset/problem/1490/A
    Problem id      e.g. 1490A, 1490a

Example:
    cp_parser cf https://codeforces.com/contest/1490/problem/A
    cp_parser cf 1490`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Exact one argument should be provided, but get %v arguments\n", len(args))
			os.Exit(1)
		}
		dir, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		URL, ty := cfProcessArg(args[0])
		if ty == 0 {
			err = util.ParseProblem(URL, dir, util.CfProb)
		} else if ty == 1 {
			err = util.ParseContest(URL, dir, util.CfContest)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var cfArgRegStr = map[string]int{
	`^https://codeforces.com/problemset/problem/\d+/[[:alpha:]]\d?$`: 0,
	`^https://codeforces.com/contest/\d+/problem/[[:alpha:]]\d?$`:    0,
	`^https://codeforces.com/contest/\d+$`:                           1,
	`^(\d+)([[:alpha:]]\d?)$`:                                        2,
	`^\d+$`:                                                          3,
}

func cfProcessArg(arg string) (string, int) {
	for regStr, ty := range cfArgRegStr {
		reg := regexp.MustCompile(regStr)
		match := reg.FindStringSubmatch(arg)
		if len(match) != 0 && match[0] != "" {
			if ty == 3 {
				return `https://codeforces.com/contest/` + arg, 1
			} else if ty == 2 {
				return fmt.Sprintf(`https://codeforces.com/problemset/problem/%v/%v`, match[1], match[2]), 0
			} else {
				return arg, ty
			}
		}
	}
	fmt.Println("Invalid problem/contest")
	os.Exit(1)
	return "", 0
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
