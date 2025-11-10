package conduct_lesson

type Request struct {
	StudentID int64  `json:"student_id"`
	Duration  int64  `json:"duration"`
	Date      string `json:"date"`
}
