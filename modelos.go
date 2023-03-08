package main

import "time"

// Definir la estructura Alumno
type Alumno struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

// Definir la estructura Comentario
type Comentario struct {
	Texto    string    `json:"texto"`
	Creador  string    `json:"creador"`
	Fecha    time.Time `json:"fecha"`
	AlumnoID string    `json:"alumno_id"`
}

// Definir la estructura AlumnoComentarios para almacenar los comentarios de un alumno
type AlumnoComentarios struct {
	Alumno      Alumno       `json:"alumno"`
	Comentarios []Comentario `json:"comentarios"`
}

// Definir la estructura ListaAlumnos para almacenar una lista de alumnos
type ListaAlumnos struct {
	Alumnos []Alumno `json:"alumnos"`
}
