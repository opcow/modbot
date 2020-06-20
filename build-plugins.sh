#!/bin/bash
cd plugins
go build -buildmode=plugin ./cdbot.go
go build -buildmode=plugin ./covidbot.go
go build -buildmode=plugin ./reactbot.go
go build -buildmode=plugin ./pongbot.go
cd ..
