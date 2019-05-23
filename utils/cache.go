package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
)

var cc cache.Cache

// InitCache 
func InitCache() {
	cacheConfig := beego.AppConfig.String("cache::cache")
	cc = nil

	if "redis" == cacheConfig {
		initRedis()
	} else {
		initMemcache()
	}


}

func initMemcache() {
	var err error
	cc, err = cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("cache::memcache_host")+`"}`)

	if err != nil {
		beego.Info(err)
	}

}

func initRedis() {
	LogOut("info", "The cache is redis")
	// cc = &cache.Cache{}
	var err error

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("initial redis error caught: %v\n", r)
			cc = nil
		}
	}()
	host := beego.AppConfig.String("cache::redis_host")
	LogOut("info", "Connection parameter :"+host)
	cc, err = cache.NewCache("redis", `{"conn":"`+host+`"}`)

	if err != nil {
		LogOut("error", err)
	}
}

// SetCache
func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			LogOut("error", r)
			cc = nil
		}
	}()
	timeouts := time.Duration(timeout) * time.Second
	err = cc.Put(key, data, timeouts)
	if err != nil {
		LogOut("error", err)
		LogOut("error", "SetCache failed，key:"+key)

		return err
	} else {

		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			LogOut("error", r)
			cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		return errors.New("Cache does not exist")
	}

	err := Decode(data.([]byte), to)
	if err != nil {
		LogOut("error", err)
		LogOut("error", "GetCache failed，key:"+key)
	}

	return err
}

// DelCache
func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	err := cc.Delete(key)
	if err != nil {
		return errors.New("Cache failed to delete")
	} else {
		return nil
	}
}

// Encode
// (Data encoding with gob)
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode
// (Data decoding with gob)
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}