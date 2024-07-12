package main

import (
	"os"
	"os/signal"
	"producer_consumer/out"
	"syscall"
	"producer_consumer/many-many"
)

func main(){
	o := out.NewOut()
	go o.OutPut()

	many_many.Exec()
	
	sig:=make(chan os.Signal,1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}