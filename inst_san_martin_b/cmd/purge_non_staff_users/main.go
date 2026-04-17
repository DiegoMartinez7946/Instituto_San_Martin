// One-off: elimina de la colección user los documentos cuyo rol no sea ADMINISTRADOR ni ADMINISTRATIVO.
// Ejecutar desde la raíz del backend con las mismas variables de entorno que el servidor (.env).
package main

import (
	"log"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/db"
)

func main() {
	if config.CheckConnection() == 0 {
		log.Fatal("sin conexión a la base de datos")
	}
	n, err := db.DeleteNonStaffUsersFromCollection()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Documentos eliminados en colección user: %d", n)
}
