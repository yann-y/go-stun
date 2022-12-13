package stun

import (
	"fmt"
	"net"
	"time"
)

func (c *Client) natTTL(conn net.PacketConn, addr *net.UDPAddr) error {
	var mappedAddr *Host
	var t time.Duration
	for i := 250; i >= 200; i = i - 10 {
		c.logger.Debugln("Do natTTL time", i)
		c.logger.Debugln("Send To:", addr)
		resp, err := c.test1(conn, addr)
		if err != nil {
			return err
		}
		c.logger.Debugln("Received:", resp)
		if resp == nil {
			return err
		}
		if mappedAddr == nil {
			mappedAddr = resp.mappedAddr
		} else {
			c.logger.Debugln("mappedAddr post ", mappedAddr.port, " resp post ", resp.mappedAddr.Port())
			if mappedAddr.port == resp.mappedAddr.Port() {
				c.logger.Info(mappedAddr.port, "   ", resp.mappedAddr.Port())
				goto LOOP
			}
			mappedAddr = resp.mappedAddr
		}
		t = time.Duration(i)
		time.Sleep(time.Second * t)
	}
LOOP:
	fmt.Println("nat ttl time: ", t*time.Second)
	return nil
}
