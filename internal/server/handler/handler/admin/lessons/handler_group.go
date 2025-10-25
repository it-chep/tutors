package lessons

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/lessons/update_lesson"
)

type HandlerGroup struct {
	UpdateLesson *update_lesson.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		UpdateLesson: update_lesson.NewHandler(adminModule),
	}
}
