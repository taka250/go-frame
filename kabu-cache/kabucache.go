package kabucache

import (
	"fmt"
	"sync"
)

//定义接口
type Getter interface {
	Get(key string) ([]byte, error)
}

//定义函数类型并实现接口的get方法
type GetterFunc func(key string) ([]byte, error)

//实现Get方法  这里是接口型函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//最核心的group
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//实例化group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter ")
	}
	mu.Lock()
	defer mu.Unlock()
	g := Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = &g
	return &g
}

//用名字查找一个group
func GetGroup(name string) *Group {
	mu.RLocker()
	g := groups[name]
	mu.RUnlock()
	return g
}

//核心方法get实现了查找缓存，如果存在就返回不存在就调用load
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getlocally(key)
}

func (g *Group) getlocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
