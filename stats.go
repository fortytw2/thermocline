package thermocline

// Stats are the statistics unit that should be returned by a broker
type Stats struct {
	Failures map[string]int `json:"failures"`
	Retries  map[string]int `json:"retries"`
	Finished map[string]int `json:"finished"`
	Total    map[string]int `json:"total"`
}
