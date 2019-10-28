import (
	"io"
	"net"
	"sync"
)

//NONIDIOMATIC
type Mux struct {
	mu    sync.Mutex
	conns map[net.Addr]net.Conn
}

func (m *Mux) Add(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conns[conn.RemoteAddr()] = conn
}

func (m *Mux) Remove(addr net.Addr) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, addr)
}

func (m *Mux) SendMsg(msg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, conn := range m.conns {
		err := io.WriteString(conn, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Because this is a server, all of these methods will be called concurrently,
//  so we need to use a mutex to protect the conns map and prevent data races. Is this what you’d call idiomatic Go code?