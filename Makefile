all:
	go build .

watch:
	EBDEMO_URL=localhost EBDEMO_USER=gajduler EBDEMO_DB=ebdemo EBDEMO_PASSWORD=omedbe fresh .

docker: ebdemo
	docker build . -t=ebdemo && docker run --rm -p 1234:3000 -e EBDEMO_URL=localhost -e EBDEMO_USER=gajduler -e EBDEMO_DB=ebdemo -e EBDEMO_PASSWORD=omedbe ebdemo

ebdemo:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ebdemo .

clean:
	rm ebdemo
