# egomap

Egomap is a lock-free, eventually consistent, concurrent map.

It is my naive implementation of Rust's [`evmap`](https://github.com/jonhoo/evmap) in go, for educational purposes. I also don't intend to implement the exact same algorithm and API as `evmap`. There is another implementation of [go-evmap](https://github.com/clarkmcc/go-evmap), but I believe some implementation details are not correct.

## Usage

```go
readHandle, writer := egomap.New[string, int]()
reader := readHandle.Reader()

// insert element into map
writer.Set("alice", 25)

// attempt to read
age, ok := reader.Get("alice") // 0, false

// refresh the writer to apply changes
writer.Refresh()
age, ok := reader.Get("alice") // 25, true

reader.Close()
```

### alt usage
```go
studentScores := egomap.NewHandle[string, int](1) // refresh after every 1 operation

// readers still need to be instanstiated
reader := studentScores.Reader()

// insert element
studentScores.Set("alice", 78)

// attempt to read
reader.Get("alice") // 78, true

// writer handle can be refreshed manually
studentScores.Refresh()

```

## Benchmarks

The benchmark is set up as such:
- Initialize array of size 1M
- Generate a random array of keys [0, 1M) over a zipfian distribution
    - For each key, pair with a randomly selected read or write operation (99%, 99.9%, 99.99%, 99.999% reads)
- Initialize `GOMAXPROCS` goroutines (using `-cpu 1, 2, 4, 8`) with random start indices

### Results


raw output in [`bench.out`](https://github.com/dingyuchen/egomap/blob/master/bench.out)

Test Specs:
- Macbook M1 Pro (10 cores)
- 16GB RAM
- `go v1.19.9`

## Features

- zero dependencies
- 100% `go`

## Limitations

As with `evmap` and other implementations, this library assumes a singular writer. A synchronized handle for multiple writers is provided.

This map is backed by 2 hashmaps which implies there are 2 copies of data, possibly leading high memory usage. Please pass pointers into the map if you want to avoid duplicate data.

## Chores

- [ ] Add tests
- [ ] Avoid malloc for oplog and queue
- [ ] Come up with better way to register and deregister readers
- [x] Add benchmarks with `sync.Map`
- [x] Profile and optimize
