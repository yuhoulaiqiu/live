package stream

import (
	"fmt"
	"testing"
)

func TestGetChannelKey(t *testing.T) {
	c := GetChannelKey("0000001")
	fmt.Println("c=", c)
}
