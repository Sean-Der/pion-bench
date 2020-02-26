#!/bin/sh
go build && rm -rf /tmp/profile* && ./pion-bench && go tool pprof --pdf ./pion-bench /tmp/profile*/cpu.pprof > file.pdf
