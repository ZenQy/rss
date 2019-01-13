package main

type routeManual struct {
	Path string            `json:"path,omitempty"`
	Info map[string]string `json:"info,omitempty"`
}
