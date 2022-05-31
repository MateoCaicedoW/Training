package actions

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mateo/apiGo/models"
)

func ListPlayers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(models.Template{Status: 200, Data: models.P, Message: ""})
	if err != nil {
		json.NewEncoder(w).Encode(models.Template{Status: 400, Data: models.ListPlayers{}, Message: "Bad Request"})
	}

}

func Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	key := r.FormValue("id")

	players := models.P
	//log.Println(params["id"])
	for i := 0; i < len(players); i++ {
		if key == players[i].ID.String() {
			json.NewEncoder(w).Encode(models.Template{Status: 200, Data: models.ListPlayers{players[i]}, Message: ""})
			return
		}
	}
	json.NewEncoder(w).Encode(&models.Template{Status: 400, Data: models.ListPlayers{}, Message: "Player was not found"})

}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	key := r.FormValue("id")
	for i, item := range models.P {
		if key == item.ID.String() {
			models.P = append(models.P[:i], models.P[i+1:]...)
			json.NewEncoder(w).Encode(models.Template{Status: 200, Data: models.P, Message: ""})
			return
		}
	}
	json.NewEncoder(w).Encode(&models.Template{Status: 400, Data: models.ListPlayers{}, Message: "ID is invalid"})

}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	players := models.P
	var newPlayers models.Player
	json.NewDecoder(r.Body).Decode(&newPlayers)
	newPlayers.ID = uuid.New()
	if int(newPlayers.Level) == 0 || newPlayers.FirstName == "" || newPlayers.LastName == "" {
		json.NewEncoder(w).Encode(models.Template{Status: 400, Data: models.ListPlayers{}, Message: "Insert a valid Data"})
	} else {
		models.P = append(players, newPlayers)
		json.NewEncoder(w).Encode(models.Template{Status: 200, Data: models.ListPlayers{newPlayers}, Message: ""})

	}

}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	key := r.FormValue("id")

	var tempUpdate models.Player
	json.NewDecoder(r.Body).Decode(&tempUpdate)

	for i, item := range models.P {
		if key == item.ID.String() {
			if int(tempUpdate.Level) == 0 || tempUpdate.FirstName == "" || tempUpdate.LastName == "" {
				json.NewEncoder(w).Encode(models.Template{Status: 400, Data: models.ListPlayers{models.P[i]}, Message: "Insert a valid Data"})
				return
			} else {
				models.P[i].FirstName = tempUpdate.FirstName
				models.P[i].LastName = tempUpdate.LastName
				models.P[i].Level = tempUpdate.Level
				json.NewEncoder(w).Encode(models.Template{Status: 200, Data: models.ListPlayers{models.P[i]}, Message: ""})
				return
			}
		}
	}
	json.NewEncoder(w).Encode(&models.Template{Status: 400, Data: models.ListPlayers{}, Message: "ID is invalid"})

}
