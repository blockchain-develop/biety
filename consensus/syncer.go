package consensus

type Syncer struct {
	server  *Server
}

func NewSyncer(server *Server) *Syncer {
	sync := &Syncer {
		server : server,
	}

	return sync
}

func (self *Syncer) run() {

}
