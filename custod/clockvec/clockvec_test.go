package clockvec

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

func Test1(t *testing.T) {
	fmt.Println("------PRIMER TEST----------")
	var c = NewClock(1)
	var c2 = NewClock(1)
	c.SumValue(1,"Escribo1")
	c.SumValue(2,"Escribo2")

	c2.SumValue(1,"Escribo1")
	c2.SumValue(2,"Escribo2")
	c2.SumValue(2,"Escribo3")
	fmt.Println("Clock1:", c)
	fmt.Println("Clock2:", c2)
	
	err,check := c.CheckClock(c2)

	fmt.Println(err,check)
}

func Test2(t *testing.T) {
	fmt.Println("------SEGUNDO TEST----------")
	var c = NewClock(1)
	var c2 = NewClock(1)
	c.SumValue(1,"Escribo1")
	c.SumValue(1,"Escribo2")

	c2.SumValue(1,"Escribo1")
	c2.SumValue(2,"Escribo2")
	c2.SumValue(2,"Escribo3")
	fmt.Println("Clock1:", c)
	fmt.Println("Clock2:", c2)
	
	err,check := c.CheckClock(c2)

	fmt.Println(err,check)
}

func Test3(t *testing.T) {
	fmt.Println("------TERCER TEST----------")
	var c = NewClock(1)
	var c2 = NewClock(1)
	c.SumValue(1,"Escribo1")
	c.SumValue(2,"Escribo2")
	c.SumValue(2,"Escribo3")
	c.SumValue(3,"Escribo4")
	c.SumValue(1,"Escribo5")

	c2.SumValue(1,"Escribo1")
	c2.SumValue(2,"Escribo2")
	c2.SumValue(2,"Escribo3")
	fmt.Println("Clock1:", c)
	fmt.Println("Clock2:", c2)
	
	err,check := c.CheckClock(c2)

	fmt.Println(err,check)
}


func Test4(t *testing.T) {
	fmt.Println("------CUARTO TEST----------")
	var c = NewClock(1)
	var c2 = NewClock(1)
	var mutexObject = &sync.Mutex{}


	c.SumValue(1,"Escribo1")
	c.SumValue(2,"Escribo2")
	c.SumValue(2,"Escribo3")
	c.SumValue(3,"Escribo4")
	c.SumValue(1,"Escribo5")

	c2.SumValue(1,"Escribo1")
	c2.SumValue(2,"Escribo2")
	c2.SumValue(2,"Escribo3")
	fmt.Println("Clock1:", c)
	fmt.Println("Clock2:", c2)
	
	err,check := c.CheckClock(c2)

	fmt.Println(err,check)
	for i:=1; i<5;i++{
		go func(i int){ 
			mutexObject.Lock()
			c2.SumValue(i,"Escribo1")
			c.SumValue(i,"Escribo1")
			fmt.Println(c,c2)
			err,check := c2.CheckClock(c)
			fmt.Println(err,check)
			mutexObject.Unlock()
		}(i)
	}
	time.Sleep(2 * time.Second)
}
