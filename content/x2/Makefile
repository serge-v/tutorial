helloworld:
	go build helloworld.go
	./helloworld -local

worldtime:
	go build worldtime.go
	./worldtime -local

prod:
	@echo building for FreeBSD...
	GOOS=freebsd GOARCH=amd64 go build helloworld.go
	GOOS=freebsd GOARCH=amd64 go build worldtime.go
	@echo deploying to ocean...
	ocean -deploy helloworld
	ocean -deploy worldtime
	rm worldtime

clean:
	rm -f helloworld worldtime
