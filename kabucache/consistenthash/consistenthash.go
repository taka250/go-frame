package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//byte to uint
type Hash func(data []byte) uint32

//map 有所有的key
type Map struct {
	Hash     Hash
	replicas int            //虚ls拟节点倍数
	keys     []int          //stored 哈希环
	hashMap  map[int]string //虚拟节点映射表
}

//new map
func New(replicas int, fn Hash) *Map {
	m := &Map{
		Hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.Hash == nil {
		m.Hash = crc32.ChecksumIEEE
	}
	return m

}

//添加真实节点的Add方法
func (n *Map) Add(keys ...string) {
	for _, key := range keys { //每一个真实节点对应多个虚拟节点
		for i := 0; i < n.replicas; i++ {
			hash := int(n.Hash([]byte(strconv.Itoa(i) + key))) //转换以后加入到
			n.keys = append(n.keys, hash)
			n.hashMap[hash] = key
		}
	}
	sort.Ints(n.keys)
}

//选择节点的方法
func (n *Map) Get(key string) string {
	if len(n.keys) == 0 {
		return ""
	}

	hash := int(n.Hash([]byte(key)))
	idx := sort.Search(len(n.keys), func(i int) bool { //二分法顺时针查找第一个匹配的下标
		return n.keys[i] >= hash
	})
	return n.hashMap[n.keys[idx%len(n.keys)]] //如果求余为0说明是最大的要选择第一个节点也就是n.keys[0],同时通过映射得到真实的节点
}
