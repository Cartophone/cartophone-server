package owntone

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) PlayPlaylist(playlistID string) error {
	// Call OwnTone API to play the playlist
	return nil
}
