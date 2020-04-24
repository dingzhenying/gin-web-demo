package server

import (
	"errors"
	"fmt"
	"sync"
)

// Event+Conn(事件)
type eventConn struct {
	Event string
	Conn  *Conn
}

// 存储userID 和 eventConn关系的结构体
type binder struct {
	//读写锁（保证读写安全）
	mu sync.RWMutex

	// Map（key:userID,Value:eventConn）
	userID2EventConnMap map[string]*[]eventConn

	// Map（key: connID , value: userID）

	connID2UserIDMap map[string]string
}

// 绑定关系
func (b *binder) Bind(userID, event string, conn *Conn) error {
	if userID == "" {
		return errors.New("userID 不能为空！")
	}

	if event == "" {
		return errors.New("event 不能为空！")
	}

	if conn == nil {
		return errors.New("conn 不能为空！")
	}
	//写锁
	b.mu.Lock()
	defer b.mu.Unlock()

	// map the eConn if it isn't be put.
	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		for i := range *eConns {
			if (*eConns)[i].Conn == conn {
				return nil
			}
		}

		newEConns := append(*eConns, eventConn{event, conn})
		b.userID2EventConnMap[userID] = &newEConns
	} else {
		b.userID2EventConnMap[userID] = &[]eventConn{{event, conn}}
	}
	b.connID2UserIDMap[conn.GetID()] = userID

	return nil
}

//解绑
func (b *binder) Unbind(conn *Conn) error {
	if conn == nil {
		return errors.New("conn can't be empty")
	}
	//写锁
	b.mu.Lock()
	defer b.mu.Unlock()

	// 查询之前是否绑定
	userID, ok := b.connID2UserIDMap[conn.GetID()]
	if !ok {
		return fmt.Errorf("can't find userID by connID: %s", conn.GetID())
	}

	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		for i := range *eConns {
			if (*eConns)[i].Conn == conn {
				newEConns := append((*eConns)[:i], (*eConns)[i+1:]...)
				b.userID2EventConnMap[userID] = &newEConns
				delete(b.connID2UserIDMap, conn.GetID())

				// delete the key of userID when the length of the related
				// eventConn slice is 0.
				if len(newEConns) == 0 {
					delete(b.userID2EventConnMap, userID)
				}

				return nil
			}
		}

		return fmt.Errorf("can't find the conn of ID: %s", conn.GetID())
	}

	return fmt.Errorf("can't find the eventConns by userID: %s", userID)
}

// 获取conn
func (b *binder) FindConn(connID string) (*Conn, bool) {
	if connID == "" {
		return nil, false
	}

	userID, ok := b.connID2UserIDMap[connID]
	// if userID been found by connID, then find the Conn using userID
	if ok {
		if eConns, ok := b.userID2EventConnMap[userID]; ok {
			for i := range *eConns {
				if (*eConns)[i].Conn.GetID() == connID {
					return (*eConns)[i].Conn, true
				}
			}
		}

		return nil, false
	}

	// userID not found, iterate all the conns
	for _, eConns := range b.userID2EventConnMap {
		for i := range *eConns {
			if (*eConns)[i].Conn.GetID() == connID {
				return (*eConns)[i].Conn, true
			}
		}
	}

	return nil, false
}

// 查询userID 和event下的coon
func (b *binder) FilterConn(userID, event string) ([]*Conn, error) {
	if userID == "" {
		return nil, errors.New("userID can't be empty")
	}

	//加读锁
	b.mu.RLock()
	//执行完后解读锁
	defer b.mu.RUnlock()

	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		ecs := make([]*Conn, 0, len(*eConns))
		for i := range *eConns {
			if event == "" || (*eConns)[i].Event == event {
				ecs = append(ecs, (*eConns)[i].Conn)
			}
		}
		return ecs, nil
	}

	return []*Conn{}, nil
}
