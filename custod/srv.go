package main

import (
	"./clockvec"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Servidor struct {
	Name   string
	Status int
}

type Ask struct {
	ID int
}

type Object struct {
	Name string
	Len  int
	Text []byte
}

var idagente = 0
var idobjeto = 0

var mutexAgent = &sync.Mutex{}
var mutexObject = &sync.Mutex{}
var mutexCheck = &sync.Mutex{}
var home = os.Getenv("HOME")

func (s *Servidor) Agente(n *Ask, numAgent *int) error {
	mutexAgent.Lock()
	defer mutexAgent.Unlock()

	idagente++
	*numAgent = idagente
	fmt.Println("Agente numero: ", idagente)
	f, err := os.OpenFile(home+"/srvcustod/agents", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err)
		os.Exit(1)
	}
	f.Write([]byte(strconv.Itoa(idagente)))
	f.WriteString("\n")
	return nil
}

func (s *Servidor) Objeto(n *Ask, numObject *int) error {
	if idagente == 0 || n.ID == 0 {
		fmt.Println("No se puede suministrar objeto, agente sin id")
	} else {
		mutexObject.Lock()
		defer mutexObject.Unlock()
		idobjeto++
		*numObject = idobjeto
		fmt.Println("Objeto numero: ", *numObject)
		f, err := os.OpenFile(home+"/srvcustod/objects", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s \n", err)
			os.Exit(1)
		}
		f.Write([]byte(strconv.Itoa(idobjeto)))
		f.WriteString("\n")
	}
	return nil
}

func openFile(file string) *clockvec.Clock {
	var c2 *(clockvec.Clock)

	text, _ := ioutil.ReadFile(file)
	dec := json.NewDecoder(strings.NewReader(string(text[:len(text)])))
	for {
		var c1 *(clockvec.Clock)
		if err2 := dec.Decode(&c1); err2 == io.EOF {
			break
		} else if err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: %s \n", err2)
			os.Exit(1)
		}
		c2 = c1
	}
	return c2
}

func (s *Servidor) ReceiveFile(n *Object, objeto *string) error {
	//DEVOLVER FALLO DE SUBIDA AL CLIENTE
	mutexCheck.Lock()
	defer mutexCheck.Unlock()
	if status, err := os.Stat(home+ "/srvcustod/" + n.Name); err == nil {
		if n.Len == int(status.Size()) {
			err := errors.New("Objeto existente en srv")
			fmt.Fprintf(os.Stderr, "OBJETO DESECHADO: %s --> %s \n", n.Name, err)
			*objeto = ("DESECHADO: " + err.Error())
			return nil
		}
		var c2 *(clockvec.Clock)
		dec := json.NewDecoder(strings.NewReader(string(n.Text)))
		for {
			var c1 *(clockvec.Clock)
			if err2 := dec.Decode(&c1); err2 == io.EOF {
				break
			} else if err2 != nil {
				fmt.Fprintf(os.Stderr, "Error: %s \n", err2)
				os.Exit(1)
			}
			c2 = c1
		}
		c3 := openFile(home + "/srvcustod/" + n.Name)
		//Ultimo objeto pasado por cliente: c2
		//Ultimo objeto servidor: c3
		if err, check := c2.CheckClock(c3); check == true {
			fmt.Println("OBJETO SUBIDO: ", n.Name)
			ioutil.WriteFile(home+"/srvcustod/"+n.Name, n.Text, 0664)
			*objeto = "OK"
		} else {
			fmt.Fprintf(os.Stderr, "OBJETO DESECHADO: %s --> %s \n", n.Name, err)
			*objeto = ("DESECHADO: " + err.Error())
			//MANDAR AL CLIENTE EL FALLO
		}
	} else {
		f, _ := os.OpenFile(home+ "/srvcustod/"+n.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		f.Write(n.Text)
		fmt.Println("OBJETO SUBIDO: ", n.Name)
		*objeto = "OK"
	}
	return nil
}

func recover(path string) int {
	var line int
	if _, err := os.Stat(path); err == nil {
		fd, _ := os.Open(path)
		for {
			_, err2 := fmt.Fscanf(fd, "%d\n", &line)
			if err2 != nil {
				if err2 == io.EOF {
					break
				}
				fmt.Fprintf(os.Stderr, "Error: %s \n", err)
				os.Exit(1)
			}
		}
		fd.Close()
	}
	return line
}

func main() {
	s := new(Servidor)
	rpc.Register(s)
	rpc.HandleHTTP()
	ip := os.Args[1]
	port := os.Args[2]
	if _,err := os.Stat(home+"/srvcustod/"); os.IsNotExist(err){
		os.MkdirAll(home+"/srvcustod/",os.ModePerm)
	}
	idagente = recover(home+"/srvcustod/agents")
	if idagente > 0 {
		fmt.Println("Agente guardado: ", idagente)
	}
	idobjeto = recover(home+"/srvcustod/objects")
	if idobjeto > 0 {
		fmt.Println("Ultimo objeto suministrado: ", idobjeto)
	}

	l, e := net.Listen("tcp", ip+":"+port)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", e)
		os.Exit(1)
	}
	fmt.Println("CONECTADO")

	go http.Serve(l, nil)

	time.Sleep(1 * time.Hour)
}
