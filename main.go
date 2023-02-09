package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/dereulenspiegel/timew-jira-sync/timewarrior"
	"gopkg.in/ini.v1"
)

var debug = false

var loggedTag = "logged"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to obtain users homedir: %w", err)
	}

	configDirPath := path.Join(homeDir, ".timewarrior", "extensions", "config.ini")
	cfg, err := ini.Load(configDirPath)
	if err != nil {
		log.Fatal("failed to read config file: %w", err)
	}
	jiraSection := cfg.Section("jira")

	wlAppender, err := NewDumpWorklogAppender(jiraSection.Key("base_url").String(),
		jiraSection.Key("username").String(), jiraSection.Key("token").String())

	if err != nil {
		log.Fatal("failed to create worklog appender: %w", err)
	}

	config, intervals, err := timewarrior.Parse(os.Stdin)
	if err != nil {
		log.Fatal("failed to parse timewarrior input: %w", err)
	}

	for _, conf := range config {
		switch conf.Name {
		case "debug":
			if conf.Value == "on" {
				debug = true
			}
		}
	}

	wlAppender.debug = debug
	for _, interval := range intervals {
		issue := ""
		isLogged := false
		comment := ""
		debuglog("inspecting interval %d with tags %s", interval.Id, interval.Tags)
		for _, tag := range interval.Tags {
			prefix, value, found := strings.Cut(tag, ":")
			switch prefix {
			case "iss":
				if found {
					issue = strings.TrimSpace(value)
				} else {
					errorLog("interval %d has invalid issue tag %s", interval.Id, tag)
				}

			case loggedTag:
				isLogged = true

			case "comment":
				comment = value
			}
		}

		if issue != "" && !isLogged && !interval.End.Time().IsZero() {

			delta := int64(interval.End.Time().Sub(interval.Start.Time()).Minutes())
			debuglog("Sending interval %d to issue %s with %d minutes", interval.Id, issue, delta)
			wl := Worklog{
				Started:   JiraDate(interval.Start.Time()),
				TimeSpent: fmt.Sprintf("%dm", delta),
				Comment:   comment,
			}
			if err := wlAppender.Append(issue, wl); err != nil {
				errorLog("failed to send worklog for interval %d to issue %s: %w", interval.Id, issue, err)
			} else {
				timewarrior.Cli.Tag(interval.Id, loggedTag)
			}
		} else {
			debuglog("not sending interval %d", interval.Id)
		}
	}
}

func errorLog(format string, vals ...any) {
	fmt.Printf("[ERROR] "+format+"\n", vals...)
}

func debuglog(format string, vals ...any) {
	if debug {
		fmt.Printf("[DEBUG] "+format+"\n", vals...)
	}
}
