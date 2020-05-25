package hw10_program_optimization //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

var (
	reg        *regexp.Regexp
	emails     [100_000]string
	p          fastjson.Parser
	ErrNoEmail = errors.New("line does not contain email")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	// Get Users
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		v, err := p.Parse(line)
		if err != nil {
			return nil, err
		}
		emails[i] = string(v.GetStringBytes("Email"))
		if (emails[i] == "") || (emails[i] == " ") {
			return nil, fmt.Errorf("%w: %s", ErrNoEmail, line)
		}
	}

	// Count domains
	reg, err = regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	result := make(DomainStat)

	for _, e := range emails {
		if e == "" {
			break
		}
		if reg.MatchString(e) {
			result[strings.ToLower(strings.SplitN(e, "@", 2)[1])]++
		}
	}
	return result, nil
}
