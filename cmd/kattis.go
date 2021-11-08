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

var kattisProbInfo = &util.ProblemInfo{
	NameRg:   regexp.MustCompile(`<div class="headline-wrapper"><h1>([[:print:]]+?)</h1>`),
	InputRg:  regexp.MustCompile(`Sample Input[\s\S]*?<pre>([\s\S]*?)</pre>`),
	OutputRg: regexp.MustCompile(`Sample Output[\s\S]*?<pre>[\s\S]*?</pre>[\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var kattisContestInfo = &util.ContestInfo{
	ProbsRg:  regroup.MustCompile(`<th class="problem_letter">(?P<index>\w+?)<[\s\S]*?href="(?P<link>[[:print:]]+?)">`),
	NameRg:   regexp.MustCompile(`<div class="header-title">([[:print:]]+?)</div>`),
	ProbInfo: kattisProbInfo,
	BaseURL:  `https://open.kattis.com`,
}

var kattisArgRegStr = map[string]int{
	`^https://open.kattis.com/contests/\w+/problems/\w+$`: util.ProbURL,
	`^https://\w+.kattis.com/problems/\w+$`:               util.ProbURL,
	`^https://\w+.kattis.com/problems$`:                   util.ContestURL,
	`^https://open.kattis.com/contests/\w+$`:              util.ContestURL,
	`^\w+$`:                                               util.ProbID,
}

func kattisArgToURL(arg string, ty int, match []string) (string, int) {
	if ty == util.ProbID {
		return `https://open.kattis.com/problems/` + arg, 0
	} else if ty == util.ContestURL {
		return arg + `/problems`, 1
	} else {
		return arg, ty
	}
}

// kattisCmd represents the kattis command
var kattisCmd = &cobra.Command{
	Use:   "kattis",
	Short: "Parse problems/contests from open.kattis.com",
	Long: `Usage: 
    cp_parse kattis [contest/problem]
Contest can be:
    URL             e.g. https://open.kattis.com/contests/nar20practice14
                         https://open.kattis.com/contests/nar20practice14/problems
Problem can be:
    URL             e.g. https://open.kattis.com/contests/nar20practice14/problems/arrayofdiscord
                         https://open.kattis.com/problems/sequences
    Problem id      e.g. sequences, 10kindsofpeople

Example:
    cp_parse kattis 10kindsofpeople
    cp_parse kattis https://open.kattis.com/contests/nar20practice14`,
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
		URL, ty := util.ProcessArg(args[0], &kattisArgRegStr, kattisArgToURL)
		if ty == 0 {
			err = util.ParseProblem(URL, dir, kattisProbInfo)
		} else if ty == 1 {
			err = util.ParseContest(URL, dir, kattisContestInfo)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
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
