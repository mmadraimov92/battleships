package battleships

import "io"

func (b *Battleships) setConnection(conn io.ReadWriter) {
	b.conn = conn
}
