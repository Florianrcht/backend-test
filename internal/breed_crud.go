package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/models"
)

type BreedHandler struct {
	DB *sql.DB
}

//SELECT toutes les races
func (h *BreedHandler) GetAllBreeds(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT * FROM breeds")
	if err != nil {
		http.Error(w, "DB Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var breeds []models.Breed
	for rows.Next() {
		var b models.Breed
		if err := rows.Scan(&b.ID, &b.Species, &b.PetSize, &b.Name,
			&b.AverageMaleAdultWeight, &b.AverageFemaleAdultWeight,
			&b.Hypoallergenic, &b.ApartmentFriendly); err != nil {
			http.Error(w, "Scan Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		breeds = append(breeds, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

//INSERT une race
func (h *BreedHandler) AddBreed(w http.ResponseWriter, r *http.Request) {
	var b models.Breed

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	stmt, err := h.DB.Prepare(`
		INSERT INTO breeds (
			species, pet_size, name,
			average_male_adult_weight, average_female_adult_weight,
			hypoallergenic, apartment_friendly
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		http.Error(w, "Erreur de préparation de la requête", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		b.Species, b.PetSize, b.Name,
		b.AverageMaleAdultWeight, b.AverageFemaleAdultWeight,
		b.Hypoallergenic, b.ApartmentFriendly,
	)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	insertedID, _ := result.LastInsertId()
	b.ID = int(insertedID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

//UPDATE une race
func (h *BreedHandler) UpdateBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var b models.Breed
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	stmt, err := h.DB.Prepare(`
		UPDATE breeds SET
			species = ?,
			pet_size = ?,
			name = ?,
			average_male_adult_weight = ?,
			average_female_adult_weight = ?,
			hypoallergenic = ?,
			apartment_friendly = ?
		WHERE id = ?
	`)
	if err != nil {
		http.Error(w, "Erreur de préparation de la requête", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		b.Species, b.PetSize, b.Name,
		b.AverageMaleAdultWeight, b.AverageFemaleAdultWeight,
		b.Hypoallergenic, b.ApartmentFriendly,
		id,
	)
	if err != nil {
		http.Error(w, "Erreur lors de la modification : "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Aucune race trouvée avec cet ID", http.StatusNotFound)
		return
	}

	b.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

//DELETE une race
func (h *BreedHandler) DeleteBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	stmt, err := h.DB.Prepare("DELETE FROM breeds WHERE id = ?")
	if err != nil {
		http.Error(w, "Erreur de préparation de la requête", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression : "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Aucune race trouvée avec cet ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
