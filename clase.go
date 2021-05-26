package main

import (
	. "fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Html(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

func Root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	Fprintf(res, Html("index.html"))
}

type Calificacion struct {
	Calif   float64
	Alumno  string
	Materia string
}

type Alumno struct {
	Nombre string
	Calif  map[string]float64
	Suma   float64
}

type Clase struct {
	Calificaciones []Calificacion
	Alumnos        map[string]*Alumno
	Materias       map[string]float64
	Suma           float64
}

func (clase *Clase) Agregar(calif Calificacion) {
	clase.Suma -= clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Materias[calif.Materia] -= clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Suma -= clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Alumnos[calif.Alumno].Suma -= clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Alumnos[calif.Alumno].Calif[calif.Materia] = calif.Calif
	clase.Suma += clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Materias[calif.Materia] += clase.Alumnos[calif.Alumno].Calif[calif.Materia]
	clase.Alumnos[calif.Alumno].Suma += clase.Alumnos[calif.Alumno].Calif[calif.Materia]
}

var clase Clase

func AgregarCalif(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var err error
			if err = req.ParseForm(); err != nil {
				Fprintf(res, "ParseForm() error %v", err)
				return
			}
			Println(req.PostForm)
			calif := Calificacion{0, req.FormValue("alumno"), req.FormValue("materia")}

			if calif.Calif, err = strconv.ParseFloat(req.FormValue("calif"), 64); err != nil {
				Fprint(res, "Solo aceptamos numeros")
			}
			clase.Agregar(calif)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(res, Html("response.html"), "La calificacion fue modificada.")
		}
	case "GET":
		{
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(res, Html("agregarc.html"), ParseMateria(), ParseAlumno())
		}
	}
}

func ParseAlumno() string {
	alu := ""
	for key, _ := range clase.Alumnos {
		alu += "<option value=" + "\"" + key + "\"" + ">" + key + "</option>" + "\n"
	}
	return alu
}
func PromedioAlumno(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var err error
			if err = req.ParseForm(); err != nil {
				Fprintf(res, "ParseForm() error %v", err)
				return
			}
			Println(req.PostForm)
			name := req.FormValue("alumno")
			promedio := clase.Alumnos[name].Suma / float64(len(clase.Materias))
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(res, Html("response.html"), strconv.FormatFloat(promedio, 'f', -1, 64))
		}
	case "GET":
		{
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(res, Html("promalu.html"), ParseAlumno())
		}
	}
}

func ParseMateria() string {
	mat := ""
	for key, _ := range clase.Materias {
		mat += "<option value=" + "\"" + key + "\"" + ">" + key + "</option>" + "\n"
	}
	return mat
}

func PromedioMateria(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var err error
			if err = req.ParseForm(); err != nil {
				Fprintf(res, "ParseForm() error %v", err)
				return
			}
			Println(req.PostForm)
			materia := req.FormValue("materia")
			promedio := clase.Materias[materia] / float64(len(clase.Alumnos))
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(res, Html("response.html"), strconv.FormatFloat(promedio, 'f', -1, 64))
		}
	case "GET":
		{
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			Fprintf(
				res,
				Html("prommat.html"),
				ParseMateria(),
			)
		}
	}
}

func PromedioGeneral(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var err error
			if err = req.ParseForm(); err != nil {
				Fprintf(res, "ParseForm() error %v", err)
				return
			}
			Println(req.PostForm)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
		}
	case "GET":
		{
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			promedio := clase.Suma / float64(len(clase.Materias)*len(clase.Alumnos))
			Fprintf(res, Html("response.html"), strconv.FormatFloat(promedio, 'f', -1, 64))
		}
	}
}

func genAlumno() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdedfghijklmnopqrst")
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		output.WriteRune(chars[rand.Intn(len(chars))])
	}
	return output.String()
}
func main() {
	clase = Clase{
		Calificaciones: []Calificacion{},
		Alumnos:        map[string]*Alumno{},
		Materias: map[string]float64{
			"Historia":     0,
			"Programacion": 0,
			"Algebra":      0,
			"Algoritmia":   0,
		},
		Suma: 0,
	}
	for i := 0; i < 10; i += 1 {
		aux := genAlumno()
		clase.Alumnos[aux] = &Alumno{
			Nombre: aux,
			Calif:  map[string]float64{},
			Suma:   0,
		}
	}

	http.HandleFunc("/", Root)
	http.HandleFunc("/AgregarCalif", AgregarCalif)
	http.HandleFunc("/PromedioAlumno", PromedioAlumno)
	http.HandleFunc("/PromedioMateria", PromedioMateria)
	http.HandleFunc("/PromedioGeneral", PromedioGeneral)
	Println("Corriendo servidor...")
	http.ListenAndServe(":9000", nil)

}
