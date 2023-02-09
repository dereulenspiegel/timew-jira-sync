package timewarrior

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Config holds timewarrior configuration statements given at the beginning of the Std in stream
type Config struct {
	Name  string
	Value string
}

// Time represents time stamps as formatted by timewarrior
type Time time.Time

const timewarriorLayout = "20060102T150405Z"

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(timewarriorLayout, s)
	*t = Time(nt)
	return
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

// Interval represents timewarrior intervals as created with start and stop
type Interval struct {
	Id    uint64   `json:"id"`
	Start Time     `json:"start"`
	End   Time     `json:"end"`
	Tags  []string `json:"tags"`
}

// Parse parses a data stream like std in from timewarrior and returns the received configuration and intervals
func Parse(in io.Reader) (config []Config, intervals []Interval, err error) {

	inBuf, err := io.ReadAll(in)
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
		config = append(config, Config{
			Name:  name,
			Value: value,
		})
	}
	if err = sc.Err(); err != nil {
		return nil, nil, err
	}

	// Reader should now be advanced to the JSON encoded intervals
	intervals = []Interval{}
	if err := json.Unmarshal([]byte(inParts[1]), &intervals); err != nil {
		return nil, nil, fmt.Errorf("failed to decode intervals: %w", err)
	}

	return
}
