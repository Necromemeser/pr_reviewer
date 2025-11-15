package storage

import (
	"errors"
	// "fmt"
	"database/sql"
	"pr_reviewer/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

func (s *Storage) GetPR(userID string) ([]models.PullRequestShort, error) {
	// Получаем все pull_request_id, где пользователь ревьювер
	rows, err := s.DB.Query(`
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
        FROM pull_requests pr
        JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
        WHERE r.reviewer_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := []models.PullRequestShort{}
	for rows.Next() {
		var pr models.PullRequestShort
		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (s *Storage) SetIsActive(userID string, isActive bool) (*models.User, error) {
	// Проверяем, что пользователь существует
	var name string
	err := s.DB.QueryRow(`SELECT username FROM users WHERE user_id = $1`, userID).Scan(&name)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = s.DB.Exec(`
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2`,
		isActive, userID)

	if err != nil {
		return nil, err
	}

	row := s.DB.QueryRow(`SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`, userID)

	var updatedUser models.User
	err = row.Scan(&updatedUser.UserID, &updatedUser.Username, &updatedUser.TeamName, &updatedUser.IsActive)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
