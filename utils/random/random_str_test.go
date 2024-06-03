package random

import (
	"fmt"
	"testing"
)

func TestRandomStr(t *testing.T) {
	fmt.Println(RandomStr(10))
	fmt.Println(RandomStr(10))
	fmt.Println(RandomStr(10))
	fmt.Println(RandomStr(10))
}
