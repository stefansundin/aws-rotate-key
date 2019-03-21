package replacer

import (
	"io/ioutil"
	"strings"
)

type Replacer strings.Replacer

func New(pairs map[string]string) *Replacer {
	s := make([]string, 0)
	for k, v := range pairs {
		s = append(s, k, v)
	}
	return (*Replacer)(strings.NewReplacer(s...))
}

func (r *Replacer) RewriteFile(filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return err
	}
	return writeFile(filename, r.replace(string(data)))
}

func (r *Replacer) replace(s string) string {
	return (*strings.Replacer)(r).Replace(s)
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0600)
}
