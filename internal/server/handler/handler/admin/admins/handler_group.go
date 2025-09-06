package admins

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins/create_admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins/delete_admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins/get_admin_by_id"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins/get_admins"
)

type HandlerGroup struct {
	CreateAdmin  *create_admin.Handler
	DeleteAdmin  *delete_admin.Handler
	GetAdminByID *get_admin_by_id.Handler
	GetAdmins    *get_admins.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		CreateAdmin:  create_admin.NewHandler(adminModule),
		DeleteAdmin:  delete_admin.NewHandler(adminModule),
		GetAdminByID: get_admin_by_id.NewHandler(adminModule),
		GetAdmins:    get_admins.NewHandler(adminModule),
	}
}
