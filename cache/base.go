package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var LongTermExpire = 86400 * 31 // 最大缓存31天

var conn redis.UniversalClient

// Redis 设置默认redis连接
func Redis(connUrl string, otherAddrs ...string) redis.UniversalClient {
	opts, err := redis.ParseURL(connUrl)
	if err != nil {
		panic(err)
	}
	if len(otherAddrs) == 0 {
		conn = redis.NewClient(opts)
		return conn
	}

	addrs := append([]string{opts.Addr}, otherAddrs...)
	conn = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       addrs,
		Username:    opts.Username,
		Password:    opts.Password,
		DB:          opts.DB,
		PoolSize:    opts.PoolSize,
		Dialer:      opts.Dialer,
		DialTimeout: opts.DialTimeout,
		OnConnect:   opts.OnConnect,
		MaxRetries:  opts.MaxRetries,
	})
	return conn
}

// Client 获取当前redis连接
func Client() redis.UniversalClient {
	if conn == nil {
		Redis("redis://127.0.0.1:6379")
	}
	return conn
}

type RedisBase struct {
	conn    redis.UniversalClient
	ctx     context.Context
	name    string
	timeout int
}

// NewRedisBase 创建redis数据
func NewRedisBase(ctx context.Context, name string, secs int) *RedisBase {
	if secs == 0 {
		secs = LongTermExpire
	}
	return &RedisBase{
		conn: Client(), ctx: ctx, name: name, timeout: secs,
	}
}

// GetName 获取键名
func (r *RedisBase) GetName() string {
	return r.name
}

// Rename 修改名称
func (r *RedisBase) Rename(name string) *RedisBase {
	r.name = name
	return r
}

// SetCtx 更换ctx
func (r *RedisBase) SetCtx(ctx context.Context) *RedisBase {
	r.ctx = ctx
	return r
}

// Ok 返回是否成功
func (r *RedisBase) Ok(ok bool, err error) bool {
	return ok && err == nil
}

// Err 返回错误
func (r *RedisBase) Err(_ string, err error) error {
	return err
}

// Int 返回int结果
func (r *RedisBase) Int(n int64, err error) int {
	if err != nil {
		return -2 // -1 key不存在 -2 操作出错
	}
	return int(n)
}

// ExpireOnce 第一次写入时设置超时时间
func (r *RedisBase) ExpireOnce() (secs int) {
	if secs = r.timeout; secs > 0 {
		r.Expire(secs)
		r.timeout = 0
	}
	return
}

// ExpireKey 设置新的过期时间
func (r *RedisBase) ExpireKey(key string, secs int) bool {
	expire := time.Second * time.Duration(secs)
	op := r.conn.Expire(r.ctx, key, expire)
	return r.Ok(op.Result())
}

// Expire 设置新的过期时间
func (r *RedisBase) Expire(secs int) bool {
	return r.ExpireKey(r.name, secs)
}

// TimeoutKey 有效时间 -1 无限 -2 不存在 -3 出错
func (r *RedisBase) TimeoutKey(key string) int {
	ttl, err := r.conn.TTL(r.ctx, key).Result()
	if err != nil {
		return -3
	} else if ttl <= 0 {
		return int(ttl)
	}
	return int(int64(ttl) / int64(time.Second))
}

// Timeout 有效时间
func (r *RedisBase) Timeout() int {
	return r.TimeoutKey(r.name)
}

// DeleteKey 删除
func (r *RedisBase) DeleteKey(key string) bool {
	_, err := r.conn.Del(r.ctx, key).Result()
	return err == nil
}

// Delete 删除
func (r *RedisBase) Delete() bool {
	return r.DeleteKey(r.name)
}

// ExistsKeys 存在几个key
func (r *RedisBase) ExistsKeys(keys ...string) int {
	op := r.conn.Exists(r.ctx, keys...)
	return r.Int(op.Result())
}

// Keys 查找key
func (r *RedisBase) Keys(pattern string) []string {
	op := r.conn.Keys(r.ctx, pattern)
	if keys, err := op.Result(); err == nil {
		return keys
	}
	return nil
}

// TypeKey 数据类型
func (r *RedisBase) TypeKey(key string) string {
	op := r.conn.Type(r.ctx, key)
	if dt, err := op.Result(); err == nil {
		return dt
	}
	return ""
}

// Type 数据类型
func (r *RedisBase) Type() string {
	return r.TypeKey(r.name)
}

// DataSize 获取数据长度
func (r *RedisBase) DataSize(dk, dt string) int {
	var op *redis.IntCmd
	switch dt {
	default:
		return -2
	case "none":
		return 0
	case "db":
		op = r.conn.DBSize(r.ctx)
	case "string":
		op = r.conn.StrLen(r.ctx, dk)
	case "hash":
		op = r.conn.HLen(r.ctx, dk)
	case "list":
		op = r.conn.LLen(r.ctx, dk)
	case "set":
		op = r.conn.SCard(r.ctx, dk)
	case "zset":
		op = r.conn.ZCard(r.ctx, dk)
	case "stream":
		op = r.conn.XLen(r.ctx, dk)
	}
	return r.Int(op.Result())
}

// Size 获取数据长度
func (r *RedisBase) Size() int {
	return r.DataSize(r.name, r.Type())
}
