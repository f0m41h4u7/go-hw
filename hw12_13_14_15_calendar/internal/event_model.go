package internal

//easyjson:json
type EventList []Event

//easyjson:json
type Event struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	NotifyIn    string `json:"notify_in"`
}
