package lru

import (
	"container/list"
)

//最近最少使用算法
//创建一个包含字典和双向链表的结构体类型 Cache，方便实现后续的增删查改操作

// Cache is a LRU cache. It is not safe for concurrent access.
// 这是一个最近最少使用算法 并发访问是不安全的
type Cache struct {
	maxBytes int64                    //最大内存
	nbytes   int64                    //已使用内存
	ll       *list.List               //双向链表
	cache    map[string]*list.Element //值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	//可选，并在清除 entry 时执行
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 使用长度来统计他带上了多少字节
// 返回值所占内存大小
type Value interface {
	Len() int
}

// New 方便实例化 Cache，实现 New() 函数：构造函数
// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找功能
// 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾。
// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok { //如果key存在于cache中
		c.ll.MoveToFront(ele) //将c.cache[key]对应的结点ele 移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除 实际上是缓存淘汰。即移除最近最少访问的节点（队首）
// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //Back returns the last element of list ll or nil if the list is empty.
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                //从字典中 c.cache 删除该节点的映射关系。
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //c所占内存为
		if c.OnEvicted != nil {                                //如果回调函数 OnEvicted 不为 nil，则调用回调函数
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
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
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
