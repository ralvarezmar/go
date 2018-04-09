package main

import (
	"bufio"
	"./clockvec"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Object struct {
	Name string
	Len  int
	Text []byte
}

type Ask struct {
	ID int
}

var idAgente int
var home = os.Getenv("HOME")

func askAgent(client *rpc.Client) {
	if _, err := os.Stat(home+"/clicustod/myid"); err == nil {
		fmt.Println("Agente ID: ", idAgente)
	} else {
		argumentos := &Ask{0}
		var reply int
		err = client.Call("Servidor.Agente", argumentos, &reply)
		if err != nil {
			fmt.Println("Agente error:", err)
		}
		idAgente = reply
		fmt.Println("Agente ID: ", int(idAgente))
		ioutil.WriteFile(home+"/clicustod/myid", []byte(strconv.Itoa(idAgente)), 0664)
	}
}

func askObject(client *rpc.Client) {
	argumentos := &Ask{idAgente}
	var reply int
	err := client.Call("Servidor.Objeto", argumentos, &reply)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err)
	}
	idObjeto := reply
	if idAgente == 0 {
		fmt.Fprintf(os.Stderr,"Fallo al suministrar ID. Intentelo de nuevo")
	} else { //CREAR OBJETO
		fmt.Println("Objeto: ", idObjeto)
		c := clockvec.NewClock(idAgente)
		var mark []byte

		mark, _ = json.Marshal(c)
		f, _ := os.OpenFile(home+"/clicustod/objeto"+strconv.Itoa(idObjeto)+".cst", os.O_CREATE| os.O_APPEND|os.O_WRONLY, 0664)
		f.Write(mark)
		f.WriteString("\n")
	}
}

func sendFile(client *rpc.Client, file string) {
	var reply string

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err)
		return
	}
	fileInfo, _ := f.Stat()

	fileSize, _ := strconv.Atoi(strconv.FormatInt(fileInfo.Size(), 10))
	fileName := fileInfo.Name()
	nombre := (fileName)
	buffer := make([]byte, fileSize)
	for {
		_, err = f.Read(buffer)
		if err == io.EOF {
			break
		}
	}
	arguments := &Object{nombre, fileSize, buffer}
	fmt.Printf("Subiendo %s --> ", nombre)
	//errObj := new(ErrObject)
	err = client.Call("Servidor.ReceiveFile", arguments, &reply)
	fmt.Println(reply)
	//fmt.Fprintf(os.Stderr,"Objeto no subido al server -->: %s \n", &reply)
}

func searchObjects(client *rpc.Client) {
	files, err1 := filepath.Glob(home+ "/clicustod/*.cst")
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err1)
	}
	for _, v := range files {
		sendFile(client, v)
	}
}

func writeObject(id string) {
	var c2 *(clockvec.Clock)
	if _, err := os.Stat(home+"/clicustod/myid"); err != nil {
		panic("Agente sin id")
	}

	if _, err := os.Stat(home+"/clicustod/objeto" + id + ".cst"); err == nil {
		text, _ := ioutil.ReadFile(home+"/clicustod/objeto" + id + ".cst")
		//text = strings.Replace(text, "\n", "", -1)
		dec := json.NewDecoder(strings.NewReader(string(text[:len(text)])))
		for {
			var c1 *(clockvec.Clock)
			if err2 := dec.Decode(&c1); err2 == io.EOF {
				break
			} else if err2 != nil {
				panic(err2)
			}
			c2 = &clockvec.Clock{c1.Vector, ""}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Introduzca descripcion: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.Replace(desc, "\n", "", -1)
	c2.SumValue(idAgente, desc)
	mark, _ := json.Marshal(c2)
	f, _ := os.OpenFile(home+"/clicustod/objeto"+id+".cst", os.O_APPEND|os.O_WRONLY, 0664)
	f.Write(mark)
	f.WriteString("\n")
}

func main() {
	//Recupero estado

	if _,err := os.Stat(home+"/clicustod/"); os.IsNotExist(err){
		os.MkdirAll(home+"/clicustod/",os.ModePerm)
	}

		//-------------------PRUEBA
		//---------------
	if len(os.Args)==5{
		home="/home/ruben/Distribuidos"
	}

	if _, err := os.Stat(home+"/clicustod/myid"); err == nil {
		dat, e := ioutil.ReadFile(home+"/clicustod/myid")
		if e != nil {
			panic(e)
		}
		idAgente, _ = strconv.Atoi(string(dat))
	}

	if os.Args[1] == "-m" {
		idobjeto := os.Args[2]
		writeObject(idobjeto)
		return
	}

	ip := os.Args[2]
	port := os.Args[3]


	client, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("dialing:", err)
	}

	if os.Args[1] == "-a" {
		searchObjects(client)
		return
	}

	if os.Args[1] == "-g" {
		askAgent(client)
	}

	if os.Args[1] == "-c" {
		askObject(client)
	}

	time.Sleep(2 * time.Second)
}
