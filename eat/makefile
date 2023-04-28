run:
	go build 
	go run ./petsAdopt

build:
	go build

develop:
	git pull
# 	go build
	pm2 delete "petsAdopt --c config.develop"
	pm2 start "./petsAdopt --c config.develop.toml"

meixi:
	git pull
# 	go build
	pm2 delete "petsAdopt --c config.meixi"
	pm2 start "./petsAdopt --c config.meixi.toml"

production:
	git pull
# 	go build
	pm2 delete "petsAdopt --c config.production"
	pm2 start "./petsAdopt --c config.production.toml"
	
win2linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build

mac2linux:
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	go build

macToLinux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

logs:
	pm2 logs

log:
	pm2 log

status:
	pm2 status
