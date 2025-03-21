package response

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(data interface{}, message ...string) APIResponse {
	resp := APIResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		resp.Message = message[0]
	}

	return resp
}

func Error(message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   message,
	}
}
