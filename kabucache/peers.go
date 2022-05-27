package kabucache

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(sey string) (peer PeerGetter, ok bool)
}

//getter则是要被节点实现的
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
