package hash

func StrToUint32(str string, mod uint32) (v uint32) {
	v = 2166136261
	for _, c := range []byte(str) {
		v ^= uint32(c)
		v *= 16777619
	}
	return v % mod
}
