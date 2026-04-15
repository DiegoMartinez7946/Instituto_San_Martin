package services

import (
	"log"
	"strconv"
	stdtime "time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
)

var diasES = []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
var mesesES = []string{
	"", "enero", "febrero", "marzo", "abril", "mayo", "junio",
	"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
}

/* BuildArgentinaTimeDisplay calcula fecha/hora en Argentina (America/Argentina/Buenos_Aires),
   equivalente a https://time.is/es/Argentina para uso en la app (sin scraping). */
func BuildArgentinaTimeDisplay() models.ArgentinaTimeDisplay {
	loc, err := stdtime.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		loc = stdtime.FixedZone("ART", -3*3600)
	}
	now := stdtime.Now().In(loc)
	diaNombre := diasES[int(now.Weekday())]
	mesNombre := mesesES[int(now.Month())]
	return models.ArgentinaTimeDisplay{
		Hora:  now.Format("15:04:05"),
		Fecha: now.Format("02/01/2006"),
		Dia:   diaNombre,
		Mes:   mesNombre,
		Ano:   strconv.Itoa(now.Year()),
	}
}

/* GetArgentinaClockService construye la hora argentina y la persiste en la colección time */
func GetArgentinaClockService() (models.ArgentinaTimeDisplay, bool, error) {
	row := BuildArgentinaTimeDisplay()
	ok, err := db.UpsertArgentinaTimeDisplay(row)
	if err != nil {
		log.Println("UpsertArgentinaTimeDisplay:", err)
	}
	return row, ok, nil
}
