package app

import (
	"sync"
)

type ConnectionInfo interface {
	comparable
}

type Connection interface {
	bool
}

type ConnectionPool[K ConnectionInfo, V Connection] struct {
	connections map[K]V
	connMux     sync.RWMutex
}

func (cp *ConnectionPool[K, V]) NewConnection(connInfo K) {
	cp.connMux.Lock()
	if cp.connections == nil {
		cp.connections = make(map[K]V)
	}
	cp.connections[connInfo] = true
	cp.connMux.Unlock()
}

func (cp *ConnectionPool[K, V]) CloseConnection(connInfo K) {
	cp.connMux.Lock()
	delete(cp.connections, connInfo)
	cp.connMux.Unlock()
}

func (cp *ConnectionPool[K, V]) GetConnections() (connections []K) {
	cp.connMux.RLock()
	for conn := range cp.connections {
		connections = append(connections, conn)
	}
	cp.connMux.RUnlock()
	return
}
