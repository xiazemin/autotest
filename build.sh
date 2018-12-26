#!/bin/bash
go build -o /Users/didi/goLang/bin/autotest /Users/didi/goLang/src/github.com/xiazemin/autotest/autotest/main.go

GOOS=linux GOARCH=amd64 go build -o autotest /Users/didi/goLang/src/github.com/xiazemin/autotest/autotest/main.go
