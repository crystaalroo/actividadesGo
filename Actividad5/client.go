package main

import (
	"bufio"
	"fmt"
	. "fmt"
	"net/rpc"
	"os"
	"strings"
)

type NuevaC struct {
	Name, Materia string
	Calif         float64
}

func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	in := bufio.NewReader(os.Stdin)
	var op int64
	for {
		Println("1) Agregar calificación de una materia")
		Println("2) Mostrar el promedio de un Alumno")
		Println("3) Mostrar el promedio general")
		Println("4) Mostrar el promedio de una materia")
		Println("0) Exit")
		Scanln(&op)

		switch op {
		case 1:
			Print("Nombre del alumno: ")
			name, _ := in.ReadString('\n')
			name = strings.ReplaceAll(name, "\n", "")
			Print("Nombre de la materia: ")
			materia, _ := in.ReadString('\n')
			materia = strings.ReplaceAll(materia, "\n", "")
			var calif float64
			Print("Calificacion: ")
			Scanln(&calif)
			var result string
			nc := NuevaC{name, materia, calif}
			err = c.Call("Server.AgregarCalif", nc, &result)
			if err != nil {
				Println(err)
			} else {
				Println(result)
			}
		case 2:
			Print("Nombre del alumno: ")
			name, _ := in.ReadString('\n')
			name = strings.ReplaceAll(name, "\n", "")
			var result float64
			_ = c.Call("Server.PromedioAlumno", name, &result)
			Print("El promedio del alumno " + name + " es: ")
			Println(result)
		case 4:
			Print("Nombre de la materia: ")
			materia, _ := in.ReadString('\n')
			materia = strings.ReplaceAll(materia, "\n", "")
			var result float64
			_ = c.Call("Server.PromedioMateria", materia, &result)
			Print("El promedio de la materia " + materia + " es: ")
			Println(result)
		case 3:
			var result float64
			_ = c.Call("Server.PromedioGeneral", 0, &result)
			Print("El promedio general es: ")
			Println(result)
		case 0:
			return
		default:
			Println("Opcion no valida")
		}
		Println()
	}
}

func main() {
	client()
}
