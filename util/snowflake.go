package util

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"sync"
)

var node *snowflake.Node
var once sync.Once

func NextID() uint64 {

	return uint64(getNode().Generate().Int64())
}

func getNode() *snowflake.Node {

	once.Do(func() {

		var err error
		node, err = snowflake.NewNode(rand.Int63() % 1024)
		if err != nil {

			panic(err)
		}
	})

	return node
}
