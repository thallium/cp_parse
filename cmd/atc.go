/*
Copyright © 2021 Gengchen Tuo <tuogengchen@gmail.com>

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

var atcoderProbInfo = &util.ProblemInfo{
	regexp.MustCompile(`<span class="h2">\s*(.+?)\s*?<`),
	regexp.MustCompile(`Sample Input [\s\S]*?<pre>([\s\S]*?)</pre>`),
	regexp.MustCompile(`Sample Output [\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var atcoderContestInfo = &util.ContestInfo{
	regroup.MustCompile(`<a href="(?P<link>.+?)">(?P<index>\w{1,2})</a>`),
	regexp.MustCompile(`<a class="contest-title".*?>(.+?)</a>`),
	atcoderProbInfo,
	`https://atcoder.jp`,
}

var atcArgRegStr = map[string]int{
	`^https://atcoder.jp/contests/\w+?/tasks/\w+$`: util.ProbURL,
	`^https://atcoder.jp/contests/\w+?$`:           util.ContestURL,
	`^([\w-]+)_[[:alpha:]]$`:                       util.ProbID,
	`^[\w-]+$`:                                     util.ContestID,
}

func atcArgToURL(arg string, ty int, match []string) (string, int) {
	if ty == util.ContestID {
		return `https://atcoder.jp/contests/` + arg + `/tasks`, 1
	} else if ty == util.ProbID {
		return fmt.Sprintf(`https://atcoder.jp/contests/%v/tasks/%v`, match[1], arg), 0
	} else if ty == util.ContestURL {
		return arg + `/tasks`, 1
	} else {
		return arg, ty
	}
}

// atcCmd represents the atc command
var atcCmd = &cobra.Command{
	Use:   "atc",
	Short: "Parse problems/contests from atcoder.jp",
	Long: `Usage: 
    cp_parse atc [contest/problem]
Contest can be:
    URL             e.g. https://atcoder.jp/contests/arc112
    Contest id      e.g. arc112

Problem can be:
    URL             e.g. https://atcoder.jp/contests/arc112/tasks/arc112_a
    Problem id      e.g. arc112_a

Example:
    cp_parse atc https://atcoder.jp/contests/arc112/tasks/arc112_a
    cp_parse atc arc112`,
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
		URL, ty := util.ProcessArg(args[0], &atcArgRegStr, atcArgToURL)
		if ty == 0 {
			err = util.ParseProblem(URL, dir, atcoderProbInfo)
		} else if ty == 1 {
			err = util.ParseContest(URL, dir, atcoderContestInfo)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(atcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// atcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// atcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
