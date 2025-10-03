package models

import (
	"github.com/servusdei2018/shards/v2"
)

type Client struct {
	ShardsMgr *shards.Manager

	Token  string
	Owners []int
}

var Bot Client
