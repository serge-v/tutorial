tutorial project contains tutorials for my wife to start with Go language.
The goal is to teach her to

    - create simple web apps
    - deploy apps to the cloud web server
    - generate html using templates
    - process html forms
    - interact with mariadb
    - create browser notifications in javascript

blog program itself is a CGI handler for golang.org/x/tools/blog package.
Locally it runs under debug-server. In production it runs under Apache web
server.

debug-server runs blog locally on http://localhost:8081. Start it and edit
articles.

ocean-params.go is a file with sshParams structure containing credentials
for connecting to the ocean. This file is not stored in git.

ocean is a tool to deploy golang application to digitalocean cloud server.
It builds the app for freebsd-amd64 and copies to /usr/local/www/wet virtual
directory. Uses ocean-params.go file which contains credentials to connect
to the ocean.

wethome is a home page application for wet.voilokov.com.
