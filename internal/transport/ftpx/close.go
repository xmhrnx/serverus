package ftpx

// Close signs out.
func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
		c.mu.Lock()
		_ = c.sc.Quit()
		c.mu.Unlock()
	})
	return nil
}
