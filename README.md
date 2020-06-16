# modbot

| Command  | Description  | Req. Op  |
|---|---|---|
| !op \<userid\> \| \<@user\> | add one or more users to the operators  | yes  |
| !deop \<userid\> \| \<@user\> | remove one or more users from the operators  | yes  |
| !delmsg \<server id\> \<message id\> | delete a message  | no  |
| !config | print the current config via direct message | yes  |
| !quit  | kill the bot  | yes  |


---
`go build -buildmode=plugin covidbot.go`
