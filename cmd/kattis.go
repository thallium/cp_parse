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

// kattisCmd represents the kattis command
var kattisCmd = &cobra.Command{
	Use:   "kattis",
	Short: "Parse problems/contests from open.kattis.com",
    Long:`Usage: 
    cp_parser kattis [contest/problem]
Contest can be:
    URL             e.g. https://open.kattis.com/contests/nar20practice14
                         https://open.kattis.com/contests/nar20practice14/problems
Problem can be:
    URL             e.g. https://open.kattis.com/contests/nar20practice14/problems/arrayofdiscord
                         https://open.kattis.com/problems/sequences
    Problem id      e.g. sequences, 10kindsofpeople

Example:
    cp_parser kattis 10kindsofpeople
    cp_parser kattis https://open.kattis.com/contests/nar20practice14`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Printf("Exact one argument should be provided, but get %v arguments\n", len(args))
			os.Exit(1)
		}
		dir, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		if args[0][len(args[0])-1] == '/' {
			args[0] = args[0][:len(args[0])-1]
		}
		URL, ty := kattisProcessArg(args[0])
		if ty == 0 {
			err = util.ParseProblem(URL, dir, util.KattisProb)
		} else if ty == 1 {
			err = util.ParseContest(URL, dir, util.KattisContest)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var kattisArgRegStr = map[string]int{
	`^https://open.kattis.com/contests/\w+/problems/\w+$`: 0,
	`^https://open.kattis.com/problems/\w+$`:              0,
	`^https://open.kattis.com/contests/\w+/problems$`:     1,
	`^https://open.kattis.com/contests/\w+$`:              3,
	`^\w+$`:                                               2,
}

func kattisProcessArg(arg string) (string, int) {
	for regStr, ty := range kattisArgRegStr {
		reg := regexp.MustCompile(regStr)
		match := reg.FindStringSubmatch(arg)
		if len(match) != 0 && match[0] != "" {
			if ty == 2 {
				return `https://open.kattis.com/problems/` + arg, 0
			} else if ty == 3 {
				return arg + `/problems`, 1
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
	rootCmd.AddCommand(kattisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kattisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kattisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
