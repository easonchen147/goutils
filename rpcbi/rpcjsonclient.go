package rpcbi

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"encoding/json"

	"errors"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
)

type RpcJsonClient struct {
	rpcServer *rpc.Server

	ClientId int32
	*rpc.Client
}

func NewRpcJsonClient() *RpcJsonClient {
	return &RpcJsonClient{
		rpcServer: rpc.NewServer(),
	}
}

func (c *RpcJsonClient) doServer(sess *yamux.Session) {
	clientConn, err := sess.Accept()
	if err != nil {
		log.Panic(err)
		return
	}
	c.rpcServer.ServeCodec(jsonrpc.NewServerCodec(clientConn))
}

func (c *RpcJsonClient) doClient(sess *yamux.Session) {
	clientConn, err := sess.Open()
	if err != nil {
		log.Panic(err)
		return
	}
	c.Client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(clientConn))
}

func (c *RpcJsonClient) Start(conn net.Conn) error {
	err := c.handshake(conn)
	if err != nil {
		conn.Close()
		return err
	}

	sess, err := yamux.Client(conn, nil)
	if err != nil {
		log.Panic(err)
	}

	safego.Go(func() {
		c.doServer(sess)
	})
	c.doClient(sess)
	return nil
}

func (c *RpcJsonClient) handshake(conn net.Conn) error {
	bs, _ := json.Marshal(&HandshakeReq{Version: PROTOCOL_VERSION, Id: "11"})
	handshake := &utils.ControlPacket{ProtocolId: PROTOCOL_HANDSHAKE, Data: bs}
	handshake.Pack(conn)

	handshakeAck := &utils.ControlPacket{}
	handshakeAck.Unpack(conn)

	if handshakeAck.ProtocolId != PROTOCOL_HANDSHAKE_ACK {
		return errors.New("protocol not match")
	}

	resp := &HandshakeResp{}
	err := json.Unmarshal(handshakeAck.Data, resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}

	return nil
}

func (c *RpcJsonClient) RegisterService(name string, service interface{}) error {
	return c.rpcServer.RegisterName(name, service)
}

func (c *RpcJsonClient) Close() {
	c.Client.Close()
}