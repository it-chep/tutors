package lessons

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/lessons/delete_lesson"
)

type HandlerGroup struct {
	DeleteLesson *delete_lesson.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		DeleteLesson: delete_lesson.NewHandler(adminModule),
	}
}
