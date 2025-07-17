package internal

import (
	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"database/sql"

)

type App struct {
	logger *charmLog.Logger
	db     *sql.DB 
}

func NewApp(logger *charmLog.Logger, db *sql.DB) *App {
	return &App{
		logger: logger,
		db:     db,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	breedHandler := &BreedHandler{DB: a.db}

	//SELECT
	r.HandleFunc("/getAllBreeds", breedHandler.GetAllBreeds).Methods("GET")
	r.HandleFunc("/getBreedById/{id}", breedHandler.GetBreedById).Methods("GET")
	r.HandleFunc("/getBreedsByParams", breedHandler.GetBreedsByParams).Methods("GET")

	//INSERT
	r.HandleFunc("/addBreed", breedHandler.AddBreed).Methods("POST")		
	
	//UPDATE
	r.HandleFunc("/updateBreed/{id}", breedHandler.UpdateBreed).Methods("PATCH")

	//DELETE
	r.HandleFunc("/deleteBreed/{id}", breedHandler.DeleteBreed).Methods("DELETE")
}
