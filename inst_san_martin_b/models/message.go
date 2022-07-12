// Package models provides ...
package models

/* Response model for the response */
type Response struct {
	Message string `json:"message,omitempty"`
	Code int `json:"code,omitempty"`
	Ok bool `json:"ok,omitempty"`
	Data interface{} `json:"data"`
}

