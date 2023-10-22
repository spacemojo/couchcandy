test: export dbname=test
test: export dbhost=http://localhost
test: export dbusername=admin
test: export dbpassword=nimda

get:
	go get -t -v ./...

# See go help testflag for details in the currently used flags
test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

showcoverage:
	go tool cover -func=coverage.out

htmlcoverage:
	go tool cover -html=coverage.out
