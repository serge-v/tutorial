GOOS=freebsd GOARCH=amd64 go build blog.go
host=$(cat host~)
rsync -rvaz content template blog ${host}:
ssh ${host} sudo cp -r content /usr/local/www/wet/
ssh ${host} sudo cp -r template /usr/local/www/wet/
ssh ${host} sudo cp blog /usr/local/www/wet/
