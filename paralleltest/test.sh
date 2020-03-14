#!/bin/sh

# デフォルトの並列度はコア数
echo "go test -v ."
go test -v .

# busy loopの場合は遅くなるが、
# time.Sleepだとゴールーチンが切り替わるのでそれなりに早い
echo "go test -v -parallel=1 ."
go test -v -parallel=1 .

# busy loopでも速くなる（1コアマシンを除く）
echo "go test -v -parallel=2 ."
go test -v -parallel=2 .
