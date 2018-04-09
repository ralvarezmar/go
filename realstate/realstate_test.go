package realstate

import (
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var (
		r *RealState = NewRealState()
	)
	dni := DNI{123, "x"}
	nombre := Name{"Pepe", "Perez", nil}
	o := Owner{dni, nombre}

	dni2 := DNI{456, "x"}
	nombre2 := Name{"Juan", "Manolo", nil}
	o2 := Owner{dni2, nombre2}
	dni3 := DNI{987, "v"}
	nombre3 := Name{"Pablo", "Hernandez", nil}
	o3 := Owner{dni3, nombre3}
	direccion := Address{"Agua", 9, "B", 23123, nil}
	h := House{01, direccion}
	direccion2 := Address{"Fuego", 45, "Z", 234, nil}
	h2 := House{02, direccion2}
	direccion3 := Address{"Aire", 28, "V", 345, nil}
	h3 := House{03, direccion3}

	r.AddNewPair(&o, &h)

	go fmt.Println(r.AddOwner(&o2))
	go fmt.Println(r.AddOwner(&o3))
	go fmt.Println(r.AddHouse(dni, &h2))
	go fmt.Println(r.AddHouse(dni3, &h3))

	go fmt.Println(r.GetHouses(dni))
	go fmt.Println(r.GetOwners(02))

	go fmt.Println(r.GetOwner(dni2))

	go fmt.Println(r.GetHouse(01))

	go fmt.Println(r.ChangeOwner(&o2))

	go fmt.Println(r.ChangeHouse(&h2))

	go fmt.Println(r.DelOwner(dni3))
	go fmt.Println(r.DelHouse(02))

	go fmt.Println(r.NukeOwner(dni))

	time.Sleep(2 * time.Second)
}
