package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/auth0-community/auth0"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
)

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
}

type Organization struct {
	Name         string `json:"name"`
	CreationDate string `json:"creation_date"`
}

type Bootcamp struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Organization Organization
}

type Student struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Bootcamp Bootcamp
}

var (
	audience string
	domain   string
)

func main() {
	r := mux.NewRouter()

	// Obtener los valores de configuración de Auth0 del entorno
	audience = os.Getenv("AUTH0_AUDIENCE")
	domain = os.Getenv("AUTH0_DOMAIN")

	// Rutas de la API
	r.HandleFunc("/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/bootcamps", getBootcampsHandler).Methods("GET")
	r.HandleFunc("/bootcamps/{id}", updateBootcampHandler).Methods("PUT")

	// Inicializar el servidor
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Middleware para validar el token de acceso
func authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretProvider := auth0.NewKeyProvider(os.Getenv("AUTH0_CLIENT_SECRET"))
		secret, err := secretProvider.GetSecret()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		configuration := auth0.NewConfiguration(
			secretProvider,
			[]string{audience},
			domain,
			jose.HS256,
		)

		validator := auth0.NewValidator(configuration)
		token, err := validator.ValidateRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims := map[string]interface{}{}
		err = validator.Claims(token, secret, &claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}

// Manejador para obtener todos los usuarios
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Autenticar al usuario
	authorize(func(w http.ResponseWriter, r *http.Request) {
		// Obtener todos los usuarios
		users := []User{
			{ID: "1", Name: "John Doe", OrganizationID: "1", Email: "john@example.com"},
			{ID: "2", Name: "Jane Smith", OrganizationID: "1", Email: "jane@example.com"},
			{ID: "3", Name: "Bob Johnson", OrganizationID: "2", Email: "bob@example.com"},
		}

		// Enviar la lista de usuarios como respuesta
		json.NewEncoder(w).Encode(users)
	})(w, r)
}

// Manejador para obtener un usuario por ID
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Autenticar al usuario
	authorize(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el ID del usuario de la URL
		vars := mux.Vars(r)
		userID := vars["id"]

		// Obtener el usuario correspondiente al ID
		user := User{ID: userID, Name: "John Doe", OrganizationID: "1", Email: "john@example.com"}

		// Enviar el usuario como respuesta
		json.NewEncoder(w).Encode(user)
	})(w, r)
}

// Manejador para obtener todos los bootcamps
func getBootcampsHandler(w http.ResponseWriter, r *http.Request) {
	// Autenticar al usuario
	authorize(func(w http.ResponseWriter, r *http.Request) {
		// Obtener todos los bootcamps
		bootcamps := []Bootcamp{
			{ID: "1", Name: "Full-Stack Web Development", Organization: Organization{Name: "Acme Inc", CreationDate: "2022-01-01"}},
			{ID: "2", Name: "Data Science", Organization: Organization{Name: "Globex Corp", CreationDate: "2021-12-01"}},
		}

		// Enviar la lista de bootcamps como respuesta
		json.NewEncoder(w).Encode(bootcamps)
	})(w, r)
}

// Manejador para actualizar un bootcamp
func updateBootcampHandler(w http.ResponseWriter, r *http.Request) {
	// Autenticar al usuario
	authorize(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el ID del bootcamp de la URL
		vars := mux.Vars(r)
		bootcampID := vars["id"]

		// Obtener los datos del bootcamp a actualizar del cuerpo de la solicitud
		var bootcamp Bootcamp
		err := json.NewDecoder(r.Body).Decode(&bootcamp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Actualizar el bootcamp correspondiente al ID
		updatedBootcamp := Bootcamp{ID: bootcampID, Name: bootcamp.Name, Organization: bootcamp.Organization}

		// Enviar el bootcamp actualizado como respuesta
		json.NewEncoder(w).Encode(updatedBootcamp)
	})(w, r)
}
