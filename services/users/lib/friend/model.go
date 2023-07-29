package friend

type Friend struct {
	UserID   string `bson:"user_id"`
	FriendID string `bson:"friend_id"`
	Approved bool
}
