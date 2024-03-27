@ECHO OFF
for /f "tokens=*" %%d in ('date /t') do (set BUILD_DATE=%%d)
for /f "tokens=*" %%t in ('time /t') do (set BUILD_TIME=%%t)
for /f "tokens=*" %%g in ('go version') do (set GO_VER=%%g)
set VERSION=0.2.1


CMD /C go build -ldflags "-s -w -X sm-status/version.Version=%VERSION% -X 'sm-status/version.BuildDate=%BUILD_DATE%' -X 'sm-status/version.BuildTime=%BUILD_TIME%' -X 'sm-status/version.GO_Version=%GO_VER%'" -o sm-status.exe