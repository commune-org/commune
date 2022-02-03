all: clean build 
build: 
	go build -o bin/commune cmd/commune/main.go
	./build.sh
vendor: clean vendorbuild 
vendorbuild:
	go build -mod=vendor -o bin/commune cmd/commune/main.go
clean: 
	rm -f bin/commune
	rm -rf static/css
	mkdir static/css
rebuild: stopUnit clean buildJS build startUnit
production: stopUnit reset clean buildJS build startUnit
cleanBuild: clean build
dev: devReset clean buildJS build modd
modd:
	-modd
refresh: reset clean build run
run:
	./bin/commune
reset: stopSynapse dropDB flushRedis createDB runMigrations startSynapse
devReset: stopDevSynapse dropDB flushRedis createDB runMigrations startDevSynapse
createDB:
	-createdb commune
	-createdb --encoding=UTF8 --locale=C --template=template0 synapse
dropDB:
	-dropdb commune
	-dropdb synapse
stopSynapse: SHELL := /bin/bash
stopSynapse:
	-cd ..;cd synapse;source env/bin/activate;synctl stop;
startSynapse: SHELL := /bin/bash
startSynapse:
	-cd ..;cd synapse;source env/bin/activate;synctl start;
runMigrations:
	-cd db/migrations;goose postgres "postgres://commune:@localhost:5432/commune?sslmode=disable" up;
flushRedis: SHELL := /bin/bash
flushRedis:
	-REDISCLI_AUTH=$$REDISAUTH redis-cli -n 1 flushdb
buildJS:
	-cd ui/js;npm run production;
startUnit:
	-systemctl --user start commune-dev.service
stopUnit:
	-systemctl --user stop commune-dev.service
setup:
	-go get -d github.com/cortesi/modd/cmd/modd;
	-go install github.com/pressly/goose/v3/cmd/goose@latest;
	-cd ..;git clone https://github.com/matrix-org/synapse;cd synapse;cp ../commune/docs/alt* demo/;python3 -m venv ./env;source ./env/bin/activate;pip install -e ".[all,dev]";
stopDevSynapse: SHELL := /bin/bash
stopDevSynapse:
	-cd ..;cd synapse;source env/bin/activate;./demo/stop.sh;
startDevSynapse: SHELL := /bin/bash
startDevSynapse:
	-cd ..;cd synapse;source env/bin/activate;./demo/alt-start.sh;
cleanDevSynapse: SHELL := /bin/bash
cleanDevSynapse:
	-cd ..;cd synapse;source env/bin/activate;./demo/alt-clean.sh;
