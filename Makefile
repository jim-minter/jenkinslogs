jenkinslogs: $(shell find cmd pkg vendor)
	go build ./cmd/jenkinslogs

clean:
	@rm -f jenkinslogs

.PHONY: clean
