package business

import (
	"fmt"
)

func ProcessBookMessageEvent(message string) {
	fmt.Printf("Received Message Body %s \n", message)
}
