package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/shadowsocks"
	"v2ray.com/core/testing/assert"
)

func TestUDPEncoding(t *testing.T) {
	assert := assert.On(t)

	request := &protocol.RequestHeader{
		Version: Version,
		Command: protocol.RequestCommandUDP,
		Address: v2net.LocalHostIP,
		Port:    1234,
		User: &protocol.User{
			Email: "love@v2ray.com",
			Account: loader.NewTypedSettings(&Account{
				Password:   "shadowsocks-password",
				CipherType: CipherType_AES_128_CFB,
				Ota:        Account_Disabled,
			}),
		},
	}

	data := alloc.NewLocalBuffer(256).Clear().AppendString("test string")
	encodedData, err := EncodeUDPPacket(request, data)
	assert.Error(err).IsNil()

	decodedRequest, decodedData, err := DecodeUDPPacket(request.User, encodedData)
	assert.Error(err).IsNil()
	assert.Bytes(decodedData.Value).Equals(data.Value)
	assert.Address(decodedRequest.Address).Equals(request.Address)
	assert.Port(decodedRequest.Port).Equals(request.Port)
}

func TestTCPRequest(t *testing.T) {
	assert := assert.On(t)

	request := &protocol.RequestHeader{
		Version: Version,
		Command: protocol.RequestCommandTCP,
		Address: v2net.LocalHostIP,
		Option:  RequestOptionOneTimeAuth,
		Port:    1234,
		User: &protocol.User{
			Email: "love@v2ray.com",
			Account: loader.NewTypedSettings(&Account{
				Password:   "tcp-password",
				CipherType: CipherType_CHACHA20,
			}),
		},
	}

	data := alloc.NewLocalBuffer(256).Clear().AppendString("test string")
	cache := alloc.NewLargeBuffer().Clear()

	writer, err := WriteTCPRequest(request, cache)
	assert.Error(err).IsNil()

	writer.Write(data)

	decodedRequest, reader, err := ReadTCPSession(request.User, cache)
	assert.Error(err).IsNil()
	assert.Address(decodedRequest.Address).Equals(request.Address)
	assert.Port(decodedRequest.Port).Equals(request.Port)

	decodedData, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bytes(decodedData.Value).Equals([]byte("test string"))
}
