package main

import (
	"fmt"
	"val-snitch/internal/auth"
)

func main() {

	log_info := auth.Get_client_info()
	fmt.Printf("logInfo: %+v\n", log_info)

}
