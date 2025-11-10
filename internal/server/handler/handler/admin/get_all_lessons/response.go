package get_all_lessons

type Lesson struct {
	ID                int64   `json:"id"`
	CreatedAt         string  `json:"created_at"`
	DurationInMinutes float64 `json:"duration_in_minutes"`

	StudentID   int64  `json:"student_id"`
	StudentName string `json:"student_name"`

	TutorID   int64  `json:"tutor_id"`
	TutorName string `json:"tutor_name"`
}

type Response struct {
	Lessons      []Lesson `json:"lessons"`
	LessonsCount int64    `json:"lessons_count"`
}
