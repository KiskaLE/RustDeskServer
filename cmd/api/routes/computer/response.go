package computer

import "github.com/KiskaLE/RustDeskServer/cmd/api/database"

type GetComputersResponse struct {
	Computers []database.Computers `json:"computers"`
}
