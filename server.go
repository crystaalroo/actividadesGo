package main

import (
	"errors"
	"fmt"
	. "fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Alumno struct {
	Grades map[string]float64
	Sum    float64
	sync.Mutex
}

type Server struct {
	Alumnos  map[string]*Alumno
	Materias map[string]float64
	B        sync.Mutex
}

type NuC struct {
	Name, Materia string
	Calif         float64
}

func (u *Server) agregarAlumno(name string) {
	u.B.Lock()
	u.Alumnos[name] = &Alumno{make(map[string]float64), 0.0, sync.Mutex{}}
	u.B.Unlock()
}

func (u *Server) agregarMateria(name string) {
	u.B.Lock()
	u.Materias[name] = 0.0
	u.B.Unlock()
}

func (u *Server) AgregarCalif(c NuC, reply *string) error {
	_, b := u.Alumnos[c.Name]
	if !b {
		u.agregarAlumno(c.Name)
	}
	_, b = u.Materias[c.Materia]
	if !b {
		u.agregarMateria(c.Materia)
	}
	a := u.Alumnos[c.Name]
	a.Lock()
	_, b = a.Grades[c.Materia]
	if !b {
		a.Grades[c.Materia] = c.Calif
		a.Sum += c.Calif
		u.B.Lock()
		u.Materias[c.Materia] += c.Calif
		u.B.Unlock()
		a.Unlock()
		*reply = "La calificacion del alumno " + c.Name + " en la materia " + c.Materia + " ha sido registrada correctamente"
		return nil
	}
	a.Unlock()
	return errors.New("El alumno " + c.Name + " ya tenia una calificacion en la materia " + c.Materia)
}

func (u *Server) PromedioAlumno(nombre string, reply *float64) error {
	_, b := u.Alumnos[nombre]
	if !b {
		u.agregarAlumno(nombre)
	}
	u.B.Lock()
	totm := len(u.Materias)
	sum := u.Alumnos[nombre].Sum
	if totm == 0 {
		totm = 1
	}
	u.B.Unlock()
	*reply = sum / float64(totm)
	return nil
}

func (u *Server) PromedioMateria(nombre string, reply *float64) error {
	_, b := u.Materias[nombre]
	if !b {
		u.agregarMateria(nombre)
	}
	u.B.Lock()
	tota := len(u.Alumnos)
	sum := u.Materias[nombre]
	if tota == 0 {
		tota = 1
	}
	u.B.Unlock()
	*reply = sum / float64(tota)
	return nil
}

func (u *Server) PromedioGeneral(idk int, reply *float64) error {
	if idk != 0 {
		return nil
	}
	u.B.Lock()
	tota := len(u.Alumnos)
	totm := len(u.Materias)
	if tota == 0 {
		tota = 1
	}
	if totm == 0 {
		totm = 1
	}
	u.B.Unlock()
	var sum float64
	for _, e := range u.Alumnos {
		sum += e.Sum / float64(totm)
	}
	*reply = sum / float64(tota)
	return nil
}

func (u *Server) print() {
	for {
		Println(time.Now())
		Println("Alumnos registrados:")
		for a, _ := range u.Alumnos {
			Println(a)
		}
		Println("Materias registrados:")
		for a, _ := range u.Materias {
			Println(a)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func server() {
	s := Server{make(map[string]*Alumno), make(map[string]float64), sync.Mutex{}}
	rpc.Register(&s)
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	var inn string
	Scanln(&inn)
}
