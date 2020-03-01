#!/bin/sh
go build -o pion-bench "$1" && rm -rf /tmp/profile* && ./pion-bench && go tool pprof --pdf ./pion-bench /tmp/profile*/cpu.pprof > results.pdf
