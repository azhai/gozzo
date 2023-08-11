package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Dict 字典类型
type Dict = map[string]any

// HandlerFunc 处理消息的方法
type HandlerFunc func(mq *RedisStream, msg Dict, id, name string)

// NewRedisMQ 简易消息队列
func NewRedisMQ(ctx context.Context, name, group string) *RedisStream {
	mq := &RedisStream{RedisBase: NewRedisBase(ctx, name, -1)}
	if group != "" {
		_ = mq.CreateGroup(group)
	}
	return mq
}

// RedisStream redis消息队列
type RedisStream struct {
	*RedisBase
	customerGroup string
}

// Type 获取类型
func (r *RedisStream) Type() string {
	return "stream"
}

// Size 获取长度
func (r *RedisStream) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Receive 接收消息
func (r *RedisStream) Receive(workers int, handler HandlerFunc) {
	for i := 0; i < workers; i++ {
		customerName := fmt.Sprintf("customer-%04d", i)
		go func(name string) {
			for {
				topic, msgs := r.ReadMessages(name, 1)
				if topic != "" && len(msgs) == 1 {
					handler(r, msgs[0].Values, msgs[0].ID, name)
				}
			}
		}(customerName)
	}
}

// Send 发送多条消息
func (r *RedisStream) Send(msgs ...Dict) (msgid string) {
	for _, msg := range msgs {
		msgid = r.Publish(msg)
	}
	return
}

// SendPairs 发送单条消息，参数必须偶数个
func (r *RedisStream) SendPairs(pairs ...any) string {
	return r.Publish(pairs)
}

// CreateGroup 创建消费组
// 已有同名消费组时会报错：BUSYGROUP Consumer CreateGroup name already exists
func (r *RedisStream) CreateGroup(group string) error {
	r.customerGroup = group
	op := r.conn.XGroupCreate(r.ctx, r.name, r.customerGroup, "0")
	return r.Err(op.Result())
}

// DestroyGroup 删除消费组
func (r *RedisStream) DestroyGroup() int {
	if r.customerGroup == "" {
		return 0
	}
	op := r.conn.XGroupDestroy(r.ctx, r.name, r.customerGroup)
	return r.Int(op.Result())
}

// Ack 确认消息，防止消息被重复读取
func (r *RedisStream) Ack(ids ...string) int {
	op := r.conn.XAck(r.ctx, r.name, r.customerGroup, ids...)
	return r.Int(op.Result())
}

// Remove 删除消息
func (r *RedisStream) Remove(needAck bool, ids ...string) int {
	if needAck && len(ids) != r.Ack(ids...) {
		return 0
	}
	op := r.conn.XDel(r.ctx, r.name, ids...)
	return r.Int(op.Result())
}

// Trim 保留最新的一些消息，使用XTrimApprox更高效
func (r *RedisStream) Trim(size int, isApprox bool) int {
	var op *redis.IntCmd
	if isApprox {
		op = r.conn.XTrimMaxLenApprox(r.ctx, r.name, int64(size), 0)
	} else {
		op = r.conn.XTrimMaxLen(r.ctx, r.name, int64(size))
	}
	return r.Int(op.Result())
}

// Publish 发布消息
func (r *RedisStream) Publish(data any) string {
	args := &redis.XAddArgs{Stream: r.name, ID: "*"}
	switch data.(type) {
	default:
		panic("the data type is not correct.")
	case []string, []any, Dict:
		args.Values = data
	}
	id, err := r.conn.XAdd(r.ctx, args).Result()
	if err != nil {
		fmt.Println("RedisStream Publish error:", err)
		return ""
	}
	return id
}

// Subscribe 接收消息
func (r *RedisStream) Subscribe(consumer string, ack bool, count, secs int) []redis.XStream {
	block := time.Second * time.Duration(secs)
	args := &redis.XReadGroupArgs{
		Consumer: consumer, Group: r.customerGroup, Streams: []string{r.name, ">"},
		Block: block, Count: int64(count), NoAck: ack == false,
	}
	streams, err := r.conn.XReadGroup(r.ctx, args).Result()
	if err == nil {
		return streams
	} else if strings.HasPrefix(err.Error(), "NOGROUP") {
		r.DestroyGroup()
		return nil
	} else {
		fmt.Println("RedisStream Subscribe error:", err)
		return nil
	}
}

// ReadMessages 读取消息
func (r *RedisStream) ReadMessages(consumer string, count int) (string, []redis.XMessage) {
	streams := r.Subscribe(consumer, true, count, 1)
	if len(streams) == 1 {
		return streams[0].Stream, streams[0].Messages
	}
	return "", nil
}

// MoveMessages 转移消息
func (r *RedisStream) MoveMessages(src, dst string, count, secs int) int {
	idle := time.Second * time.Duration(secs)
	pendingArgs := &redis.XPendingExtArgs{
		Consumer: src, Group: r.customerGroup, Stream: r.name,
		Idle: idle, Count: int64(count), Start: "-", End: "+",
	}
	pendLst, err := r.conn.XPendingExt(r.ctx, pendingArgs).Result()
	if err != nil || len(pendLst) == 0 {
		fmt.Println("RedisStream MoveMessages XPendingExt error:", err)
		return 0
	}

	var ids []string
	for _, pend := range pendLst {
		ids = append(ids, pend.ID)
	}
	claimArgs := &redis.XClaimArgs{
		Consumer: dst, Group: r.customerGroup, Stream: r.name,
		MinIdle: idle, Messages: ids,
	}
	msgLst, err := r.conn.XClaim(r.ctx, claimArgs).Result()
	if err != nil {
		fmt.Println("RedisStream MoveMessages XClaim error:", err)
		return 0
	}
	return len(msgLst)
}
