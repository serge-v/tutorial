all: wethome debug-server ocean blog ocean.zip

wethome: wethome.go
	go build wethome.go

debug-server: debug-server.go
	go build debug-server.go

ocean: ocean.go ocean-params.go
	go build ocean.go ocean-params.go

ocean.exe: ocean.go ocean-params.go
	GOOS=windows GOARCH=amd64 go build ocean.go ocean-params.go

ocean.zip: ocean.exe
	zip ocean.zip ocean.exe

blog: blog.go
	go build blog.go

clean:
	rm -f wethome debug-server ocean blog

deploy_blog: ocean
	./deploy.sh

deploy_wethome: ocean
	GOOS=freebsd GOARCH=amd64 go build 
