type Mux struct {
	ops chan func(map[net.Addr]net.Conn)
}

func (m *Mux) Add(conn net.Conn) {
	m.ops <- func(m map[net.Addr]net.Conn) {
			m[conn.RemoteAddr()] = conn
	}
}

// In this case the signature of the operation is a function which takes a map of net.Addr’s to net.Conn’s. 
// In a real program you’d probably have a much more complicated type to represent a client connection, but it’s sufficient for the purpose of this example.

func (m *Mux) Remove(addr net.Addr) {
	m.ops <- func(m map[net.Addr]net.Conn) {
			delete(m, addr)
	}
}

// Remove is similar, we send a function that deletes its connection’s address from the supplied map.

func (m *Mux) SendMsg(msg string) error {
	m.ops <- func(m map[net.Addr]net.Conn) {
			for _, conn := range m {
					io.WriteString(conn, msg)
			}
	}
	return nil
}

// SendMsg is a function which iterates over all connections in the supplied map and calls io.WriteString to send each a copy of the message.

func (m *Mux) loop() { 
	conns := make(map[net.Addr]net.Conn)
	for op := range m.ops {
			op(conns)
	}
}

// You can see that we’ve moved the logic from the body of loop into anonymous functions created by our helpers. 
// So the job of loop is now to create a conns map, wait for an operation to be provided on the ops channel, then invoke it, passing in its map of connections.

// But there are a few problems still to fix. The most pressing is the lack of error handling in SendMsg; an error
//  writing to a connection will not be communicated back to the caller. So let’s fix that now.

func (m *Mux) SendMsg(msg string) error {
	result := make(chan error, 1)
	m.ops <- func(m map[net.Addr]net.Conn) {
			for _, conn := range m.conns {
					err := io.WriteString(conn, msg)
					if err != nil {
							result <- err
							return
					}
			}
			result <- nil
	}
	return <-result
}

// To handle the error being generated inside the anonymous function we pass to loop we need to create a channel to 
// communicate the result of the operation. This also creates a point of synchronisation, the last line of SendMsg blocks until the function we passed into loop has been executed.

func (m *Mux) loop() {
	conns := make(map[net.Addr]net.Conn)
	for op := range m.ops {
			op(conns)
	}
}

// Note that we didn’t have the change the body of loop at all to incorporate this error handling. And now we know how to do this, 
// we can easily add a new function to Mux to send a private message to a single client.

func (m *Mux) PrivateMsg(addr net.Addr, msg string) error {
	result := make(chan net.Conn, 1)
	m.ops <- func(m map[net.Addr]net.Conn) {
			result <- m[addr]
	}
	conn := <-result
	if conn == nil {
			return errors.Errorf("client %v not registered", addr)
	}
	return io.WriteString(conn, msg)
}

// To do this we pass a “lookup function” to loop via the ops channel, which will look in the map provided to it—this is loop‘s conns map—and 
// return the value for the address we want on the result channel.

// In the rest of the function we check to see if the result was nil—the zero value from the map lookup implies that the client is not 
// registered. Otherwise we now have a reference to the client and we can call io.WriteString to send them a message.

// And just to reiterate, we did this all without changing the body of loop, or affecting any of the other operations.