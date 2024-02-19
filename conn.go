package ptp

import (
  "fmt"
  "log"
  "net"

  "github.com/google/uuid"
)

type Conn struct {
  conn net.Conn
  version [4]byte
  guid uuid.UUID
  name string
}

func Dial(network string, addr string, guid uuid.UUID, name string) (ptpConn Conn, err error) {
  conn, err := net.Dial(network, addr)
  if err != nil {
    return
  }
  log.Printf("Connected to %s", addr)

  ptpConn = Conn{conn, [4]byte{0, 0, 1, 0}, guid, name}
  err = ptpConn.Initiailize()
  return
}

func (c Conn) MakePacket(typ uint32, data ...[]byte) Packet {
  totalLen := 0
  for _, part := range data {
    totalLen += len(part)
  }

  packData := make([]byte, totalLen)
  start := 0
  for _, part := range data {
    copy(packData[start:], part)
    start += len(part)
  }

  return Packet{
    typ, packData, c.version,
  }
}

func (c Conn) Initiailize() error {
  log.Printf("Initializing connection")
  _, err := c.MakePacket(PacketTypeInitReq, c.guid[:], EncodeString(c.name)).WriteTo(c.conn)
  if err != nil {
    return err
  }
  log.Println("Sent init req")

  resp := &Packet{}
  _, err = resp.ReadFrom(c.conn)
  if err != nil {
    return err
  }
  log.Println("Received response")

  if resp.PacketType != PacketTypeInitAck {
    return fmt.Errorf("Expected init ack, got %d", resp.PacketType)
  }
  log.Println("Response was valid")

  return nil
}

func (c Conn) Close() error {
  return c.conn.Close()
}
