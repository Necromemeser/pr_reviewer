package storage

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/models"
)

var ErrTeamNotFound = errors.New("team not found")

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
