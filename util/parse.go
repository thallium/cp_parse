package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var BodyByExtension string

func GetBody(URL string) ([]byte, error) {
	if len(BodyByExtension) != 0 {
		return []byte(BodyByExtension), nil
	}
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func findSample(body []byte, info *ProblemInfo) (input [][]byte, output [][]byte, err error) {
	in := info.InputRg.FindAllSubmatch(body, -1)
	ou := info.OutputRg.FindAllSubmatch(body, -1)
	if in == nil || ou == nil || len(in) != len(ou) {
		return nil, nil, fmt.Errorf("parse sample failed")
	}
	processString := func(s []byte) []byte {
		brRg := regexp.MustCompile(`<br ?/>`)
		s = brRg.ReplaceAll(s, []byte("\n"))
		if s[0] == '\n' {
			s = s[1:]
		}
		if s[len(s)-1] != '\n' {
			s = append(s, '\n')
		}
		return s
	}
	for _, s := range in {
		input = append(input, processString(s[1]))
	}
	for _, s := range ou {
		output = append(output, processString(s[1]))
	}
	return
}

func findName(body []byte, info *ProblemInfo) (string, error) {
	name := info.NameRg.FindSubmatch(body)
	if len(name) == 0 {
		return "", fmt.Errorf("Can't find problem name!\n")
	}
	return string(name[1]), nil
}

func ParseProblem(URL, path string, info *ProblemInfo) error {
	body, err := GetBody(URL)
	if err != nil {
		return err
	}
	name, err := findName(body, info)
	name = strings.Replace(name, `<br/>`, " ", -1)
	if err != nil {
		return err
	}
	input, output, err := findSample(body, info)

	if err != nil {
		return err
	}

	// delete old input and output files
	old_in, err_in := filepath.Glob("*.in")
	old_out, err_out := filepath.Glob("*.out")
	if err_in != nil || err_out != nil {
		return err
	}
	delete_files := func(files []string) {
		for _, file := range files {
			os.Remove(file)
		}
	}
	delete_files(old_in)
	delete_files(old_out)

	for i := 0; i < len(input); i++ {
		fileIn := filepath.Join(path, fmt.Sprintf("%v.in", i+1))
		fileOut := filepath.Join(path, fmt.Sprintf("%v.out", i+1))

		e := ioutil.WriteFile(fileIn, input[i], 0644)
		if e != nil {
			return e
		}
		e = ioutil.WriteFile(fileOut, output[i], 0644)
		if e != nil {
			return e
		}
	}
	if len(input) == 1 {
		fmt.Printf("Parsed %v with 1 sample\n", name)
	} else {
		fmt.Printf("Parsed %v with %v samples\n", name, len(input))

	}
	return nil
}

func ParseContest(URL, path string, info *ContestInfo) error {
	body, err := GetBody(URL)
	if err != nil {
		return err
	}
	name := string(info.NameRg.FindSubmatch(body)[1])
	fmt.Println("Parsing contest " + name)
	name = strings.Replace(name, " ", "_", -1)
	nonWord := regexp.MustCompile(`\W`)
	name = nonWord.ReplaceAllString(name, "")
	os.Mkdir(name, 01755)
	os.Chdir(filepath.Join(path, name))
	type prob struct {
		Link  string `regroup:"link"`
		Index string `regroup:"index"`
	}
	target := &prob{}
	rets, err := info.ProbsRg.MatchAllToTarget(string(body), -1, target)
	if err != nil {
		return err
	}
	for _, suffix := range rets {
		index := suffix.(*prob).Index
		os.Mkdir(index, 01755)
		ioutil.WriteFile(filepath.Join(path, name, index, index+".cpp"), nil, 0644)
		err := ParseProblem(info.BaseURL+suffix.(*prob).Link, filepath.Join(path, name, index), info.ProbInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetWebsiteName(url string) string {
	parts := strings.Split(url, ".")
	return parts[len(parts)-2]
}
