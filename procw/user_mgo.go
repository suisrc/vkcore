package procw

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserByMgo(ctx context.Context, cll *mongo.Collection) UserUpdate {
	return func(user *UserData, keys ...string) error {

		filter := bson.M{"id": user.Id}
		// 更新所有字段
		if len(keys) == 0 {
			// cll.ReplaceOne(ctx, filter, usr)
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": user})
			return err
		}

		// 更新username
		if keys[0] == "username" {
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"username": user.Username}})
			return err
		}

		if keys[0] == "email" {
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"email": user.Email}})
			return err
		}

		// 更新tokens
		if keys[0] == "tokens" {
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"tokens": user.Tokens}})
			return err
		}

		// 更新operate_data
		if keys[0] == "operate_data" {
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"operate_data": user.OperateData}})
			return err
		}

		// 更新operate_data, 但是只更新operate_data中的某个字段
		if strings.HasPrefix(keys[0], "operate_data.") {
			data := bson.M{}
			for _, k := range keys {
				if !strings.HasPrefix(k, "operate_data.") {
					return fmt.Errorf("unknown key: %v", keys)
				}
				optkey := k[len("operate_data."):]
				data[k] = user.OperateData[optkey]
			}
			_, err := cll.UpdateOne(ctx, filter, bson.M{"$set": data})
			return err
		}
		return fmt.Errorf("unknown key: %v", keys)
	}
}
