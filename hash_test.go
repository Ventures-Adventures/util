package hash

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var (
	sum      = 10000000
	sTime    time.Time
	duration int64
)

func Test_t(t *testing.T) {
	var mp = make(map[uint32]int)
	sTime = time.Now()
	for i := 0; i < sum; i++ {
		v := StrToUint32(strconv.Itoa(i)+"DSLSDFASLFDELGDSLGSDASFADSFDSFEWsdfsdfwezsd", 512)
		mp[v]++
	}
	duration = time.Now().Sub(sTime).Nanoseconds()
	fmt.Printf("StrToUint32 Benchmark %.2f ns/op\n", float64(duration)/float64(sum))
	fmt.Printf("%d: %+v\n", len(mp), mp)
}
