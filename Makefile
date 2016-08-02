all:
	go build .

watch:
	EBDEMO_URL=localhost EBDEMO_USER=gajduler EBDEMO_DB=ebdemo EBDEMO_PASSWORD=omedbe fresh .
