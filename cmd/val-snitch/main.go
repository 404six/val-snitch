package main

import (
	"fmt"
	"val-snitch/internal/auth"
)

func main() {

	log_info := auth.Get_client_info()
	fmt.Printf("log_info: %+v\n", log_info)

	access_token, err := auth.Auth_from_client()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	entitlement_token, err := auth.Get_entitlement(access_token)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("entitlement_token: %+v\n", entitlement_token)
}
