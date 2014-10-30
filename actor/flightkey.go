package actor

// flight lock key for user/player
type User struct {
	Uid int64
}

// flight lock key for kingdom map tile
type Tile struct {
	Geohash int
}
