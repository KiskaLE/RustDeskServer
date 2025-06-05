package computer

type RefreshComputerPayload struct {
	ComputerName string `json:"computerName"`
	IP           string `json:"ip"`
	OS           string `json:"os"`
	OSVersion    string `json:"osVersion"`
	RustDeskID   string `json:"rustDeskID"`
}
