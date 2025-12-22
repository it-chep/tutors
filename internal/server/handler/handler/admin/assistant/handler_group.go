package assistant

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/add_available_tg"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/create_assistant"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/delete_assistant"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/delete_available_tg"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/get_assistant_by_id"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/assistant/get_assistants"
)

type HandlerGroup struct {
	CreateAssistant   *create_assistant.Handler
	DeleteAssistant   *delete_assistant.Handler
	GetAssistantByID  *get_assistant_by_id.Handler
	GetAssistants     *get_assistants.Handler
	AddAvailableTG    *add_available_tg.Handler
	DeleteAvailableTG *delete_available_tg.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		CreateAssistant:   create_assistant.NewHandler(adminModule),
		DeleteAssistant:   delete_assistant.NewHandler(adminModule),
		GetAssistantByID:  get_assistant_by_id.NewHandler(adminModule),
		GetAssistants:     get_assistants.NewHandler(adminModule),
		AddAvailableTG:    add_available_tg.NewHandler(adminModule),
		DeleteAvailableTG: delete_available_tg.NewHandler(adminModule),
	}
}
