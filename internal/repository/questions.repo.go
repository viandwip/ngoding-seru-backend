package repository

import (
	"math/rand"

	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
)

type RepoQuestionsIF interface {
	CreateQuestion(data *models.Question) (*config.Result, error)
	FetchQuiz(types string) (*config.Result, error)
	// SearchQuestions(search string, page, limit int) (*config.Result, error)
	// FetchQuestion(id, slug string) (*config.Result, error)
	// UpdateQuestion(data *models.Question) (*config.Result, error)
	// RemoveQuestion(id string) (*config.Result, error)
}

type RepoQuestions struct {
	*sqlx.DB
}

func NewQuestion(db *sqlx.DB) *RepoQuestions {
	return &RepoQuestions{db}
}

// Create Question
func (r *RepoQuestions) CreateQuestion(data *models.Question) (*config.Result, error) {
	q := `INSERT INTO questions (
				image, 
				type,
				level,
				question,
				option_a,
				option_b,
				option_c,
				option_d,
				explanation,
				answer
			)
			VALUES(
				:image, 
				:type,
				:level,
				:question,
				:option_a,
				:option_b,
				:option_c,
				:option_d,
				:explanation,
				:answer
			)`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data question created"}, nil
}

// Get Quiz
func (r *RepoQuestions) FetchQuiz(types string) (*config.Result, error) {
	var result []models.Question

	q := `SELECT * FROM questions WHERE type = $1 LIMIT 100`

	if err := r.Select(&result, q, types); err != nil {
		return nil, err
	}

	// Create a slice with numbers 1 to 100
	count := 4
	numbers := make([]int, count)
	for i := 0; i < count; i++ {
		numbers[i] = i + 1
	}

	// Shuffle the slice
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	return &config.Result{Data: result, Numbers: numbers}, nil
}

// // Search Questions
// func (r *RepoQuestions) SearchQuestions(search string, page, limit int) (*config.Result, error) {
// 	searchStr := fmt.Sprintf("%%%s%%", search)
// 	offset := (page - 1) * limit
// 	var result []models.Question

// 	q := `SELECT * FROM questions WHERE name ILIKE $1 ORDER BY name DESC LIMIT $2 OFFSET $3`

// 	if err := r.Select(&result, r.Rebind(q), searchStr, limit, offset); err != nil {
// 		return nil, err
// 	}

// 	// Pagination
// 	var totalCount int

// 	totalCountQuery := `SELECT COUNT(*) FROM questions WHERE name ILIKE $1`

// 	if err := r.Get(&totalCount, totalCountQuery, searchStr); err != nil {
// 		return nil, err
// 	}

// 	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
// 	var next, prev interface{}
// 	if page < totalPages {
// 		next = page + 1
// 	}
// 	if page > 1 {
// 		prev = page - 1
// 	}

// 	meta := &config.Metas{
// 		Total: totalCount,
// 		Next:  next,
// 		Prev:  prev,
// 	}

// 	return &config.Result{Data: result, Meta: meta}, nil
// }

// // Get Question
// func (r *RepoQuestions) FetchQuestion(id, slug string) (*config.Result, error) {
// 	var result models.Question
// 	data := id
// 	column := "id"
// 	if id == "" {
// 		data = slug
// 		column = "slug"
// 	}
// 	fmt.Println(data, column)

// 	q := fmt.Sprintf("SELECT * FROM questions WHERE %s = $1", column)

// 	if err := r.Get(&result, r.Rebind(q), data); err != nil {
// 		return nil, err
// 	}

// 	return &config.Result{Data: result}, nil
// }

// // Update Question
// func (r *RepoQuestions) UpdateQuestion(data *models.Question) (*config.Result, error) {
// 	q := `UPDATE questions
// 				SET
// 					image = COALESCE(NULLIF(:image, ''), image),
// 					name = COALESCE(NULLIF(:name, ''), name),
// 					description = COALESCE(NULLIF(:description, ''), description),
// 					stock = COALESCE(CAST(NULLIF(:stock, '') AS INT), stock),
// 					price = COALESCE(CAST(NULLIF(:price, '') AS INT), price),
// 					slug = COALESCE(NULLIF(:slug, ''), slug),
// 					size = COALESCE(NULLIF(:size, ''), size),
// 					updated_at = NOW()
// 				WHERE
// 					id = CAST(:id AS UUID)
// 			`

// 	_, err := r.NamedExec(q, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &config.Result{Message: "1 data question updated"}, nil
// }

// // Delete Question
// func (r *RepoQuestions) RemoveQuestion(id string) (*config.Result, error) {
// 	q := `DELETE FROM questions
// 				WHERE
// 					id = $1`

// 	_, err := r.Exec(q, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &config.Result{Message: "1 data question deleted"}, nil
// }
