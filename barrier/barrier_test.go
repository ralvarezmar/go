package barrier

import (
	"fmt"
	"testing"
	"time"
	"sync/atomic"
)

func Test1(t *testing.T) {
	fmt.Println("------PRIMER TEST----------")
	var barrier = NewBarrier(5)
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		go f(barrier)
	}
	time.Sleep(5 * time.Second)
}

func f(b *Barrier) {
	time.Sleep(1 * time.Second)
	fmt.Println("Espero")
	b.Wait()
	fmt.Println("Lets go!")
}

func Test2(t *testing.T) {
	fmt.Println("------SEGUNDO TEST----------")
	var barrier = NewBarrier(5)
	var counter int32
	counter = 0
	for i := 0; i < 12; i++ {
		go func(n int){
			barrier.Wait()
			atomic.AddInt32(&counter,1);
			fmt.Println(counter)

			//counter++ //addInt32
		}(i)
	}
	time.Sleep(2 * time.Second)
	fmt.Println("COUNTER:", counter)
	if counter==10{
		fmt.Println("Test ok")
	}else{
		t.Error("Contador no esperado")
	}
}

func Test3(t *testing.T) {
	fmt.Println("------TERCER TEST----------")
	var barrier = NewBarrier(4)
	var counter int32
	counter = 0
	for i := 0; i < 16; i++ {
		go func(n int){
			barrier.Wait()
			atomic.AddInt32(&counter,1);

		}(i)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("COUNTER:", counter)
	if counter != 16{
		t.Error("Contador no esperado")
	}else{
		fmt.Println("Test ok")
	}
}
