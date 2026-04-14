package ws

func StartBroadcast() {
	for {
		msg := <-Broadcast

		for c := range Clients {
			select {
			case c.Send <- msg:
			default:
				close(c.Send)
				delete(Clients, c)
			}
		}
	}
}
