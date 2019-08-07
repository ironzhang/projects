# README.md

生成性能剖析文件

```
go test -bench=. -cpuprofile cpu.profile
```

分析性能剖析文件

```
go tool pprof spam.test cpu.profile
```

v1版本性能测试

```
$ go test -bench=.
PASS
BenchmarkTest1-2         1000000              1837 ns/op
BenchmarkTest2-2            1000           1875367 ns/op
ok      github.com/ironzhang/pearls/spam        3.928s
```
