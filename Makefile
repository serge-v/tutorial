all: wethome debug-server ocean blog readme.txt

readme.txt: *.go
	go doc > readme.txt

wethome: wethome.go
	go build wethome.go

debug-server: debug-server.go
	go build debug-server.go

ocean: ocean.go ocean-params.go
	go build ocean.go ocean-params.go

blog: blog.go
	go build blog.go

clean:
	rm -f wethome debug-server ocean blog

deploy_blog: ocean
	ocean -deploy blog
	./deploy.sh

deploy_wethome: ocean
	GOOS=freebsd GOARCH=amd64 go build wethome.go
	./ocean -deploy wethome
	rm wethome
