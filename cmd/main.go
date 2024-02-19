package main

import (
  "log"
  "time"

  "github.com/google/uuid"
  "github.com/workspace-9/ptp"
)

func main() {
  guid, err := uuid.FromBytes([]byte{
    0xc7, 0x15, 0x7c, 0x74, 0xa0, 0xa7, 0x06, 0x63, 0x05, 0x93, 0xbb, 0x05, 0x76, 0xcb, 0xd8, 0x2f,
  })
  if err != nil {
    panic(err)
  }

  conn, err := ptp.Dial("tcp", "192.168.1.2:15740", guid, "jordan.framework.arch")
  if err != nil {
    log.Fatalf(err.Error())
  }

  time.Sleep(time.Second * 5)
  log.Println(conn.Close())
}
