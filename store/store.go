package store

import (
	"context"
	"time"
)

// Storer 令牌存储接口
type Storer interface {
	// 存储令牌数据，并指定到期时间
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	// 获取存储数据
	Get(ctx context.Context, key string) (string, bool, error)
	// 检查令牌是否存在
	Check(ctx context.Context, key string) (bool, error)
	// 存放一个键值, 只用来确定是否存在
	Set1(ctx context.Context, key string, expiration time.Duration) error
	// 删除存储的令牌
	Delete(ctx context.Context, key string) error
	// 关闭存储
	Close() error
}
