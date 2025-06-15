package webhandler

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/webui/view/pages"
)

func HandleDashboardPage(w http.ResponseWriter, r *http.Request) {
	// Zatím bez chybové zprávy
	pages.DashboardPage().Render(r.Context(), w)
}
