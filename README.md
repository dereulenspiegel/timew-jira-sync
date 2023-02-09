# timew-jira-sync

Simple tool to sync timewarrior reports to Jira Tempo.

## Example config.ini

```
[jira]
base_url = https://your.jira.example.com
username = <username>
token = <apitoken>
```

## Installation

This binary needs to be installed into your Timewarrior extension folder. Either
you copy this binary there or you create a symlink to `~/.timewarrior/extensions/jira-sync`.
See the [extension docs](https://timewarrior.net/docs/api/) of Timewarrior for more details.

## Usage

To upload intervals as worklogs to Jira Tempo you need to add specifically
formatted tags to these intervals. At minimum you need to add a tag like
`iss:<Ticket-ID>` to the interval. Alternatively you can add a tag cotnaining
a description with `description:<Your Description>`, but this is not necessary.

Uploading to Jira Temp can be done with `timew report jira-sync :yesterday` for example.
See https://timewarrior.net/reference/timew-report.1/ for more details.
