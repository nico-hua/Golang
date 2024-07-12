package many_many

import (
	"producer_consumer/out"
	"time"
)

type Task struct{
	ID int64
}

func (t *Task) run(){
	out.Println(t.ID)
}

var taskCh = make(chan Task, 10)
var done = make(chan struct{})

const taskNum int64 = 10000

func producer(wo chan<- Task, done chan struct{}){
	var i int64 
	for{
		if i > taskNum {
			i = 0
		}
		i++
		t := Task{
			ID: i,
		}
		wo <- t
		select {
		case wo <- t:
		case <- done:
			out.Println("生产者退出")
			return
		}
	}
}

func consumer(ro <-chan Task, done chan struct{}){
	for{
		select{
		case t := <-ro:
			if t.ID != 0{
				t.run()
			}
		case <- done:
			for t := range ro{
				if t.ID != 0 {
					t.run()
				}
			}
			out.Println("消费者退出")
			return
		}
	}
}

func Exec(){
	go producer(taskCh, done)
	go producer(taskCh, done)
	go producer(taskCh, done)
	go producer(taskCh, done)
	go producer(taskCh, done)

	go consumer(taskCh, done)
	go consumer(taskCh, done)
	go consumer(taskCh, done)
	go consumer(taskCh, done)
	go consumer(taskCh, done)

	time.Sleep(time.Second*5)
	close(done)
	time.Sleep(time.Second * 1)
	close(taskCh)
	time.Sleep(time.Second*5)
}

