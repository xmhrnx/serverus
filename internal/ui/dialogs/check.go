package dialogs

import (
	"context"
	"time"
)

func (d *ServerDialog) testConnection() {
	s, err := d.collect()
	if err != nil {
		d.err = err.Error()
		return
	}
	d.err, d.testMsg = "", ""
	d.test.Loading(true)
	dialer := d.dialers(formCredentials{base: d.creds, password: d.password.Text(), passphrase: d.passphrase.Text()})
	start := time.Now()
	var elapsed time.Duration
	d.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		conn, err := dialer.Dial(ctx, s)
		if err != nil {
			return err
		}
		elapsed = time.Since(start)
		return conn.Close()
	}, func(err error) {
		d.test.Loading(false)
		if err != nil {
			d.testOK, d.testMsg = false, err.Error()
			return
		}
		d.testOK, d.testMsg = true, "Connected in "+elapsed.Round(time.Millisecond).String()
	})
}
