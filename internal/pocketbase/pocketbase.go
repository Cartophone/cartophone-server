package pocketbase

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) UIDExists(uid string) bool {
	// Query PocketBase API for UID existence
	return false
}

func (c *Client) GetPlaylistForUID(uid string) (string, bool) {
	// Query PocketBase API for playlist associated with UID
	return "", false
}

func (c *Client) RegisterUID(uid string) {
	// Register new UID in PocketBase
}

func (c *Client) UpdateUID(uid string) {
	// Update existing UID in PocketBase
}

func (c *Client) GetActiveAlarms() []Alarm {
	// Fetch active alarms from PocketBase
	return []Alarm{}
}

type Alarm struct {
	Hour       int
	Minute     int
	PlaylistID string
}
