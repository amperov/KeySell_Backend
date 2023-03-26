package gen

import (
	"log"
	"math/rand"
	"time"
)

type Gen struct {
}

func (g *Gen) Generate() string {
	inp := "0123456789QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"

	newPass := ""
	rand.Int()
	for i := 1; i < 17; i++ {
		r := rand.New(rand.NewSource(int64(i) * time.Now().Unix()))
		Index := r.Intn(len(inp))

		newPass += string(inp[Index])
	}

	log.Println(newPass)
	return newPass
}
