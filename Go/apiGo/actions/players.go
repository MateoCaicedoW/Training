package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mateo/apiGo/models"
)

func (handler handler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)

	var players []models.Player
	var response models.PlayerResponse

	players = listAllPlayers(handler, w, response)
	response.Status = http.StatusOK
	response.Data = players

	json, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (handler handler) ShowPlayer(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	var response models.PlayerResponse

	params := mux.Vars(r)
	idPlayer := params["id"]
	player, err := findPlayer(handler, idPlayer, w, response)
	if err != nil {
		return
	}

	response.Status = http.StatusAccepted
	response.Data = models.ListPlayers{player}
	json, _ := json.Marshal(response)
	w.WriteHeader(http.StatusAccepted)
	w.Write(json)

}

func (handler handler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	var response models.PlayerResponse
	params := mux.Vars(r)
	idPlayer := params["id"]
	player, err := findPlayer(handler, idPlayer, w, response)
	if err != nil {
		return
	}

	handler.db.Model(&player).Association("Teams").Clear()
	if result := handler.db.Delete(&player); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "Can not delete player"
		response.Status = http.StatusInternalServerError
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}

	response.Status = http.StatusOK
	players := listAllPlayers(handler, w, response)
	response.Data = players
	json, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (handler handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	var response models.PlayerResponse
	var newPlayer models.Player
	json.NewDecoder(r.Body).Decode(&newPlayer)

	if newPlayer.ID.String() == "00000000-0000-0000-0000-000000000000" {
		newPlayer.ID = uuid.New()
	}

	err := newPlayer.Validate()
	if err.Message != "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(err)
		w.Write(json)
		return
	}
	log.Println("aaas")
	teams, response := findTeamPlayer(handler, w, newPlayer)
	if response.Status != http.StatusOK {
		return
	}
	newPlayer.Teams = teams

	if result := handler.db.Create(&newPlayer); result.Error != nil {
		w.WriteHeader(http.StatusBadGateway)
		response.Message = result.Error.Error()

		response.Status = http.StatusBadGateway
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}
	response.Status = http.StatusOK
	response.Data = models.ListPlayers{newPlayer}
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)

	json, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (handler handler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	body, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	idPlayer := params["id"]
	var response models.PlayerResponse
	var newPlayer models.Player
	json.Unmarshal(body, &newPlayer)

	player, err1 := findPlayer(handler, idPlayer, w, response)
	if err1 != nil {
		return
	}

	err2 := newPlayer.Validate()
	if err2.Message != "" {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(err2)
		w.Write(json)
		return
	}

	teams, response := findTeamPlayer(handler, w, newPlayer)

	if response.Status != http.StatusOK {
		return
	}
	newPlayer.Teams = teams
	//player.Teams = teams
	if result := handler.db.Model(&player).Updates(newPlayer); result.Error != nil {
		w.WriteHeader(http.StatusBadGateway)
		response.Message = result.Error.Error()
		response.Status = http.StatusBadGateway
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}

	handler.db.Model(&player).Association("Teams").Replace(teams)

	response.Status = http.StatusOK
	json, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (handler handler) allTeams() (teams []models.Team) {
	handler.db.Find(&teams)
	return
}

func findPlayer(handler handler, idPlayer string, w http.ResponseWriter,
	response models.PlayerResponse) (player models.Player, err error) {
	handler.db.Preload("Teams").Find(&player, "id = ?", idPlayer)
	if player.FirstName == "" {
		w.WriteHeader(http.StatusNotFound)
		response.Message = "Player not found"
		response.Status = http.StatusNotFound
		json, _ := json.Marshal(response)
		w.Write(json)

		err = errors.New("player not found")
		return
	}

	return
}

func listAllPlayers(handler handler, w http.ResponseWriter,
	response models.PlayerResponse) (players []models.Player) {
	if result := handler.db.Preload("Teams").Find(&players); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = strings.ToTitle(result.Error.Error())
		response.Status = http.StatusInternalServerError
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}
	return
}

func findTeamPlayer(handler handler, w http.ResponseWriter,
	tempUpdate models.Player) (teams []models.Team, response models.PlayerResponse) {
	var count int
	for i := range tempUpdate.Teams {

		nameTeam := strings.Replace(strings.ToLower(tempUpdate.Teams[i].Name), " ", "", -1)
		var teamTemp models.Team

		handler.db.Raw("SELECT * FROM teams WHERE name = ?", nameTeam).Scan(&teamTemp)

		if len(teamTemp.Name) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = strings.ToTitle(nameTeam) + " not found"
			response.Status = http.StatusInternalServerError
			log.Println("vaciooooo")
			json, _ := json.Marshal(response)
			w.Write(json)
			return
		}
		allteams := handler.allTeams()

		for _, item := range allteams {
			if item.Name == nameTeam {
				if strings.Replace(strings.ToLower(item.Type), " ", "", -1) == models.Club {
					count++
				}
			}
		}
		if count == 2 {
			w.WriteHeader(http.StatusBadRequest)
			response.Message = "Teams must be a club and a national team"
			response.Status = http.StatusBadRequest
			json, _ := json.Marshal(response)
			w.Write(json)
			return
		}

		teams = append(teams, teamTemp)

	}
	response.Status = http.StatusOK
	return
}
