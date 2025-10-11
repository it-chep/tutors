package move_student

type Request struct {
	OldTutorID int64 `json:"old_tutor_id,omitempty"`
	NewTutorID int64 `json:"new_tutor_id"`
	StudentID  int64 `json:"student_id,omitempty"`
}
