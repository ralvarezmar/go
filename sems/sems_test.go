package sems

import (
	"fmt"
	"testing"
	"time"
	"sync/atomic"
)

func Test1(t *testing.T) {
	fmt.Println("------PRIMER TEST/BLOQUEO----------")
	var sem = NewSem(0)

	for i := 0; i < 1; i++ {
		go f2(sem)
	}
	fmt.Println("Bloqueo 2 segundos")
	time.Sleep(2 * time.Second)
	sem.Up()
	fmt.Println("Hola")
	time.Sleep(2 * time.Second)
}
func f2(s *Sem) {
	s.Down()
	fmt.Println("Adios")
}

func Test2(t *testing.T) {
	fmt.Println("------SEGUNDO TEST/MUTEX----------")
	sem := NewSem(1)
	contador := 0
	for i := 0; i < 100; i++ {
		go func() {
			sem.Down()
			time.Sleep(50 * time.Millisecond)
			contador++
			fmt.Println("Contador ", contador)
			sem.Up()
		}()
	}
	time.Sleep(6 * time.Second)
	if contador == 100{
		fmt.Println("TEST OK")
	}else{
		t.Error("Contador no esperado")
	}
}

func Test3(t *testing.T){
	fmt.Println("------TERCER TEST/N HILOS MAX----------")
	sem := NewSem(5)
	var contador int32
	contador = 0
	for i := 0; i < 50;i++{
		go func(){
			sem.Down()
			atomic.AddInt32(&contador,1);
			//contador++//addInt32
			sem.Up()
			//fmt.Println("Contador: ",contador)
			}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("Contador: ",contador)
	if contador == 50{
		fmt.Println("TEST OK")
	}else{
		t.Error("Contador no esperado")
	}
}

func Test4(t *testing.T) {
	fmt.Println("------CUARTO TEST----------")
	sem := NewSem(10)
	for i := 0; i <= 50; i++ {
		go func(n int) {
			sem.Down()
			fmt.Println(n)
		}(i)
	}
	time.Sleep(3 * time.Second)
}
