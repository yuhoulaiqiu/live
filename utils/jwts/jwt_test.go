package jwts

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(JwtPayload{
		UserID:   1,
		Nickname: "test",
		Role:     1,
	}, "123456", 3600)
	fmt.Println(token, err)
}

func TestParseToken(t *testing.T) {
	payload, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJuYW1lIjoidGVzdCIsInJvbGUiOjEsImV4cCI6MTcxNDQ5NTAxNX0.dXNBhXG1Yt0KJt6ampnaEVxYWX9dhvwl0t_8rGDlqXo", "123456")
	fmt.Println(payload, err)
}
