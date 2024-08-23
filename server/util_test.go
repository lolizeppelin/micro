package server

import (
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {

	fmt.Printf("sn 1 : %s\n", SNBase62(1))
	fmt.Printf("sn 443 : %s\n", SNBase62(443))
	fmt.Printf("sn 242424 : %s\n", SNBase62(24242))
	fmt.Printf("sn 100000 : %s\n", SNBase62(100000))
	fmt.Printf("sn 1000000 : %s\n", SNBase62(1000000))
	fmt.Printf("sn 2000000 : %s\n", SNBase62(2000000))

	fmt.Printf("now 12 se: %0*s\n", 12, utils.ToBase62(int(time.Now().Unix())))
	fmt.Printf("now 12 ms: %0*s\n", 12, utils.ToBase62(int(time.Now().UnixMicro())))
	fmt.Printf("now 15 ms: %0*s\n", 15, utils.ToBase62(int(time.Now().UnixMicro())))
	fmt.Printf("now 16 ms: %0*s\n", 16, utils.ToBase62(int(time.Now().UnixMicro())))
}
