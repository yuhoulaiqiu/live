package pwd

import (
	"fmt"
	"testing"
)

func TestHashPwd(t *testing.T) {
	hash := HashPwd("123456")
	fmt.Println(hash)
}
func TestCheckPwd(t *testing.T) {
	hash := HashPwd("123456")
	fmt.Println(CheckPwd(hash, "123456"))
}
