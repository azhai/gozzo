package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// FlatDict 将字典转为一维数组
func FlatDict(data Dict) (items []any) {
	for k, v := range data {
		items = append(items, k, v)
	}
	return
}

/*************************************/
/*************  字符串  ***************/
/*************************************/

// RedisString redis字符串
type RedisString struct {
	*RedisBase
}

// NewRedisString 创建redis字符串
func NewRedisString(ctx context.Context, name string, secs int) *RedisString {
	return &RedisString{
		RedisBase: NewRedisBase(ctx, name, secs),
	}
}

// Type 获取类型
func (r *RedisString) Type() string {
	return "string"
}

// Size 获取长度
func (r *RedisString) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Incr 自增操作
func (r *RedisString) Incr(v int) int {
	op := r.conn.IncrBy(r.ctx, r.name, int64(v))
	r.ExpireOnce()
	return r.Int(op.Result())
}

// Set 设置值和有效期
func (r *RedisString) Set(v any) bool {
	secs := LongTermExpire
	if r.timeout > 0 {
		secs = r.timeout
	}
	dur := time.Second * time.Duration(secs)
	ok, err := r.conn.Set(r.ctx, r.name, v, dur).Result()
	return ok == "ok" && err == nil
}

// AddLock 加排他锁
func (r *RedisString) AddLock(secs int) bool {
	if secs <= 0 {
		secs = LongTermExpire
	}
	dur := time.Second * time.Duration(secs)
	now := time.Now().Format(time.DateTime)
	op := r.conn.SetNX(r.ctx, r.name, now, dur)
	return r.Ok(op.Result())
}

/*************************************/
/*************  哈希表  ***************/
/*************************************/

// RedisHash redis哈希表
type RedisHash struct {
	*RedisBase
	Prefix string
}

// NewRedisHash 创建redis哈希表
func NewRedisHash(ctx context.Context, name, prefix string, secs int) *RedisHash {
	return &RedisHash{
		RedisBase: NewRedisBase(ctx, prefix+name, secs),
		Prefix:    prefix,
	}
}

// Type 获取类型
func (r *RedisHash) Type() string {
	return "hash"
}

// Size 获取长度
func (r *RedisHash) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Rename 修改名称
func (r *RedisHash) Rename(name string) *RedisHash {
	r.name = r.Prefix + name
	return r
}

// Incr 自增操作
func (r *RedisHash) Incr(m string, v int) int {
	op := r.conn.HIncrBy(r.ctx, r.name, m, int64(v))
	r.ExpireOnce()
	return r.Int(op.Result())
}

// GetAll 获取全部
func (r *RedisHash) GetAll() map[string]string {
	data, _ := r.conn.HGetAll(r.ctx, r.name).Result()
	return data
}

// Get 获取单个字段
func (r *RedisHash) Get(m string) string {
	v, _ := r.conn.HGet(r.ctx, r.name, m).Result()
	return v
}

// Set 设置单个字段
func (r *RedisHash) Set(data ...any) int {
	op := r.conn.HSet(r.ctx, r.name, data...)
	r.ExpireOnce()
	return r.Int(op.Result())
}

// Merge 设置多个数据
func (r *RedisHash) Merge(data Dict) bool {
	op := r.conn.HMSet(r.ctx, r.name, FlatDict(data)...)
	r.ExpireOnce()
	return r.Ok(op.Result())
}

/*************************************/
/*************   队列   ***************/
/*************************************/

// RedisList redis队列
type RedisList struct {
	*RedisBase
}

// NewRedisList 创建redis队列
func NewRedisList(ctx context.Context, name string, secs int) *RedisList {
	return &RedisList{
		RedisBase: NewRedisBase(ctx, name, secs),
	}
}

// Type 获取类型
func (r *RedisList) Type() string {
	return "list"
}

// Size 获取长度
func (r *RedisList) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Push 数据入栈
func (r *RedisList) Push(data ...any) int {
	if len(data) == 0 {
		return r.Size()
	}
	op := r.conn.LPush(r.ctx, r.name, data...)
	r.ExpireOnce()
	return r.Int(op.Result())
}

// Pop 单个数据出栈
func (r *RedisList) Pop(secs int, others ...string) (name, value string) {
	if secs > 0 { // 阻塞版本
		dur := time.Second * time.Duration(secs)
		others = append(others, r.name)
		data, _ := r.conn.BRPop(r.ctx, dur, others...).Result()
		if len(data) == 2 {
			name, value = data[0], data[1]
		}
	} else { // 非阻塞版本
		v, err := r.conn.RPop(r.ctx, r.name).Result()
		if err == nil && v != "" {
			name, value = r.name, v
		}
	}
	return
}

// PopN 数据出栈
func (r *RedisList) PopN(n int) (data []string) {
	data, _ = r.conn.RPopCount(r.ctx, r.name, n).Result()
	return
}

