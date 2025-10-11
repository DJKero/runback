package bot

import (
	"github.com/servusdei2018/shards/v2"
)

type Bot struct {
	ShardsMgr *shards.Manager
}

var Client Bot
