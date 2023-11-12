package corew

import "time"

// 用户信息
type UserData struct {
	Username       string `json:"username,omitempty" bson:"username,omitempty"`
	DisplayName    string `json:"display_name,omitempty" bson:"display_name,omitempty"`
	Id             string `json:"id,omitempty" bson:"id,omitempty"`
	Xp             int    `json:"xp,omitempty" bson:"xp,omitempty"`
	Email          string `json:"email,omitempty" bson:"email,omitempty"`
	FollowerCount  int    `json:"follower_count,omitempty" bson:"follower_count,omitempty"`
	FollowingCount int    `json:"following_count,omitempty" bson:"following_count,omitempty"`
	CreatedAt      string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Tokens         struct {
		AccessToken  string    `json:"access_token,omitempty" bson:"access_token,omitempty"`
		RefreshToken string    `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
		Expiration   time.Time `json:"expiration,omitempty" bson:"expiration,omitempty"`
		StreamToken  string    `json:"stream_token,omitempty" bson:"stream_token,omitempty"`
		AmsS3Token   struct {
			AccessKey    string    `json:"access_key,omitempty" bson:"access_key,omitempty"`
			SecretKey    string    `json:"secret_key,omitempty" bson:"secret_key,omitempty"`
			SessionToken string    `json:"session_token,omitempty" bson:"session_token,omitempty"`
			IdentityId   string    `json:"identity_id,omitempty" bson:"identity_id,omitempty"`
			Expiration   time.Time `json:"expiration,omitempty" bson:"expiration,omitempty"`
		} `json:"ams_s3_token,omitempty" bson:"ams_s3_token,omitempty"`
	} `json:"tokens,omitempty" bson:"tokens,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"` // 大部分不存在，只有在注册时才会有

	Zonex string `json:"zonex,omitempty" bson:"zonex,omitempty"` // 用户所在的区域
	Group string `json:"group,omitempty" bson:"group,omitempty"` // 用户组, 用于区分不同的用户
	// 最后更新时间
	UpdatedAt time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`
	// 同步时间, 与官网同步时间
	SyncAt time.Time `json:"sync_at,omitempty" bson:"sync_at,omitempty"`
	// 操作数据, 用于记录用户的操作
	OperateData map[string]interface{} `json:"operate_data,omitempty" bson:"operate_data,omitempty"`
	OperateTemp map[string]interface{} `json:"-" bson:"-"`
	USN         int                    `json:"usn,omitempty" bson:"usn,omitempty"`
}

type UserUpdate func(user *UserData, key ...string) error
