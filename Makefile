test: export dbname=test
test: export dbhost=http://localhost
test: export dbusername=admin
test: export dbpassword=nimda

get:
	go get -t -v ./...

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic