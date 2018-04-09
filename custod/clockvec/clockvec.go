package clockvec

import (
	"errors"
	"strconv"

	)

type Value struct {
	Id    int
	Count int
}

type Clock struct {
	Vector []Value
	Text   string
}

func CheckID(id int, valores []Value) (bool, int) {
	for i, _ := range valores {
		if id == valores[i].Id {
			return true, i
		}
	}
	return false, 0
}

func NewClock(id int) *Clock {
	vect := make([]Value, 1)
	vect[0] = Value{id, 1}
	return &Clock{vect, "CREATE"}
}

func (c *Clock) SumValue(id int, text string) {
	//BUSCO SI EL CLIENTE YA ESTÁ EN EL ARRAY
	check, i := CheckID(id, c.Vector)
	if check == true {
		c.Vector[i].Count++
		c.Text = text
	} else {
		//SI NO HA ENCONTRADO EL ID, CREO NUEVO CAMPO
		newClock := make([]Value, len(c.Vector)+1)
		newClock[len(c.Vector)] = Value{id, 1}
		copy(newClock, c.Vector)
		c.Vector = newClock
		c.Text = text
	}
}

func checkVec(cMax *Clock, cMin *Clock) (bool,int,bool) {
	clockOk := true
	foundId := false
	badId := 0
	for i, _ := range cMax.Vector {
		foundId = false //si hay algun id que no encuentra foundId se queda en false
		for j, _ := range cMin.Vector {
			//badId=cMax.Vector[j].Id
			if cMax.Vector[i].Id == cMin.Vector[j].Id {
				foundId = true
				if cMax.Vector[i].Count >= cMin.Vector[j].Count && clockOk == true {
					clockOk = true
				} else if cMax.Vector[i].Count < cMin.Vector[j].Count{
					clockOk = false
					badId = cMax.Vector[i].Id
				}
			}
		}
		if foundId == false {
			return true,badId,foundId
		}
	}
	return clockOk,badId,foundId
}

func (c1 *Clock) CheckClock(c2 *Clock) (error,bool) {
	lenc1 := len(c1.Vector)
	lenc2 := len(c2.Vector)
	var check bool
	var badId int
	var found bool
	var err error
	if lenc1 >= lenc2 {
		check,badId,found = checkVec(c1, c2)
	} else {
		check,badId,found = checkVec(c2, c1)
		check = !check
	}
	if check == false && found == true{
		err = errors.New("Reloj inconsistente(Problema en agente:"+strconv.Itoa(badId)+")")
 }else if check == false && found == false{
	 err = errors.New("Reloj concurrente")

 }
	return err,check
}
