//go:build !linux && !windows

package zsyscall

func (c *SockSOOption) Controllers() []Controller {
	return nil
}

func (c *SockIPOption) Controllers() []Controller {
	return nil
}

func (c *SockIPV6Option) Controllers() []Controller {
	return nil
}

func (c *SockTCPOption) Controllers() []Controller {
	return nil
}

func (c *SockUDPOption) Controllers() []Controller {
	return nil
}
