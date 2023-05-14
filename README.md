# egomap

Egomap is a lock-free, eventually consistent, concurrent map.

It is my naive implementation of Rust's [`evmap`](https://github.com/jonhoo/evmap) in go, for educational purposes. I also don't intend to implement the exact same algorithm and API as `evmap`. There is another implementation of [go-evmap](https://github.com/clarkmcc/go-evmap), but I believe some implementation details are not correct.


