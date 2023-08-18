package playlist

type Playlist struct {
	// Interval sets the time between switching views in a playlist. FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
	Interval string `json:"interval"`
	// The ordered list of items that the playlist will iterate over.
	Items []PlaylistItem `json:"items"`
	// Name of the playlist.
	Name string `json:"name"`
	// dummy value so thema allows a breaking change version.
	Xxx string `json:"xxx"`
}

type PlaylistItem struct {
	// Value depends on type and describes the playlist item.
	//  - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This  is not portable as the numerical identifier is non-deterministic between different instances.  Will be replaced by dashboard_by_uid in the future. (deprecated)  - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All  dashboards behind the tag will be added to the playlist.  - dashboard_by_uid: The value is the dashboard UID
	Value string `json:"value"`
	// Type of the item.
	Type PlaylistItemType `json:"type"`
}

type PlaylistItemType string

const (
	Dashboard_by_tag PlaylistItemType = "dashboard_by_tag"
	Dashboard_by_uid PlaylistItemType = "dashboard_by_uid"
)
