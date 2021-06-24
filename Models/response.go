package models

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}
