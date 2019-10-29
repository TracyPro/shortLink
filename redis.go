package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base62"
	"github.com/go-redis/redis/v7"
	"io"
	"time"
)

// redis key
const (
	URLIDKEY               = "next.url.id"         // 全局计数器 自增ID
	ShortLinkKey           = "shortlink:%s:url"    // 映射短地址和长地址之间的关系（通过key拿到短地址对应的长地址）
	URLHashKey             = "urlhash:%s:url"      // 映射长地址和短地址之间的关系（通过key拿到短地址对应的长地址）
	ShortLinkDetailInfoKey = "shortlink:%s:detail" // 映射短地址和它的详细信息
)

// redis client
type RedisClient struct {
	Client *redis.Client
}

type URLDetail struct {
	URL                 string        `json:"url"`
	Created             string        `json:"created"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisClient{Client: client}
}

// 长地址转换为短地址
func (r *RedisClient) Shorten(url string, exp int64) (string, error) {
	// url转为sha1 hash
	h := toSha1(url)
	// 检查url是否在缓冲
	d, err := r.Client.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// 不存在
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// 过期
		} else {
			// 存在
			return url, nil
		}
	}
	// 自增
	err = r.Client.Incr(URLIDKEY).Err()
	if err != nil {
		return "", err
	}
	// 将key转为base64
	id, err := r.Client.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}

	eid := base62.EncodeInt64(id)
	err = r.Client.Set(fmt.Sprintf(ShortLinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	err = r.Client.Set(fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(&URLDetail{
		URL:                 url,
		Created:             time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})

	if err != nil {
		return "", err
	}

	// 存储到redis
	err = r.Client.Set(fmt.Sprintf(ShortLinkDetailInfoKey, eid), detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	// 返回短地址
	return eid, nil
}

func (r *RedisClient) ShortLinkInfo(eid string) (interface{}, error) {
	d, err := r.Client.Get(fmt.Sprintf(ShortLinkDetailInfoKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{
			404,
			errors.New("Unknown short URL"),
		}
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}

}

func (r *RedisClient) Unshort(eid string) (string, error) {
	url, err := r.Client.Get(fmt.Sprintf(ShortLinkKey, eid)).Result()
	// 未找到
	if err == redis.Nil {
		return "", StatusError{
			404,
			errors.New("Unknown short URL"),
		}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func toSha1(url string) string {
	t := sha1.New()
	io.WriteString(t, url)
	return fmt.Sprintf("%x", t.Sum(nil));
}
