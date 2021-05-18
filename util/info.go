package util

import (
	"regexp"

	"github.com/oriser/regroup"
)

type contestInfo struct {
	probsRg  *regroup.ReGroup
	nameRg   *regexp.Regexp
	probInfo *problemInfo
	baseURL  string
}

type problemInfo struct {
	nameRg, inputRg, outputRg *regexp.Regexp
}

var CfProb = &problemInfo{
	regexp.MustCompile(`<div class="title">([[:print:]]+?)<`),
	regexp.MustCompile(`class="input"[\s\S]*?<pre>([\s\S]*?)</pre>`),
	regexp.MustCompile(`class="output"[\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var CfContest = &contestInfo{
	regroup.MustCompile(`<td class="id">\s*?<a href="(?P<link>[[:print:]]+?)">\s*(?P<index>\w+?)\s*<`),
	regexp.MustCompile(`<table class="rtable ">[\s\S]*?<a.*?href.*?>([[:print:]]+?)</a>`),
	CfProb,
	`https://codeforces.com`,
}

var AtcoderProb = &problemInfo{
	regexp.MustCompile(`<span class="h2">\s*?([[:print:]]+?)\s*?<`),
	regexp.MustCompile(`Sample Input [\s\S]*?<pre>([\s\S]*?)</pre>`),
	regexp.MustCompile(`Sample Output [\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var AtcoderContest = &contestInfo{
	regroup.MustCompile(`<a href="(?P<link>[[:print:]]+?)">(?P<index>\w{1,2})</a>`),
	regexp.MustCompile(`<a class="contest-title".*?>(.+?)</a>`),
	AtcoderProb,
	`https://atcoder.jp`,
}

var KattisProb = &problemInfo{
	regexp.MustCompile(`<div class="headline-wrapper"><h1>([[:print:]]+?)</h1>`),
	regexp.MustCompile(`Sample Input[\s\S]*?<pre>([\s\S]*?)</pre>`),
	regexp.MustCompile(`Sample Output[\s\S]*?<pre>[\s\S]*?</pre>[\s\S]*?<pre>([\s\S]*?)</pre>`),
}

var KattisContest = &contestInfo{
	regroup.MustCompile(`<th class="problem_letter">(?P<index>\w+?)<[\s\S]*?href="(?P<link>[[:print:]]+?)">`),
	regexp.MustCompile(`<div class="header-title">([[:print:]]+?)</div>`),
	KattisProb,
	`https://open.kattis.com`,
}
