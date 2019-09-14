# slackwatch
Run some commands or hit some URLs when DMed on Slack

[![GoDoc](https://godoc.org/github.com/mikegrb/slackwatch/pkg/slackwatch?status.svg)](https://godoc.org/github.com/mikegrb/slackwatch/pkg/slackwatch)

Create a .slackwatch config file in your home directory. See config.json for an example. If a URL is specified, the presence of a Body determines an HTTP GET or POST.

```
   {
     "SlackToken": "xoxp-123-543",
     "WatchedChannels": [
       "SuperImportant"
     ],
     "Actions": [
       { "Command": "/usr/bin/afplay", "Args": "klaxon.wav" },
       {
         "URL": "https://hassio.local/api/services/homeassistant/turn_on?api_password=letmein",
         "Body": "{\"entity_id\":\"switch.bat_signal\"}",
         "ContentType": "application/javascript"
       }
     ]
   }
```

Alternatively, you can create your own actions to preform that conform to the Action interface and pass your config to the slackwatch.New constructor.

klaxon.wav is a public domain recording provided by the US Navy.
