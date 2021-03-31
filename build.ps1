$now = Get-Date -UFormat "%Y-%m-%d_%T"
$ver = "0.1.0"
$filePath = "bin/mangasee-$ver.exe" 

go build -o $filePath -ldflags "-X main.Version=$ver -X main.BuildTime=$now" .\src\main.go