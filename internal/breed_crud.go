package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

//SELECT race par ID
func (h *BreedHandler) GetBreedById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	stmt, err := h.DB.Prepare("SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight, hypoallergenic, apartment_friendly FROM breeds WHERE id = ?")
	if err != nil {
		http.Error(w, "Erreur de préparation de la requête", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var b models.Breed
	err = row.Scan(
		&b.ID,
		&b.Species,
		&b.PetSize,
		&b.Name,
		&b.AverageMaleAdultWeight,
		&b.AverageFemaleAdultWeight,
		&b.Hypoallergenic,
		&b.ApartmentFriendly,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Aucune race trouvée avec cet ID", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erreur lors de la récupération : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

//SELECT race par params
func (h *BreedHandler) GetBreedsByParams(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	species := queryParams.Get("species")
    petSize := queryParams.Get("pet_size")
    name := queryParams.Get("name")
    maleWeightMin := queryParams.Get("average_male_adult_weight_min")
    maleWeightMax := queryParams.Get("average_male_adult_weight_max")
    femaleWeightMin := queryParams.Get("average_female_adult_weight_min")
    femaleWeightMax := queryParams.Get("average_female_adult_weight_max")
    hypoallergenicParam := queryParams.Get("hypoallergenic")
    apartmentFriendlyParam := queryParams.Get("apartment_friendly")

    query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight, hypoallergenic, apartment_friendly FROM breeds"
    conditions := []string{}
    params := []interface{}{}

    if species != "" {
        conditions = append(conditions, "species = ?")
        params = append(params, species)
    }
    if petSize != "" {
        conditions = append(conditions, "pet_size = ?")
        params = append(params, petSize)
    }
    if name != "" {
        conditions = append(conditions, "name LIKE ?")
        params = append(params, "%"+name+"%")
    }
    if maleWeightMin != "" {
        conditions = append(conditions, "average_male_adult_weight >= ?")
        params = append(params, maleWeightMin)
    }
    if maleWeightMax != "" {
        conditions = append(conditions, "average_male_adult_weight <= ?")
        params = append(params, maleWeightMax)
    }
    if femaleWeightMin != "" {
        conditions = append(conditions, "average_female_adult_weight >= ?")
        params = append(params, femaleWeightMin)
    }
    if femaleWeightMax != "" {
        conditions = append(conditions, "average_female_adult_weight <= ?")
        params = append(params, femaleWeightMax)
    }
    if hypoallergenicParam != "" {
        hypo, err := strconv.ParseBool(hypoallergenicParam)
        if err == nil {
            conditions = append(conditions, "hypoallergenic = ?")
            params = append(params, hypo)
        }
    }
    if apartmentFriendlyParam != "" {
        apt, err := strconv.ParseBool(apartmentFriendlyParam)
        if err == nil {
            conditions = append(conditions, "apartment_friendly = ?")
            params = append(params, apt)
        }
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    rows, err := h.DB.Query(query, params...)
    if err != nil {
        http.Error(w, "Erreur DB: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var breeds []models.Breed
    for rows.Next() {
        var b models.Breed
        if err := rows.Scan(
            &b.ID,
            &b.Species,
            &b.PetSize,
            &b.Name,
            &b.AverageMaleAdultWeight,
            &b.AverageFemaleAdultWeight,
            &b.Hypoallergenic,
            &b.ApartmentFriendly,
        ); err != nil {
            http.Error(w, "Erreur scan: "+err.Error(), http.StatusInternalServerError)
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

	query := "UPDATE breeds SET "
	params := []interface{}{}
	setParts := []string{}

	if b.Species != "" {
		setParts = append(setParts, "species = ?")
		params = append(params, b.Species)
	}
	if b.PetSize != "" {
		setParts = append(setParts, "pet_size = ?")
		params = append(params, b.PetSize)
	}
	if b.Name != "" {
		setParts = append(setParts, "name = ?")
		params = append(params, b.Name)
	}
	if b.AverageMaleAdultWeight != 0 {
		setParts = append(setParts, "average_male_adult_weight = ?")
		params = append(params, b.AverageMaleAdultWeight)
	}
	if b.AverageFemaleAdultWeight != 0 {
		setParts = append(setParts, "average_female_adult_weight = ?")
		params = append(params, b.AverageFemaleAdultWeight)
	}
	if b.Hypoallergenic {
		setParts = append(setParts, "hypoallergenic = ?")
		params = append(params, b.Hypoallergenic)
	}
	if b.ApartmentFriendly {
		setParts = append(setParts, "apartment_friendly = ?")
		params = append(params, b.ApartmentFriendly)
	}
	if len(setParts) == 0 {
		http.Error(w, "Aucun champ à mettre à jour", http.StatusBadRequest)
		return
	}

	query += strings.Join(setParts, ", ") + " WHERE id = ?"
	params = append(params, id)

	stmt, err := h.DB.Prepare(query)
	if err != nil {
		http.Error(w, "Erreur préparation requête: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(params...)
	if err != nil {
		http.Error(w, "Erreur exécution PATCH: "+err.Error(), http.StatusInternalServerError)
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
