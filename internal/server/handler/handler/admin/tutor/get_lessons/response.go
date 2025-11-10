package get_lessons

type Lesson struct {
	Id              int64  `json:"id"`
	StudentId       int64  `json:"student_id"`
	TutorId         int64  `json:"tutor_id"`
	StudentFullName string `json:"student_full_name"`
	Date            string `json:"date"`
	DurationMinutes int64  `json:"duration_minutes"`
}

type Response struct {
	Lessons      []Lesson `json:"lessons"`
	LessonsCount int64    `json:"lessons_count"`
}
