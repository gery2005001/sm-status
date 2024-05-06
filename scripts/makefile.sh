#! /bin/bash
###
 # @Author: gery2005
 # @Date: 2022-08-07 02:06:33
 # @LastEditors: gery2005
 # @LastEditTime: 2022-09-26 18:38:36
 # @FilePath: \scripts\makefile.sh
 # @Description: 编译项目脚本
 # 
 # Copyright (c) 2022 by gery2005@gmail.com, All Rights Reserved. 
### 

current_date=`date -d "1 minute ago" +"%Y-%m-%d"`
current_time=`date -d "1 minute ago" +"%H:%M:%S%:z"`
go_version=`go version`
app_version="0.4.4"

go build  -o sm-status -ldflags "-s -w -X sm-status/version.Version=$app_version -X 'sm-status/version.BuildDate=$current_date' -X 'sm-status/version.BuildTime=$current_time' -X 'sm-status/version.GO_Version=$go_version'" .