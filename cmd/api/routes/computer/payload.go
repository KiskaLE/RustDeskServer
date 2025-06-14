package computer

type RefreshComputerPayload struct {
	ComputerName string `json:"computer_name"`
	IP           string `json:"ip"`
	OS           string `json:"os"`
	OSVersion    string `json:"os_version"`
	RustDeskID   string `json:"rustdesk_id"`
}
