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

func GetBody(URL string) ([]byte, error) {
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
func findSample(body []byte, info *problemInfo) (input [][]byte, output [][]byte, err error) {
	in := info.inputRg.FindAllSubmatch(body, -1)
	ou := info.outputRg.FindAllSubmatch(body, -1)

	if in == nil || ou == nil || len(in) != len(ou) {
		return nil, nil, fmt.Errorf("Parse sample failed")
	}
	processString := func(s []byte) []byte {
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

func findName(body []byte, info *problemInfo) (string, error) {
	name := info.nameRg.FindSubmatch(body)
	if len(name) == 0 {
		return "", fmt.Errorf("Can't find problem name!\n")
	}
	return string(name[1]), nil
}

func ParseProblem(URL, path string, info *problemInfo) error {
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
	fmt.Printf("Parsed %v with %v sample(s)\n", name, len(input))
	return nil
}

func ParseContest(URL, path string, info *contestInfo) error {
	body, err := GetBody(URL)
	if err != nil {
		return err
	}
	name := string(info.nameRg.FindSubmatch(body)[1])
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
	rets, err := info.probsRg.MatchAllToTarget(string(body), -1, target)
	if err != nil {
		return err
	}
	for _, suffix := range rets {
		index := suffix.(*prob).Index
		os.Mkdir(index, 01755)
		ioutil.WriteFile(filepath.Join(path, name, index, index+".cpp"), nil, 0644)
		err := ParseProblem(info.baseURL+suffix.(*prob).Link, filepath.Join(path, name, index), info.probInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
