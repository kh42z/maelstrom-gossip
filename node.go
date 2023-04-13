package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/rs/xid"
	"sync"
)

type node struct {
	Node     *maelstrom.Node
	idsMutex sync.RWMutex
	ids      []int
}

func NewNode() *node {
	n := &node{
		Node: maelstrom.NewNode(),
		ids:  make([]int, 0),
	}
	n.Node.Handle("echo", n.echoHandler)
	n.Node.Handle("generate", n.generateHandler)
	n.Node.Handle("broadcast", n.broadcastHandler)
	n.Node.Handle("read", n.readHandler)
	n.Node.Handle("topology", n.topologyHandler)
	return n
}

func (n *node) echoHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	body["type"] = "echo_ok"
	return n.Node.Reply(msg, body)
}

func (n *node) generateHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	body["type"] = "generate_ok"
	body["id"] = xid.New().String()
	return n.Node.Reply(msg, body)
}

func (n *node) broadcastHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	n.idsMutex.Lock()
	defer n.idsMutex.Unlock()
	body["type"] = "broadcast_ok"
	id := body["message"].(float64)
	n.ids = append(n.ids, int(id))
	delete(body, "message")
	return n.Node.Reply(msg, body)
}

func (n *node) readHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	n.idsMutex.RLock()
	defer n.idsMutex.RUnlock()
	body["type"] = "read_ok"
	body["messages"] = n.ids
	return n.Node.Reply(msg, body)
}

func (n *node) topologyHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	body["type"] = "topology_ok"
	delete(body, "topology")
	return n.Node.Reply(msg, body)
}
