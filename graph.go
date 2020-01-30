package grabana

type AlertChannel struct {
	ID   uint   `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"Name"`
	Type string `json:"type"`
}
