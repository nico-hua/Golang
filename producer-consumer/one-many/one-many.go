package one_many

import (
	"producer_consumer/out"
	"sync"
)

type Task struct{
	ID int64
}

func (t *Task) run(){
	out.Println(t.ID)
}

var taskCh = make(chan Task, 10)

const taskNum int64 = 1000

func producer(wo chan<- Task){
	var i int64
	for i = 1; i <= taskNum;i++{
		t:=Task{
			ID:i,
		}
		wo<-t
	}
	close(wo)
}

func consumer(ro <-chan Task){
	for t:= range ro {
		if t.ID != 0 {
			t.run()
		}
	}
}

func Exec(){
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func (wg *sync.WaitGroup)  {
		defer wg.Done()
		producer(taskCh)
	}(wg)

	var i int64 
	for i=0; i<taskNum; i++ {
		if i%100 == 0 {
			wg.Add(1)
			go func (wg *sync.WaitGroup)  {
				defer wg.Done()
				consumer(taskCh)
			}(wg)
		}
	}
	
	wg.Wait()
	out.Println("执行结束")
}