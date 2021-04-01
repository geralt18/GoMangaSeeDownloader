$now = Get-Date -UFormat "%Y-%m-%d_%T"
$ver = "0.1.0"
$filePath = "bin/mangasee-$ver.exe" 
$zipPath = "bin/mangasee-$ver.zip" 

go build -o $filePath -ldflags "-X main.Version=$ver -X main.BuildTime=$now" main.go

Compress-Archive -Path $filePath -DestinationPath $zipPath