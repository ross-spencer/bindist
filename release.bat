:: release.bat 
:: quick and easy release script for golang releases
::
@ECHO OFF
GOTO START

(we set environment variables for GOARCH and GOOS
    build the project and then rename the files as appropriate)

:START
SET GOOS=linux
SET GOARCH=386
SET FNAME=bindist-linux-i80386-v2.0.1-beta
echo built %FNAME%

go build
rename bindist %FNAME%

SET GOOS=linux
SET GOARCH=amd64
SET FNAME=bindist-linux-amd64-v2.0.1-beta
echo built %FNAME%

go build
rename bindist %FNAME%

SET GOOS=windows
SET GOARCH=386
SET FNAME=bindist-win-i80386-v2.0.1-beta.exe
echo built %FNAME%

go build
rename bindist.exe %FNAME%

SET GOOS=windows
SET GOARCH=amd64
SET FNAME=bindist-win-amd64-v2.0.1-beta.exe
echo built %FNAME%

go build
rename bindist.exe %FNAME%

SET GOOS=
SET GOARCH=

echo release complete