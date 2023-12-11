

.PHONY: ai2ps
ai2ps:
	go build -ldflags "-w -s" -o ai2ps cmd/ai2ps/ai2ps.go
