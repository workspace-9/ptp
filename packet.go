package ptp

import (
  "bytes"
  "encoding/binary"
  "io"
  "log"
)

const (
  PacketTypeInitReq = 0x1
  PacketTypeInitAck = 0x2
)

type Packet struct {
  PacketType uint32
  Data []byte
  Version [4]byte
}

func (p Packet) Len() uint32 {
  // 4 byte len + 4 byte type + data bytes + 4 byte version
  return 4 + 4 + uint32(len(p.Data) + len(p.Version[:]))
}

func (p Packet) WriteTo(w io.Writer) (n int64, err error) {
  buf := bytes.NewBuffer(make([]byte, 0, p.Len()))
  write := func(b []byte) {
    written, _ := buf.Write(b)
    n += int64(written)
  }

  var header [8]byte
  binary.LittleEndian.PutUint32(header[:4], p.Len())
  binary.LittleEndian.PutUint32(header[4:], p.PacketType)
  write(header[:])
  write(p.Data)
  write(p.Version[:])
  return buf.WriteTo(w)
}

func (p *Packet) ReadFrom(r io.Reader) (n int64, err error) {
  read := func(b []byte) error {
    nRead, err := io.ReadFull(r, b)
    n += int64(nRead)
    return err
  }

  var header [8]byte
  if err = read(header[:]); err != nil {
    return
  }
  log.Println("read header")

  length := binary.LittleEndian.Uint32(header[:4])
  log.Println("total length", length)
  p.PacketType = binary.LittleEndian.Uint32(header[4:])
  log.Println("packet type", p.PacketType)

  dataLen := length - 12
  log.Println("data len", dataLen)
  if dataLen > 0 {
    p.Data = make([]byte, dataLen)
    if err = read(p.Data); err != nil {
      return
    }
  }

  err = read(p.Version[:])
  return
}
