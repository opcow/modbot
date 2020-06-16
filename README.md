# modbot

| Command  | Description  | Req. Op  |
|---|---|---|
| !op \<userid\> \| \<@user\> | add one or more users to the operators  | yes  |
| !deop \<userid\> \| \<@user\> | remove one or more users from the operators  | yes  |
| !quit  | kill the bot  | yes  |

Building a plugin:

`go build -buildmode=plugin pongbot.go`

---
