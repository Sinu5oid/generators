# generators package

A package duplicating standard 'rand' package

## Content:
 
* immediate package `generators`: various generators, benchmarks
    * `cmd/single_dimensional` contains demo usage of single-dimensional distribution generators
    * `cmd/two_dimensional` contains demo usage of two-dimensional distribution generators
* package `stat`: statistics analysis package (contains Pearson test support for single-component distributions)
    * `cmd` contains demo usage of Pearson test function and utilities

## Usage in go

* clone package directly into your `GOPATH` or with `go get -u github.com/Sinu5oid/generators/..`
* use in your code

## Standalone versions

* cd to chosen `/cmd` directory
* run `go build .`
* run `cmd` executable