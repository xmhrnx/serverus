package monitor

// Err is the last sampling error, nil when the last sample worked.
func (c *Collector) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}
