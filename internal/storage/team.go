package storage

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/models"
)

var ErrTeamNotFound = errors.New("team not found")
var ErrTeamExists = errors.New("team already exists")

func (s *Storage) GetTeam(teamName string) (*models.Team, error) {
	var name string
	err := s.DB.QueryRow(`SELECT team_name FROM teams WHERE team_name = $1`, teamName).Scan(&name)
	if err == sql.ErrNoRows {
		return nil, ErrTeamNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`
        SELECT user_id, username, is_active
        FROM users
        WHERE team_name = $1`,
		teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.TeamMember{}
	for rows.Next() {
		var u models.TeamMember
		err := rows.Scan(&u.UserID, &u.Username, &u.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	team := &models.Team{
		TeamName: teamName,
		Members:  members,
	}

	return team, nil
}

func (s *Storage) AddTeam(team models.Team) (*models.Team, error) {
	var existing string
	err := s.DB.QueryRow(`SELECT team_name FROM teams WHERE team_name = $1`, team.TeamName).Scan(&existing)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existing != "" {
		return nil, ErrTeamExists
	}

	_, err = s.DB.Exec(`INSERT INTO teams(team_name) VALUES($1)`, team.TeamName)
	if err != nil {
		return nil, err
	}

	for _, u := range team.Members {
		var exists string
		err := s.DB.QueryRow(`SELECT user_id FROM users WHERE user_id = $1`, u.UserID).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if exists != "" {
			_, err = s.DB.Exec(`UPDATE users SET username=$1, is_active=$2, team_name=$3 WHERE user_id=$4`,
				u.Username, u.IsActive, team.TeamName, u.UserID)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = s.DB.Exec(`INSERT INTO users(user_id, username, team_name, is_active) VALUES($1,$2,$3,$4)`,
				u.UserID, u.Username, team.TeamName, u.IsActive)
			if err != nil {
				return nil, err
			}
		}
	}

	rows, err := s.DB.Query(`SELECT user_id, username, is_active FROM users WHERE team_name = $1`, team.TeamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var u models.TeamMember
		err := rows.Scan(&u.UserID, &u.Username, &u.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	return &models.Team{
		TeamName: team.TeamName,
		Members:  members,
	}, nil
}
