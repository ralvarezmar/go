package realstate

import (
	"errors"
	"sync"
)

type DNI struct {
	number int
	letter string
}

type Owner struct {
	dni  DNI
	name Name
}

type Name struct {
	fname   string
	surname string
	casa    []int
}

type House struct {
	id        int
	direccion Address
}

type Address struct {
	street      string
	number      int
	letter      string
	pc          int
	propietario []DNI
}

type RealState struct {
	ownerList  map[DNI]Name
	housesList map[int]Address
	mutexOwner *sync.Mutex
	mutexHouse *sync.Mutex
}

func NewRealState() *RealState {
	ownerList := make(map[DNI]Name)
	housesList := make(map[int]Address)
	mutexOwner := &sync.Mutex{}
	mutexHouse := &sync.Mutex{}
	return &RealState{ownerList, housesList, mutexOwner, mutexHouse}
}

func (r *RealState) AddNewPair(o *Owner, h *House) error {
	r.mutexOwner.Lock()
	r.mutexHouse.Lock()
	defer r.mutexOwner.Unlock() //Primero unlock de House, despues owner
	defer r.mutexHouse.Unlock()

	if _, ok := r.ownerList[o.dni]; ok {
		return errors.New("Owner error")
	}
	if _, ok := r.housesList[h.id]; ok {
		return errors.New("House error")
	}
	o.name.casa = append(o.name.casa, h.id)
	h.direccion.propietario = append(h.direccion.propietario, o.dni)
	r.ownerList[o.dni] = o.name
	r.housesList[h.id] = h.direccion

	return nil
}

func (r *RealState) AddOwner(o *Owner) error {
	r.mutexOwner.Lock()
	defer r.mutexOwner.Unlock()
	if _, ok := r.ownerList[o.dni]; ok {
		return errors.New("Owner error")
	}
	r.ownerList[o.dni] = o.name
	return nil
}

func (r *RealState) AddHouse(dni DNI, h *House) error {
	r.mutexOwner.Lock()
	r.mutexHouse.Lock()
	defer r.mutexOwner.Unlock() //Primero unlock de House, despues owner
	defer r.mutexHouse.Unlock()
	if _, ok := r.housesList[h.id]; ok {
		return errors.New("Casa existe")
	}
	if _, ok := r.ownerList[dni]; !ok {
		return errors.New("Owner error, usuario no existe")
	}
	persona := r.ownerList[dni]
	persona.casa = append(persona.casa, h.id)
	h.direccion.propietario = append(h.direccion.propietario, dni)

	r.ownerList[dni] = persona
	r.housesList[h.id] = h.direccion
	return nil
}
func (r *RealState) GetHouses(dni DNI) (listId []int, err error) {
	r.mutexOwner.Lock()
	defer r.mutexOwner.Unlock()
	if _, ok := r.ownerList[dni]; !ok {
		return listId, errors.New("Owner error, usuario no existe")
	}
	persona := r.ownerList[dni]
	listId = persona.casa
	return listId, nil
}

func (r *RealState) GetOwners(id int) (listDni []DNI, err error) {
	r.mutexHouse.Lock()
	defer r.mutexHouse.Unlock()
	if _, ok := r.housesList[id]; !ok {
		return nil, errors.New("House error, casa no existe")
	}
	casa := r.housesList[id]
	return casa.propietario, nil
}

func (r *RealState) GetOwner(dni DNI) (o Owner, err error) {
	r.mutexOwner.Lock()
	defer r.mutexOwner.Unlock()
	if _, ok := r.ownerList[dni]; !ok {
		return o, errors.New("Owner error, usuario no existe")
	}
	persona := r.ownerList[dni]
	person := Owner{dni, persona}
	return person, nil
}

func (r *RealState) GetHouse(id int) (h House, err error) {
	r.mutexHouse.Lock()
	defer r.mutexHouse.Unlock()
	if _, ok := r.housesList[id]; !ok {
		return h, errors.New("House error, casa no existe")
	}
	casa := r.housesList[id]
	home := House{id, casa}
	return home, nil
}

func (r *RealState) ChangeOwner(o *Owner) error {
	r.mutexOwner.Lock()
	defer r.mutexOwner.Unlock()
	if _, ok := r.ownerList[o.dni]; !ok {
		return errors.New("Owner error, usuario no existe")
	}
	o.name.fname = "Paco"
	r.ownerList[o.dni] = o.name
	return nil
}

func (r *RealState) ChangeHouse(h *House) error {
	r.mutexHouse.Lock()
	defer r.mutexHouse.Unlock()
	if _, ok := r.housesList[h.id]; !ok {
		return errors.New("Owner error, usuario no existe")
	}
	h.direccion = r.housesList[h.id]
	h.direccion.letter = "C"
	h.direccion.number = 56
	r.housesList[h.id] = h.direccion
	return nil
}

func (r *RealState) DelOwner(dni DNI) error {
	r.mutexOwner.Lock()
	defer r.mutexOwner.Unlock()
	if _, ok := r.ownerList[dni]; !ok {
		return errors.New("Owner error, usuario no existe")
	}
	user := r.ownerList[dni]
	r.mutexHouse.Lock()
	for _, v := range user.casa {
		h := r.housesList[v]
		for j, k := range h.propietario {
			if k.number == dni.number && k.letter == dni.letter {
				h.propietario = append(h.propietario[:j], h.propietario[j+1:]...)
			}
		}
		r.housesList[v] = h //añado el array editado sin el usuario
	}
	r.mutexHouse.Unlock()
	delete(r.ownerList, dni)
	return nil
}

func (r *RealState) DelHouse(id int) error {
	r.mutexOwner.Lock()
	r.mutexHouse.Lock()
	defer r.mutexOwner.Unlock() //Primero unlock de House, despues owner
	defer r.mutexHouse.Unlock()
	if _, ok := r.housesList[id]; !ok {
		return errors.New("House error, casa no existe")
	}
	casa := r.housesList[id]
	for _, v := range casa.propietario {
		o := r.ownerList[v]
		for j, _ := range o.casa {
			if o.casa[j] == id {
				o.casa = append(o.casa[:j], o.casa[j+1:]...)
			}
		}
		r.ownerList[v] = o
	}
	delete(r.housesList, id)
	return nil
}

func (r *RealState) NukeOwner(dni DNI) error {
	r.mutexOwner.Lock()
	r.mutexHouse.Lock()
	defer r.mutexOwner.Unlock() //Primero unlock de House, despues owner
	defer r.mutexHouse.Unlock()
	if _, ok := r.ownerList[dni]; !ok {
		return errors.New("Owner error, usuario no existe")
	}
	user := r.ownerList[dni]
	for _, v := range user.casa { //user.casa ids de viviendas
		casas := r.housesList[v]
		for _, k := range casas.propietario { //casas.propietario dnis de los usuarios
			propietario := r.ownerList[k]
			for i, _ := range propietario.casa {
				if propietario.casa[i] == v {
					propietario.casa = append(propietario.casa[:i], propietario.casa[i+1:]...)
				}
			}
			r.ownerList[k] = propietario
		}
		delete(r.housesList, v)
	}
	delete(r.ownerList, dni)
	return nil
}
