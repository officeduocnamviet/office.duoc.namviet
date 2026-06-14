package main

import (
	"fmt"
	"log"
	"github.com/namviet/backend-core/internal/features/users"
	"github.com/namviet/backend-core/internal/platform/supabase"
)

func main() {
	supabase.InitDB()
	err := users.RegisterFCMTokenService("00000000-0000-0000-0000-000000000001", users.RegisterFCMTokenRequest{
		Token: "test-token-123",
		DeviceInfo: "test-device",
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("Success")
}
