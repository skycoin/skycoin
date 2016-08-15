# ChaCha20

[![Build Status](https://travis-ci.org/tang0th/go-chacha20.svg?branch=master)](https://travis-ci.org/tang0th/go-chacha20)

This is a go implementation of the chacha20 cipher by [djb](http://cr.yp.to). It uses the C implementation of 
the chacha20 cipher by [Ted Krovetz](http://krovetz.net/csus/). If cgo cannot be used it falls back to a go
implementation of chacha20.

## Speed

Ted Krovetz's version of chacha20 appears to be far faster than AES on my PC (encrypting 1Kb of data).

    BenchmarkChaCha2020_1K   1000000              2438 ns/op         419.86 MB/s
    BenchmarkChaCha2012_1K   1000000              1898 ns/op         539.48 MB/s
    BenchmarkAES256CTR_1K     500000              6564 ns/op         155.99 MB/s
    BenchmarkAES192CTR_1K     500000              6226 ns/op         164.46 MB/s
    BenchmarkAES128CTR_1K     500000              5938 ns/op         172.44 MB/s