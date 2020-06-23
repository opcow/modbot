#!/bin/bash
cd plugins
go build -buildmode=plugin ./cdbot/cdbot.go
go build -buildmode=plugin ./covidbot/covidbot.go
go build -buildmode=plugin ./reactbot/reactbot.go
go build -buildmode=plugin ./pongbot/pongbot.go
go build -buildmode=plugin ./saybot/saybot.go
cd ..
