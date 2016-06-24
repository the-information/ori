package account

import (
	"encoding/base64"
	"math/big"
)

const fnv_128_offset_basis = "144066263297769815596495629667062367629"
const fnv_128_prime = "309485009821345068724781371"

// fnv1a128 returns the FNV-1a 128-bit variant hash of in as
// a URL-safe unpadded base64 string as defined in RFC 4648.
func fnv1a128(in []byte) string {

	h := big.NewInt(0)
	h.SetString(fnv_128_offset_basis, 10)
	prime := big.NewInt(0)
	prime.SetString(fnv_128_prime, 10)
	lower128 := big.NewInt(0)
	lower128.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
	for _, inByte := range in {
		h.Xor(h, big.NewInt(int64(inByte)))
		h.Mul(h, prime)
		h.And(h, lower128)
	}

	return base64.RawURLEncoding.EncodeToString(h.Bytes())

}
