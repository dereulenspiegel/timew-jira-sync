# timew-jira-sync

Simple tool to sync timewarrior reports to Jira Tempo.

## Example config.ini

```
[jira]
base_url = https://your.jira.example.com
username = <username>
token = <apitoken>
```

## Usage

To upload intervals as worklogs to Jira Tempo you need to add specifically
formatted tags to these intervals. At minimum you need to add tag like
`iss:<Ticket-ID>` to the interval. Alternatively you can add a tag cotnaining
a description with `description:<Your Description>`, but this is not necessary
