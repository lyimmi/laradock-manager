package main

// Project represents project.json
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Version    string `json:"version"`
	Binaryname string `json:"binaryname"`
	Frontend   struct {
		Dir     string `json:"dir"`
		Install string `json:"install"`
		Build   string `json:"build"`
		Bridge  string `json:"bridge"`
		Serve   string `json:"serve"`
	} `json:"frontend"`
	WailsVersion string `json:"WailsVersion"`
}
