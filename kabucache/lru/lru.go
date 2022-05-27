package lru

import (
	"container/list"
)

//包含字典和双向链表的结构体cache
type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List                    //存放头指针
	cache     map[string]*list.Element      //map
	OnEvicted func(key string, value Value) //回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

//实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查找对应的双向链表的节点并移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//淘汰算法
func (c *Cache) RemoveOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.ll.Remove(ele) //删除最后一个节点
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key) //删除映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 增加一个value或者进行更新
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest() //超出容量就刷新
	}
}

//测试内存
func (c *Cache) Len() int {
	return c.ll.Len()
}
