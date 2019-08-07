package base58

import "math/big"

var (
	encodeStd = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
)

var StdEncoding = NewEncoding(encodeStd)

var (
	bigInt0  = big.NewInt(0)
	bigInt58 = big.NewInt(58)
)

type Encoding struct {
	encode [58]byte
	decode [256]*big.Int
}

func NewEncoding(encoder []byte) *Encoding {
	enc := Encoding{}
	copy(enc.encode[:], encoder)
	for i, c := range encoder {
		enc.decode[c] = big.NewInt(int64(i))
	}
	return &enc
}

func (enc *Encoding) Encode(src []byte) (dst []byte) {
	dst = make([]byte, len(src)*138/100+1)

	var m, x big.Int
	var di = len(dst) - 1

	x.SetBytes(src)
	for x.Cmp(bigInt0) != 0 {
		x.DivMod(&x, bigInt58, &m)
		dst[di] = enc.encode[m.Int64()]
		di--
	}
	for _, b := range src {
		if b != 0 {
			break
		}
		dst[di] = enc.encode[0]
		di--
	}
	return dst[di+1:]
}

func (enc *Encoding) Decode(src []byte) (dst []byte) {
	dst = make([]byte, 0, len(src)*733/1000+1)

	si := 0
	for _, b := range src {
		if b != enc.encode[0] {
			break
		}
		dst = append(dst, 0)
		si++
	}
	x := big.NewInt(0)
	for ; si < len(src); si++ {
		m := enc.decode[src[si]]
		if m == nil {
			break
		}
		x.Mul(x, bigInt58)
		x.Add(x, m)
	}
	return append(dst, x.Bytes()...)
}
