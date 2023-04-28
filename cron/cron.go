package main

import (
	"fmt"
	"github.com/robfig/cron"
)

func main() {

	i := 0
	c := cron.New()
	spec := "*/2 * * * * ?"
	err := c.AddFunc(spec, func() {
		i++
		fmt.Println("cron times : ", i)
	})
	if err != nil {
		fmt.Errorf("AddFunc error : %v", err)
		return
	}
	c.Start()

	defer c.Stop()
	select {}
}
