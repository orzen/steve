.PHONY: all clean dbg-app

all:
	go run main.go --verbose --backend-mongodb --with-client
	./_build/app --mongo-host 0.0.0.0 \
		--mongo-user admin \
		--mongo-password foobar \
		--mongo-database test \
		--log-level debug \
		--pretty

test:
	go test ./...

dbg-steve:
	dlv debug main.go -- --verbose --backend-mongodb --with-client

dbg-srv:
	dlv exec ./_build/app -- --mongo-host 0.0.0.0 \
		--mongo-user admin \
		--mongo-password foobar \
		--mongo-database test \
		--log-level debug \
		--pretty

clean:
	rm -r ./_build