/*************************************/
/************* 无序集合 ***************/
/*************************************/

// RedisSet redis无序集合
type RedisSet struct {
	*RedisBase
}

// NewRedisSet 创建redis无序集合
func NewRedisSet(ctx context.Context, name string, secs int) *RedisSet {
	return &RedisSet{
		RedisBase: NewRedisBase(ctx, name, secs),
	}
}

// Type 获取类型
func (r *RedisSet) Type() string {
	return "set"
}

// Size 获取长度
func (r *RedisSet) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Add 增加元素
func (r *RedisSet) Add(m string) int {
	op := r.conn.SAdd(r.ctx, r.name, m)
	r.ExpireOnce()
	return r.Int(op.Result())
}

// Drop 删除元素
func (r *RedisSet) Drop(m any) int {
	n, _ := r.conn.SRem(r.ctx, r.name, m).Result()
	return int(n)
}

// Move 在集合间移动元素
func (r *RedisSet) Move(dst, m string) bool {
	op := r.conn.SMove(r.ctx, r.name, dst, m)
	return r.Ok(op.Result())
}

// Rand 随机一个元素
func (r *RedisSet) Rand(isPop bool) string {
	var op *redis.StringCmd
	if isPop {
		op = r.conn.SPop(r.ctx, r.name)
	} else {
		op = r.conn.SRandMember(r.ctx, r.name)
	}
	if m, err := op.Result(); err == nil {
		return m
	}
	return ""
}

// RandN 随机多个元素
func (r *RedisSet) RandN(n int, isPop bool) []string {
	var op *redis.StringSliceCmd
	if isPop {
		op = r.conn.SPopN(r.ctx, r.name, int64(n))
	} else {
		op = r.conn.SRandMemberN(r.ctx, r.name, int64(n))
	}
	if ms, err := op.Result(); err == nil {
		return ms
	}
	return nil
}

// IsMember 是否其中一个元素
func (r *RedisSet) IsMember(m string) bool {
	op := r.conn.SIsMember(r.ctx, r.name, m)
	return r.Ok(op.Result())
}

// Members 返回所有元素
func (r *RedisSet) Members() []string {
	op := r.conn.SMembers(r.ctx, r.name)
	if ms, err := op.Result(); err == nil {
		return ms
	}
	return nil
}

/*************************************/
/*************  有序集合 ***************/
/*************************************/

// RedisZSet redis有序集合
type RedisZSet struct {
	*RedisBase
	Prefix string
}

// NewRedisZSet 创建redis有序集合
func NewRedisZSet(ctx context.Context, name string, secs int) *RedisZSet {
	return &RedisZSet{
		RedisBase: NewRedisBase(ctx, name, secs),
	}
}

// Type 获取类型
func (r *RedisZSet) Type() string {
	return "zset"
}

// Size 获取长度
func (r *RedisZSet) Size() int {
	return r.DataSize(r.name, r.Type())
}

// Incr 自增分数
func (r *RedisZSet) Incr(m string, s float64) float64 {
	s, _ = r.conn.ZIncrBy(r.ctx, r.name, s, m).Result()
	r.ExpireOnce()
	return s
}

// Add 增加元素
func (r *RedisZSet) Add(m string, s float64) int {
	args := redis.Z{Member: m, Score: s}
	op := r.conn.ZAdd(r.ctx, r.name, args)
	r.ExpireOnce()
	return r.Int(op.Result())
}

// Score 获取分数
func (r *RedisZSet) Score(m string) float64 {
	s, _ := r.conn.ZScore(r.ctx, r.name, m).Result()
	return s
}

// GetRange 按分数读取元素
func (r *RedisZSet) GetRange(min, max string) []string {
	args := &redis.ZRangeBy{Min: min, Max: max}
	op := r.conn.ZRangeByScore(r.ctx, r.name, args)
	if ms, err := op.Result(); err == nil {
		return ms
	}
	return nil
}

// GetRangeScores 读取元素和分数
func (r *RedisZSet) GetRangeScores(min, max string) []redis.Z {
	args := &redis.ZRangeBy{Min: min, Max: max}
	op := r.conn.ZRangeByScoreWithScores(r.ctx, r.name, args)
	if zs, err := op.Result(); err == nil {
		return zs
	}
	return nil
}

// DropRange 按分数删除元素
func (r *RedisZSet) DropRange(min, max string) int {
	op := r.conn.ZRemRangeByScore(r.ctx, r.name, min, max)
	return r.Int(op.Result())
}

// DropRangeInt 按分数删除元素
func (r *RedisZSet) DropRangeInt(start, stop int) int {
	min, max := strconv.Itoa(start), strconv.Itoa(stop)
	return r.DropRange(min, max)
}
