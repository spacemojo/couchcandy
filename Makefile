USERNAME=admin
PASSWORD=rosadiso
URL=http://${USERNAME}:${PASSWORD}@127.0.0.1:5984/cards
JSFORMATTER=~/.jq/jq

start:
	docker-compose -f docker-compose-${VERSION}.yml up -d

stop:
	docker-compose -f docker-compose-${VERSION}.yml down 

init: 
	sleep 5
	curl -X PUT ${URL} | ${JSFORMATTER}
	curl -X POST -H "Content-Type: application/json" ${URL}/_bulk_docs -d @./dev/cards.json | ${JSFORMATTER}

test:
	echo "Launching tests..."

run: start init test stop
