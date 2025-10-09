package httphandlers

type uploadResp struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type deleteResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type errResp struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
