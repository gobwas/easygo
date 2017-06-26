// +build linux

package netpoll

// New creates new epoll-based Poller instance with given config.
func New(c *Config) (Poller, error) {
	epoll, err := EpollCreate(c)
	if err != nil {
		return nil, err
	}

	return Epoller{epoll}, nil
}

// Epoller implements Poller interface.
type Epoller struct {
	*Epoll
}

// Start implements Poller.Start() method.
func (ep Epoller) Start(desc *Desc, cb CallbackFn) error {
	return ep.Add(desc.fd(), modeToEvent(desc.mode),
		func(events EpollEvent) {
			var mode Mode
			if events&(EPOLLIN|EPOLLRDHUP|EPOLLHUP|EPOLLERR) != 0 {
				mode |= ModeRead
			}
			if events&(EPOLLOUT|EPOLLHUP|EPOLLERR) != 0 {
				mode |= ModeWrite
			}
			if events&(EPOLLCLOSED) != 0 {
				mode |= ModeClosed
			}

			cb(mode)
		},
	)
}

// Stop implements Poller.Stop() method.
func (ep Epoller) Stop(desc *Desc) error {
	return ep.Del(desc.fd())
}

// Resume implements Poller.Resume() method.
func (ep Epoller) Resume(desc *Desc) error {
	return ep.Mod(desc.fd(), modeToEvent(desc.mode))
}

func modeToEvent(mode Mode) (events EpollEvent) {
	if mode&ModeRead != 0 {
		events |= EPOLLIN | EPOLLRDHUP
	}
	if mode&ModeWrite != 0 {
		events |= EPOLLOUT
	}
	if mode&ModeOneShot != 0 {
		events |= EPOLLONESHOT
	}
	if mode&ModeEdgeTriggered != 0 {
		events |= EPOLLET
	}
	return events
}
