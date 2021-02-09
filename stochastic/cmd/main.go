package main

import (
	"github.com/Sinu5oid/generators/stochastic"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	stochastic.BuildModel()
}
