package entcmd

import (
	"log"
	"os"
	"strings"
)

// RunDel run del
func RunDel(input string, entitys string) error {
	list := strings.Split(entitys, ",")
	for _, ct1 := range list {
		os.Remove(input + "/" + ct1 + ".go")
		os.Remove(input + "/" + ct1 + "_create.go")
		os.Remove(input + "/" + ct1 + "_delete.go")
		os.Remove(input + "/" + ct1 + "_query.go")
		os.Remove(input + "/" + ct1 + "_update.go")
		os.RemoveAll(input + "/" + ct1)
		log.Println("delete: " + input + "/" + ct1)
	}
	os.Remove(input + "/client.go")
	return nil
}
