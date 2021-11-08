package util

import (
	"regexp"

	"github.com/oriser/regroup"
)
type ContestInfo struct {
	ProbsRg  *regroup.ReGroup
	NameRg   *regexp.Regexp
	ProbInfo *ProblemInfo
	BaseURL  string
}

type ProblemInfo struct {
	NameRg, InputRg, OutputRg *regexp.Regexp
}
