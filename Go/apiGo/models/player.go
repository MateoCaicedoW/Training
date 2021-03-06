package models

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	Aplus  = "a+"
	Aminus = "a-"
	Bplus  = "b+"
	Bminus = "b-"
	Cplus  = "c+"
	Cminus = "c-"
	Dplus  = "d+"
	Dminus = "d-"

	PhysicalCondition   = []string{Aplus, Aminus, Bplus, Bminus, Cplus, Cminus, Dplus, Dminus}
	GoalKeeper          = "goalkeeper"
	Defender            = "defender"
	CentralMidfielder   = "centralmidfielder"
	Forward             = "forward"
	FullBack            = "fullback"
	HalfBack            = "halfback"
	DefensiveMidfielder = "defensivemidfielder"
	AttackingMidfielder = "attackingmidfielder"
	CentreBack          = "centreback"
	Winger              = "winger"

	Positions = []string{GoalKeeper, Defender, CentralMidfielder, Forward, FullBack, HalfBack, DefensiveMidfielder, AttackingMidfielder, CentreBack, Winger}
)

type Player struct {
	ID                uuid.UUID `gorm:"primary_key; not null; type:uuid"`
	FirstName         string
	LastName          string
	Level             int64
	Age               int64
	Position          string
	PhysicalCondition string
	Teams             []Team `gorm:"many2many:player_teams;"`
}

func (p *Player) Validate() (response PlayerResponse) {

	letters := regexp.MustCompile("^[a-zA-Z]+$")
	name := strings.Replace(strings.ToLower(p.FirstName), " ", "", -1)
	lastName := strings.Replace(strings.ToLower(p.LastName), " ", "", -1)

	response.Status = http.StatusBadRequest

	if name == "" {
		response.Message = "FirstName can not be empty."
		return
	}

	if p.numbersAndCaracters(name, "FirstName").Message != "" {
		response.Message = p.numbersAndCaracters(name, "FirstName").Message
		return
	}

	if p.numbersAndCaracters(lastName, "LastName").Message != "" {
		response.Message = p.numbersAndCaracters(lastName, "LastName").Message
		return
	}

	if lastName == "" {
		response.Message = "LastName can not be empty."
		return
	}

	if p.Level == 0 {
		response.Message = "Level can not be empty."
		return
	}

	if p.Level < 1 || p.Level > 99 && letters.MatchString(strconv.Itoa(int(p.Level))) {
		response.Message = "Level must be a number."
		return
	}
	if p.Level < 1 || p.Level > 99 {
		response.Message = "Level must be between 1 and 99."
		return
	}

	if p.Age == 0 {
		response.Message = "Age can not be empty."
		return
	}

	if p.Age < 1 && letters.MatchString(strconv.Itoa(int(p.Age))) {
		response.Message = "Age must be a number."
		return
	}

	if p.Age < 0 {
		response.Message = "Age must be greater than 0."
		return
	}

	if p.validatePosition().Message != "" {
		response.Message = p.validatePosition().Message
		return
	}
	if p.validatePhysicalCondition().Message != "" {
		response.Message = p.validatePhysicalCondition().Message
		return
	}
	if p.validateTeams().Message != "" {
		response.Message = p.validateTeams().Message
		return
	}

	response.Message = ""
	return

}

func (p *Player) validatePhysicalCondition() (response PlayerResponse) {
	numbers := regexp.MustCompile("^[0-9]+$")
	response.Status = http.StatusBadRequest
	physicalCondition := strings.Trim(strings.ToLower(p.PhysicalCondition), " ")
	if physicalCondition == "" {
		response.Message = "PhysicalCondition can not be empty."
		return
	}
	if physicalCondition != "" && numbers.MatchString(physicalCondition) {
		response.Message = "PhysicalCondition can not be a number."
		return
	}

	for _, item := range PhysicalCondition {
		if item == physicalCondition {
			response.Message = ""
			return
		}
	}
	response.Message = "Insert a valid PhysicalCondition."
	return
}

func (p *Player) numbersAndCaracters(param string, field string) (response PlayerResponse) {
	numbers := regexp.MustCompile("^[0-9]+$")
	caracters := regexp.MustCompile("^[!-/:-@[-`{-~-$]+$")
	if param != "" {
		for _, item := range strings.Split(param, "") {
			if caracters.MatchString(item) {
				response.Message = field + " can not contains caracters."
				return
			}
			if numbers.MatchString(item) {
				response.Message = field + " can not be a number."
				return
			}
		}
	}

	response.Message = ""
	return
}

func (p *Player) validateTeams() (response PlayerResponse) {
	response.Status = http.StatusBadRequest
	// if len(p.Teams) == 0 {
	// 	response.Message = "Teams can not be empty."
	// 	return
	// }
	if len(p.Teams) > 2 {
		response.Message = "Teams can not be greater than 2."
		return
	}

	if len(p.Teams) > 1 && p.Teams[0].Name == p.Teams[1].Name {
		response.Message = "Teams can not be the same."
		return
	}
	// if len(p.Teams) > 1 && len(p.Teams[1].Name) == 0 || len(p.Teams[0].Name) == 0 {
	// 	response.Message = "Teams can not be empty."
	// 	return
	// }
	response.Message = ""
	return
}

func (p *Player) validatePosition() (response PlayerResponse) {
	numbers := regexp.MustCompile("^[0-9]+$")
	response.Status = http.StatusBadRequest
	position := strings.Replace(strings.ToLower(p.Position), " ", "", -1)
	if position == "" {
		response.Message = "Position can not be empty."
		return

	}
	if position != "" && numbers.MatchString(position) {
		response.Message = "Position can not be a number."
		return
	}

	if p.numbersAndCaracters(position, "Position").Message != "" {
		response.Message = p.numbersAndCaracters(position, "Position").Message
		return
	}

	for _, pos := range Positions {
		if pos == position {
			response.Message = ""

			return
		}
	}
	response.Message = "Insert a valid Position."
	return
}

type PlayerResponse struct {
	Status  int
	Data    ListPlayers
	Message string
}

type ListPlayers []Player
