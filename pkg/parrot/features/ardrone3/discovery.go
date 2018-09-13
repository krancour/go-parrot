package ardrone3

import (
	"fmt"
	"net"
)

func (b *BebopClient) discover() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, discoveryPort))
	if err != nil {
		return err
	}
	b.discoveryClient, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	b.discoveryClient.Write(
		[]byte(
			fmt.Sprintf(
				`{
					"controller_type": "computer",
					"controller_name": "go-bebop",
					"d2c_port": "%d",
					"arstream2_client_stream_port": "%d",
					"arstream2_client_control_port": "%d",
				}`,
				d2cPort,
				rtpStreamPort,
				rtpControlPort),
		),
	)
	data := make([]byte, 10240)
	_, err = b.discoveryClient.Read(data)
	if err != nil {
		return err
	}
	return b.discoveryClient.Close()
}
