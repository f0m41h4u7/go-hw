package internal

//easyjson:json
type EventList []Event

//easyjson:json
type Event struct {
	UUID        string `db:"uuid" json:"uuid"`
	Title       string `db:"title" json:"title"`
	Start       string `db:"start" json:"start"`
	End         string `db:"end" json:"end"`
	Description string `db:"description" json:"description"`
	OwnerID     string `db:"ownerid" json:"ownerid"`
	NotifyIn    string `db:"notifyin" json:"notifyin"`
}
