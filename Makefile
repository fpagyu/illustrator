

.PHONY: ai2ps
ai2ps:
	go build -ldflags "-w -s" -o bin/ai2ps cmd/ai2ps/ai2ps.go
