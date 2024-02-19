package ptp

// EncodeString encodes a string into the wire format expected by ptp.
// Each byte is separated by a NULL char, and the end is denoted by two NULL chars.
func EncodeString(s string) []byte {
  ret := make([]byte, len(s)*2 + 2)
  for idx, b := range []byte(s) {
    ret[idx<<1] = b
    ret[(idx<<1)+1] = 0
  }
  ret[len(ret)-2] = 0
  ret[len(ret)-1] = 0
  return ret
}
