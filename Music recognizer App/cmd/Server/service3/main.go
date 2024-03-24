package main

import (
	"awesomeProject/internal/utils"
	"fmt"
)

func main() {
	err := utils.Publish("89484")
	if err != nil {
		fmt.Println(err)
	}
	//execute()
	//for range time.Tick(5 * time.Minute) {
	//	execute()
	//}
}
