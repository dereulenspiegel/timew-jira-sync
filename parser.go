package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type TimewarriorConfig struct {
	Name  string
	Value string
}

type TimewarriorTime time.Time

const timewarriorLayout = "20060102T150405Z"

func (t *TimewarriorTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(timewarriorLayout, s)
	*t = TimewarriorTime(nt)
	return
}

func (t TimewarriorTime) Time() time.Time {
	return time.Time(t)
}

type TimewarriorInterval struct {
	Id    uint64          `json:"id"`
	Start TimewarriorTime `json:"start"`
	End   TimewarriorTime `json:"end"`
	Tags  []string        `json:"tags"`
}

func ParseTimewarrior(in io.Reader) (config []TimewarriorConfig, intervals []TimewarriorInterval, err error) {

	inBuf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read input: %w", err)
	}

	inParts := strings.SplitN(string(inBuf), "\n\n", 2)
	if len(inParts) != 2 {
		return nil, nil, fmt.Errorf("failed to parse input")
	}
	sc := bufio.NewScanner(strings.NewReader(inParts[0]))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "" {
			// This should not happen, but we tolerate it
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("found invalid config line: %s", line)
		}
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config = append(config, TimewarriorConfig{
			Name:  name,
			Value: value,
		})
	}
	if err = sc.Err(); err != nil {
		return nil, nil, err
	}

	// Reader should now be advanced to the JSON encoded intervals
	intervals = []TimewarriorInterval{}
	if err := json.Unmarshal([]byte(inParts[1]), &intervals); err != nil {
		return nil, nil, fmt.Errorf("failed to decode intervals: %w", err)
	}

	return
}
