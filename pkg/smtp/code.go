package smtp

import (
	"fmt"
	"math/rand"
)

func GenerateCode() string {
	return fmt.Sprintf("%d", (rand.Intn(8998)+1000)%9999)
}
