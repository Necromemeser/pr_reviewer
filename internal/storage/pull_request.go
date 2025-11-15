package storage

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/models"
	"time"
)

var ErrPRExists = errors.New("PR id already exists")
var ErrPRNotFound = errors.New("PR id not found")
var ErrReviewerNotFound = errors.New("reviewer id not found")
var ErrPRMerged = errors.New("PR is merged")
var ErrNoCandidate = errors.New("no active replacement candidate in team")
var ErrNotAssigned = errors.New("reviewer is not assigned to this PR")

func (s *Storage) CreatePullRequest(pullRequestID, pullRequestName, authorID string) (*models.PullRequest, error) {
	// Проверяем, нет ли уже PR с таким ID
	var existing string
	err := s.DB.QueryRow(`
        SELECT pull_request_id 
        FROM pull_requests 
        WHERE pull_request_id = $1
    `, pullRequestID).Scan(&existing)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existing != "" {
		return nil, ErrPRExists
	}

	// Выбираем до двух ревьюверов с минимальным количеством назначенных ревью
	rows, err := s.DB.Query(`
		SELECT u.user_id
		FROM users u
		WHERE 
			u.team_name = (SELECT team_name FROM users WHERE user_id = $1)
			AND u.user_id != $1
			AND u.is_active = TRUE
		ORDER BY (
			SELECT COUNT(*) 
			FROM reviewers r 
			WHERE r.reviewer_id = u.user_id
		) ASC
		LIMIT 2
	`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := []string{}
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	createdAt := time.Now()

	_, err = s.DB.Exec(`
        INSERT INTO pull_requests(
            pull_request_id, pull_request_name, author_id, status, created_at
        )
        VALUES ($1, $2, $3, 'OPEN', $4)
    `, pullRequestID, pullRequestName, authorID, createdAt)
	if err != nil {
		return nil, err
	}

	// Добавляем ревьюверов в таблицу связи
	for _, reviewerID := range reviewers {
		_, err = s.DB.Exec(`
            INSERT INTO reviewers(pull_request_id, reviewer_id)
            VALUES ($1, $2)
        `, pullRequestID, reviewerID)
		if err != nil {
			return nil, err
		}
	}

	// Возвращаем модель
	return &models.PullRequest{
		PullRequestID:     pullRequestID,
		PullRequestName:   pullRequestName,
		AuthorID:          &authorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         createdAt,
		MergedAt:          nil,
	}, nil
}

func (s *Storage) MergePullRequest(pullRequestID string) (*models.PullRequest, error) {
	mergedAt := time.Now()
	res, err := s.DB.Exec(`
        UPDATE pull_requests
        SET status = 'MERGED', merged_at = $1
        WHERE pull_request_id = $2 AND status = 'OPEN'
    `, mergedAt, pullRequestID)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		var existingStatus string
		err := s.DB.QueryRow(`SELECT status, merged_at FROM pull_requests WHERE pull_request_id = $1`, pullRequestID).
			Scan(&existingStatus, &mergedAt)
		if err == sql.ErrNoRows {
			return nil, ErrPRNotFound
		}
		if err != nil {
			return nil, err
		}
	}

	var mergedPR models.PullRequest
	row := s.DB.QueryRow(`
        SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1
    `, pullRequestID)
	if err := row.Scan(&mergedPR.PullRequestID, &mergedPR.PullRequestName, &mergedPR.AuthorID,
		&mergedPR.Status, &mergedPR.CreatedAt, &mergedPR.MergedAt); err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`SELECT reviewer_id FROM reviewers WHERE pull_request_id = $1`, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := make([]string, 0, 2)
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, r)
	}
	mergedPR.AssignedReviewers = reviewers

	return &mergedPR, nil
}

func (s *Storage) ReasignReviewer(pullRequestID, oldReviewerID string) (*models.PullRequest, string, error) {
	var status string
	err := s.DB.QueryRow(`SELECT status FROM pull_requests WHERE pull_request_id = $1`, pullRequestID).Scan(&status)
	if err == sql.ErrNoRows {
		return nil, "", ErrPRNotFound
	}
	if err != nil {
		return nil, "", err
	}
	if status == "MERGED" {
		return nil, "", ErrPRMerged
	}

	var exists string
	err = s.DB.QueryRow(`SELECT reviewer_id FROM reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`,
		pullRequestID, oldReviewerID).Scan(&exists)
	if err == sql.ErrNoRows {
		return nil, "", ErrReviewerNotFound
	}
	if err != nil {
		return nil, "", err
	}

	var teamName string
	err = s.DB.QueryRow(`SELECT team_name FROM users WHERE user_id = $1`, oldReviewerID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, "", ErrUserNotFound
	}
	if err != nil {
		return nil, "", err
	}

	var pRAuthor string
	err = s.DB.QueryRow(`SELECT author_id FROM pull_requests WHERE pull_request_id = $1`, pullRequestID).Scan(&pRAuthor)
	if err == sql.ErrNoRows {
		return nil, "", ErrUserNotFound
	}
	if err != nil {
		return nil, "", err
	}

	var anotherReviewerID string
	var rows *sql.Rows
	err = s.DB.QueryRow(`
		SELECT reviewer_id FROM reviewers
		WHERE pull_request_id = $1 AND reviewer_id != $2`, pullRequestID, oldReviewerID).Scan(&anotherReviewerID)
	if err == sql.ErrNoRows {
		rows, err = s.DB.Query(`
			SELECT u.user_id
			FROM users u
			WHERE 
				u.team_name = $1
				AND u.user_id != $2
				AND u.user_id != $3
				AND u.is_active = TRUE
			ORDER BY (
				SELECT COUNT(*) 
				FROM reviewers r 
				WHERE r.reviewer_id = u.user_id
			) ASC
			LIMIT 1
		`, teamName, oldReviewerID, pRAuthor)

		if err != nil {
			return nil, "", err
		}
		defer rows.Close()
	} else {
		rows, err = s.DB.Query(`
			SELECT u.user_id
			FROM users u
			WHERE 
				u.team_name = $1
				AND u.user_id != $2
				AND u.user_id != $3
				AND u.user_id != $4
				AND u.is_active = TRUE
			ORDER BY (
				SELECT COUNT(*) 
				FROM reviewers r 
				WHERE r.reviewer_id = u.user_id
			) ASC
			LIMIT 1
		`, teamName, oldReviewerID, pRAuthor, anotherReviewerID)

		if err != nil {
			return nil, "", err
		}
		defer rows.Close()
	}
	if err != nil {
		return nil, "", err
	}

	var newReviewerID string
	if rows.Next() {
		if err := rows.Scan(&newReviewerID); err != nil {
			return nil, "", err
		}
	} else {
		return nil, "", ErrNoCandidate
	}

	_, err = s.DB.Exec(`
		UPDATE reviewers 
		SET reviewer_id = $1
		WHERE pull_request_id = $2 AND reviewer_id = $3
	`, newReviewerID, pullRequestID, oldReviewerID)
	if err != nil {
		return nil, "", err
	}

	row := s.DB.QueryRow(`SELECT * FROM pull_requests WHERE pull_request_id = $1`, pullRequestID)
	var pr models.PullRequest
	err = row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", ErrPRNotFound
		}
		return nil, "", err
	}

	rows, err = s.DB.Query(`SELECT reviewer_id FROM reviewers WHERE pull_request_id = $1`, pullRequestID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	reviewers := []string{}
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return nil, "", err
		}
		reviewers = append(reviewers, r)
	}

	pr.AssignedReviewers = reviewers

	return &pr, newReviewerID, nil
}
