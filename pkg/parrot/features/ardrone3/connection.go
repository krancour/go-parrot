package ardrone3

import (
	"fmt"
	"net"
	"time"
)

func (b *BebopClient) Connect() error {
	err := b.discover()

	if err != nil {
		return err
	}

	c2daddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, c2dPort))

	if err != nil {
		return err
	}

	b.c2dClient, err = net.DialUDP("udp", nil, c2daddr)

	if err != nil {
		return err
	}

	d2caddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", d2cPort))

	if err != nil {
		return err
	}
	b.d2cClient, err = net.ListenUDP("udp", d2caddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err := b.c2dClient.Write(<-b.writeChan)

			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	go func() {
		for {
			data := make([]byte, 40960)
			i, _, err := b.d2cClient.ReadFromUDP(data)
			if err != nil {
				fmt.Println("d2cClient error:", err)
			}

			b.packetReceiver(data[0:i])
		}
	}()

	// send pcmd values at 40hz
	go func() {
		// wait a little bit so that there is enough time to get some ACKs
		time.Sleep(500 * time.Millisecond)
		for {
			_, err := b.write(b.generatePcmd().Bytes())
			if err != nil {
				fmt.Println("pcmd c2dClient.Write", err)
			}
			time.Sleep(25 * time.Millisecond)
		}
	}()

	if err := b.GenerateAllStates(); err != nil {
		return err
	}
	if err := b.FlatTrim(); err != nil {
		return err
	}

	return nil
}
