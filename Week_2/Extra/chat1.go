import (
	"io"
	"net"
)

// Our first proverb–don’t mediate access to shared memory with locks and mutexes,
// instead share that memory by communicating. So let’s apply this advice to our chat server.

type Mux struct {
	add     chan net.Conn
	remove  chan net.Addr
	sendMsg chan string
}

func (m *Mux) Add(conn net.Conn) {
	m.add <- conn
}

// Add sends the connection to add to the add channel.

func (m *Mux) Remove(addr net.Addr) {
	m.remove <- addr
}

// Remove sends the address of the connection to the remove channel.

func (m *Mux) SendMsg(msg string) error {
	m.sendMsg <- msg
	return nil
}

// And send message sends the message to be transmitted to each connection to the sendMsg channel.

func (m *Mux) loop() {
	conns := make(map[net.Addr]net.Conn)
	for {
		select {
		case conn := <-m.add:
			m.conns[conn.RemoteAddr()] = conn
		case addr := <-m.remove:
			delete(m.conns, addr)
		case msg := <-m.sendMsg:
			for _, conn := range m.conns {
				io.WriteString(conn, msg)
			}
		}
	}
}

// Rather than using a mutex to serialise access to the conns map, loop will wait until it receives an operation in the form of a value sent over one of the add,
// remove, or sendMsg channels and apply the relevant case. We don’t need a mutex anymore because the shared state, our conns map, is local to the loop function.

// But, there’s still a lot of hard coded logic here. loop only knows how to do three things; add, remove and broadcast a message.
// As with the previous example, adding new features to our Mux type will involve:

//     creating a channel.
//     adding a helper to send the data over the channel.
//     extending the select logic inside loop to process that data.

// Just like our Calculator example we can rewrite our Mux to use first class functions to pass around behaviour we want to executed, not data to interpret.
// Now, each method sends an operation to be executed in the context of the loop function, using our single ops channel.