package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
)

type RepoStatsIF interface {
	FetchStat(user_id string) (*config.Result, error)
	CreateUpdateStat(data *models.TypeStat) (*config.Result, error)
	// UpdateStat(data *models.Stat) (*config.Result, error)
	// FetchQuiz(types string) (*config.Result, error)
	// SearchStats(search string, page, limit int) (*config.Result, error)
	// FetchStat(id, slug string) (*config.Result, error)
	// UpdateStat(data *models.Stat) (*config.Result, error)
	// RemoveStat(id string) (*config.Result, error)
}

type RepoStats struct {
	*sqlx.DB
}

func NewStat(db *sqlx.DB) *RepoStats {
	return &RepoStats{db}
}

// Get Stat
func (r *RepoStats) FetchStat(user_id string) (*config.Result, error) {
	var stat models.TotalStat
	var err error

	getScore := `
		WITH RankedScores AS (
				SELECT *, RANK() OVER (ORDER BY total_score DESC, created_at ASC) AS rank
				FROM total_stats
		),
		UserRank AS (
				SELECT user_id, total_score, highest_score, rank
				FROM RankedScores
				WHERE user_id = CAST($1 AS UUID)
		),
		Count AS (
				SELECT COUNT(*) AS count
				FROM RankedScores
		)
		SELECT 
				ur.total_score,
				ur.highest_score,
				ur.rank,
				c.count
		FROM UserRank ur, Count c;
	`

	err = r.QueryRowx(getScore, user_id).StructScan(&stat)
	if err != nil {
		return nil, err
	}

	return &config.Result{Data: stat}, nil
}

// Create or Update Stat
func (r *RepoStats) CreateUpdateStat(data *models.TypeStat) (*config.Result, error) {
	var err error

	insertQuery := `
		INSERT INTO type_stats (
			user_id,
			type,
			easy_correct,
			easy_incorrect,
			medium_correct,
			medium_incorrect,
			hard_correct,
			hard_incorrect
		)
		VALUES(
			CAST($1 AS UUID),
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)`

	updateQuery := `
		UPDATE type_stats
		SET
			easy_correct = $1 + easy_correct,
			easy_incorrect = $2 + easy_incorrect,
			medium_correct = $3 + medium_correct,
			medium_incorrect = $4 + medium_incorrect,
			hard_correct = $5 + hard_correct,
			hard_incorrect = $6 + hard_incorrect,
			updated_at = NOW()
		WHERE
			user_id = CAST($7 AS UUID) AND type = $8
		`

	// Check if a record exists in type_stats
	var exists bool
	err = r.QueryRow(`SELECT EXISTS (SELECT 1 FROM type_stats WHERE user_id = CAST($1 AS UUID) AND type = $2)`, data.User_id, data.Type).Scan(&exists)
	if err != nil {
		return nil, err
	}

	// Insert or Update in type_stats
	if !exists {
		_, err = r.Exec(insertQuery, data.User_id, data.Type, data.Easy_correct, data.Easy_incorrect, data.Medium_correct, data.Medium_incorrect, data.Hard_correct, data.Hard_incorrect)
	} else {
		_, err = r.Exec(updateQuery, data.Easy_correct, data.Easy_incorrect, data.Medium_correct, data.Medium_incorrect, data.Hard_correct, data.Hard_incorrect, data.User_id, data.Type)
	}
	if err != nil {
		return nil, err
	}

	// Fetch Total Score
	var result *config.Result
	result, err = r.FetchStat(data.User_id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CREATE OR REPLACE FUNCTION update_type_stats()
// RETURNS TRIGGER AS $$
// BEGIN
//   -- Calculate the 'score' based on various conditions
//   NEW.score := NEW.easy_correct + NEW.medium_correct * 2 + NEW.hard_correct * 3 - (NEW.easy_incorrect * 3 + NEW.medium_incorrect * 2 + NEW.hard_incorrect);

//   IF NEW.score < 0 THEN
//   	NEW.score = 0;
//   END IF;

//   -- Update 'highest_score' if the new 'score' is higher
//   IF NEW.score > COALESCE((SELECT highest_score FROM type_stats WHERE user_id = NEW.user_id), 0) THEN
//     NEW.highest_score := NEW.score;
//   ELSE
//     NEW.highest_score := COALESCE(OLD.highest_score, 0);
//   END IF;

//   RETURN NEW;
// END;
// $$ LANGUAGE plpgsql;

// CREATE OR REPLACE FUNCTION update_total_stats()
// RETURNS TRIGGER AS $$
// BEGIN

// 	UPDATE total_stats
// 	SET
// 	  total_score = (
// 	    SELECT SUM(score)
// 	    FROM type_stats
// 	    WHERE user_id = NEW.user_id
// 	  ),
// 	  highest_score = (
// 	    SELECT GREATEST(
// 	      (SELECT highest_score FROM total_stats WHERE user_id = NEW.user_id),
// 	      (SELECT SUM(score) FROM type_stats WHERE user_id = NEW.user_id)
// 	    )
// 	  ),
// 	  updated_at = NOW()
// 	WHERE user_id = NEW.user_id;

//   RETURN NEW;
// END;
// $$ LANGUAGE plpgsql;

// CREATE OR REPLACE FUNCTION insert_total_stats()
// RETURNS TRIGGER AS $$
// BEGIN
// 	INSERT INTO total_stats (user_id, total_score, highest_score)
// 		VALUES (
// 			NEW.id, 0, 0
// 	);
// RETURN NEW;
// END
// $$ LANGUAGE plpgsql;

// -- Drop existing triggers if they exist
// DROP TRIGGER IF EXISTS update_type_stats_trigger ON type_stats;
// DROP TRIGGER IF EXISTS update_total_stats_trigger ON type_stats;
// DROP TRIGGER IF EXISTS insert_total_stats_trigger ON users;

// -- Create trigger for updating 'type_stats'
// CREATE TRIGGER update_type_stats_trigger
// BEFORE INSERT OR UPDATE ON type_stats
// FOR EACH ROW
// EXECUTE FUNCTION update_type_stats();

// -- Create trigger for updating 'total_stats'
// CREATE TRIGGER update_total_stats_trigger
// AFTER INSERT OR UPDATE ON type_stats
// FOR EACH ROW
// EXECUTE FUNCTION update_total_stats();

// -- Create trigger for updating 'total_stats'
// CREATE TRIGGER insert_total_stats_trigger
// AFTER INSERT OR UPDATE ON users
// FOR EACH ROW
// EXECUTE FUNCTION insert_total_stats();
