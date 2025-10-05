package kvcache

import (
	"sync"
)

type Node struct {
	key, value string
	prev, next *Node
}

type LRUCache struct {
	capacity   int
	store      map[string]*Node
	head, tail *Node
	mtx        sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		store:    make(map[string]*Node),
	}
}

func (c *LRUCache) addToFront(node *Node) {
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache) removeNode(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

func (c *LRUCache) moveToFront(node *Node) {
	c.removeNode(node)
	c.addToFront(node)
}

func (c *LRUCache) removeTail() *Node {
	if c.tail == nil {
		return nil
	}
	removed := c.tail
	c.removeNode(removed)
	return removed
}

func (c *LRUCache) Get(key string) (string, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	node, exists := c.store[key]
	if !exists {
		return "", false
	}
	c.moveToFront(node)
	return node.value, true
}

func (c *LRUCache) Set(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if node, exists := c.store[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}
	newNode := &Node{key: key, value: value}
	c.store[key] = newNode
	c.addToFront(newNode)

	if len(c.store) > c.capacity {
		removed := c.removeTail()
		delete(c.store, removed.key)
	}
}
