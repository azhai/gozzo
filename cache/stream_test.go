package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/azhai/gozzo/cache"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

// go test -run=MQ
func Test11_MQ(t *testing.T) {
	ctx := context.Background()
	cache.Redis("redis://10.10.20.80:6379")
	mq := cache.NewRedisMQ(ctx, "word", "good")

	var ids []string
	num := 100
	for i := 0; i < num; i++ {
		msg := cache.Dict{"seq": i, "msg": time.Now().Format(time.DateTime)}
		ids = append(ids, mq.Send(msg))
		time.Sleep(time.Millisecond * 50)
	}

	mq.Receive(2, func(mq *cache.RedisStream, msg cache.Dict, id, name string) {
		mq.Remove(true, id)
		pp.Println(id, name, ">>", msg)
	})
	time.Sleep(time.Second * 1)
	assert.Len(t, ids, num)
}
