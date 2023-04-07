package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	uuid, _ := uuid.Parse("000004b1f3764dd5a3330d175e54630f")

	fmt.Println(uuid.String())
	fmt.Println(uuid.ID())
}
