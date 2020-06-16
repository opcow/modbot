#!/bin/bash
cd plugins
go build -v -buildmode=plugin ./cdbot.go
go build -v -buildmode=plugin ./covidbot.go
go build -v -buildmode=plugin ./reactbot.go
go build -v -buildmode=plugin ./pongbot.go
cd ..
