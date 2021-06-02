package main

import (
	"encoding/json"
	. "fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Alumno struct {
	ID     uint64
	Nombre string `json:"nombre"`
	Calif  map[uint64]float64
}

type Materia struct {
	ID     uint64
	Nombre string `json:"nombre"`
}

type Calif struct {
	Al   uint64  `json:"id-al"`
	Ma   uint64  `json:"id-ma"`
	Cali float64 `json:"calif"`
}

type Clase struct {
	Alumnos  map[uint64]*Alumno
	Materias map[uint64]*Materia
}

func getID() uint64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Uint64()
}

var clase Clase

func Add_Alumno(alumno Alumno) []byte {
	jsonData := []byte(`{"code": "ok"}`)
	clase.Alumnos[alumno.ID] = &alumno
	return jsonData
}

func Add_Materia(materia Materia) []byte {
	jsonData := []byte(`{"code": "ok"}`)
	clase.Materias[materia.ID] = &materia
	return jsonData
}

func Add_Calif(calif Calif) []byte {
	_, ok := clase.Alumnos[calif.Al]
	if ok == false {
		return []byte(`{"code": "noexistealumno"}`)
	}
	_, ok = clase.Materias[calif.Ma]
	if ok == false {
		return []byte(`{"code": "noexistemateria"}`)
	}
	_, ok = clase.Alumnos[calif.Al].Calif[calif.Ma]
	if ok == true {
		return []byte(`{"code": "yatienecalif"}`)
	}
	clase.Alumnos[calif.Al].Calif[calif.Ma] = calif.Cali
	return []byte(`{"code": "ok"}`)
}

func Delete_Alumno(id uint64) []byte {
	_, ok := clase.Alumnos[id]
	if ok == false {
		return []byte(`{"code": "noexiste"}`)
	}
	delete(clase.Alumnos, id)
	return []byte(`{"code": "ok"}`)
}

func Update_Calif(calif Calif) []byte {
	_, ok := clase.Alumnos[calif.Al]
	if ok == false {
		return []byte(`{"code": "noexistealumno"}`)
	}
	_, ok = clase.Alumnos[calif.Al].Calif[calif.Ma]
	if ok == false {
		return []byte(`{"code": "noexistemateria"}`)
	}
	clase.Alumnos[calif.Al].Calif[calif.Ma] = calif.Cali
	return []byte(`{"code": "ok"}`)
}

func Get() ([]byte, error) {
	jsonData, err := json.MarshalIndent(clase, "", "    ")
	if err != nil {
		return jsonData, nil
	}
	return jsonData, err
}

func GetAlumno(id uint64) ([]byte, error) {
	jsonData := []byte(`{}`)
	al, ok := clase.Alumnos[id]
	if ok == false {
		return jsonData, nil
	}
	jsonData, err := json.MarshalIndent(al, "", "    ")
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func alumnos(res http.ResponseWriter, req *http.Request) {
	Println(req.Method)
	switch req.Method {
	case "POST":
		var alumno Alumno

		err := json.NewDecoder(req.Body).Decode(&alumno)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		alumno.ID = getID()
		alumno.Calif = make(map[uint64]float64)
		res_json := Add_Alumno(alumno)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}

func materias(res http.ResponseWriter, req *http.Request) {
	Println(req.Method)
	switch req.Method {
	case "POST":
		var materia Materia

		err := json.NewDecoder(req.Body).Decode(&materia)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		materia.ID = getID()
		res_json := Add_Materia(materia)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}
func alumnoID(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(req.URL.Path, "/clase/alumnos/"), 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	Println(req.Method, id)
	switch req.Method {
	case "GET":
		res_json, err := GetAlumno(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	case "DELETE":
		res_json := Delete_Alumno(id)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}

func miclase(res http.ResponseWriter, req *http.Request) {
	Println(req.Method)
	switch req.Method {
	case "GET":
		res_json, err := Get()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}

func calif(res http.ResponseWriter, req *http.Request) {
	Println(req.Method)
	switch req.Method {
	case "POST":
		var calif Calif

		err := json.NewDecoder(req.Body).Decode(&calif)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res_json := Add_Calif(calif)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	case "PUT":
		var calif Calif
		err := json.NewDecoder(req.Body).Decode(&calif)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res_json := Update_Calif(calif)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json) // retornamos el JSON respuesta al cliente
	}
}
func main() {
	clase = Clase{
		Alumnos:  map[uint64]*Alumno{},
		Materias: map[uint64]*Materia{},
	}
	http.HandleFunc("/clase/alumnos", alumnos)
	http.HandleFunc("/clase/alumnos/", alumnoID)
	http.HandleFunc("/clase/materias", materias)
	http.HandleFunc("/clase", miclase)
	http.HandleFunc("/clase/", miclase)
	http.HandleFunc("/clase/calif", calif)
	Println("Corriendo RESTful API...")
	http.ListenAndServe(":9000", nil)
}
