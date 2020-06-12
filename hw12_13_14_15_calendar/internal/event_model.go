package internal

type Event struct {
	UUID        string `db:"uuid"`
	Title       string `db:"title"`
	Start       string `db:"start"`
	End         string `db:"end"`
	Description string `db:"description"`
	OwnerID     string `db:"ownerid"`
	NotifyIn    string `db:"notifyin"`
}
