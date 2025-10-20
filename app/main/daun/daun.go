package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Daun struct {
	User_id int
	Exp     time.Time
}

func main() {
	var d Daun = Daun{12, time.Now()}
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(b)
	var j Daun
	json.Unmarshal(b, &j)
	fmt.Println(j)
}
