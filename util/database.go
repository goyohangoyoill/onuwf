package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
type LoadDBInfo struct {
	MatchedUserList []*LoadedUser
	LastRoleSeq     []int //User로
}
*/
type SaveDBInfo struct {
	CurUserList []*UserData
	CurRoleSeq  []int
	MUserID     string
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func MongoConn(env map[string]string) (client *mongo.Client, ctx context.Context) {
	// timeout 기반의 Context 생성
	ctx, _ = context.WithTimeout(context.Background(), time.Second*4)

	// Authetication 을 위한 Client Option 구성
	clientOptions := options.Client().ApplyURI(
		env["dbURI"]).SetAuth(
		options.Credential{
			AuthSource: "",
			Username:   env["dbUserName"],
			Password:   env["dbPassword"],
		},
	)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB Connection Success")
	return client, ctx
}

/*
type LoadedUser struct {
	UserID   string
	Nick     string
	Title    string
	LastRole []int
}

type InputUser struct {
	UserID         string `bson: "UserID"`
	Nick           string `bson: "Nick"`
	Title          string `bson: "Title"`
	LastPlayedDate string `bson: "LastPlayedDate"`
	Played         int    `bson: "Played"`
	LastRole       []int  `bson: "LastRole"`
}
*/
func LoadEachUser(uid string, m bool, collection string, mongoDB *mongo.Database, ctx context.Context) (UserData, bool) {
	result := UserData{}
	filter := bson.D{{"uid", uid}}
	err := mongoDB.Collection(collection).FindOne(ctx, filter).Decode(&result)
	if err == nil {
		if m == false {
			result.LastRoleList = nil
		}
	} else {
		return result, false
	}
	return result, true
}

func SetStartUser(sDB SaveDBInfo, collection string, mongoDB *mongo.Database, ctx context.Context) string {
	uLen := len(sDB.CurUserList)
	t := time.Now()
	for i := 0; i < uLen; i++ {
		Input := UserData{}
		user := sDB.CurUserList[i]
		filter := bson.D{{"uid", user.UID}}
		update := bson.D{}
		num, err := mongoDB.Collection(collection).CountDocuments(ctx, filter)
		CheckErr(err)
		//User 정보 없을 시 db에 유저등록
		if num == 0 {
			if user.UID == sDB.MUserID {
				Input = UserData{user.UID, user.Nick, "", t, 0, 0, sDB.CurRoleSeq, nil}
			} else {
				Input = UserData{user.UID, user.Nick, "", t, 0, 0, nil, nil}
			}
			_, err := mongoDB.Collection(collection).InsertOne(ctx, Input)
			CheckErr(err)
		} else if num == 1 {
			//master user 일 경
			if user.UID == sDB.MUserID {
				update = bson.D{{"$set", bson.D{{"nick", user.Nick}, {"recentgametime", t}, {"lastrolelist", sDB.CurRoleSeq}}}}
			} else {
				update = bson.D{{"$set", bson.D{{"nick", user.Nick}, {"recentgametime", t}}}}
			}
		} else {
			fmt.Println("UserDB Overlapped")
			return "Overlapped"
		}
		_, err = mongoDB.Collection(collection).UpdateOne(ctx, filter, update)
		CheckErr(err)
	}
	return "create!"
}

func SaveGame(sGame GameData, t time.Time, collection string, mongoDB *mongo.Database, ctx context.Context) string {
	//시간은 이전 시간을 db에서 가져오고 현재 시간을 입력으로 받아서 저장해야함
	//User collection에서 game start time 을 load (별도로 저장?)
	filter := bson.D{{"uid", sGame.UserList[0].UID}}
	result := UserData{}
	err := mongoDB.Collection("User").FindOne(ctx, filter).Decode(&result)
	CheckErr(err)
	sGame.StartTime = result.RecentGameTime // 대소문자 확인해야함
	sGame.EndTime = t
	var OID primitive.ObjectID
	ret, err := mongoDB.Collection(collection).InsertOne(ctx, sGame)
	CheckErr(err)
	mapstructure.Decode(ret.InsertedID, &OID) // interface assertion
	fmt.Println(ret)
	fmt.Println(OID.String())
	return OID.String()
}

func SaveEachUser(user *UserData, curGameOID string, win bool, t time.Time, collection string, mongoDB *mongo.Database, ctx context.Context) {
	filter := bson.D{{"uid", user.UID}}
	update := bson.D{}
	if win == true {
		//t := time.Now()
		update = bson.D{{"$set", bson.D{{"recentgametime", t}, {"cntplay", user.CntPlay + 1}, {"cntwin", user.CntWin + 1}, {"playedgameoid", append(user.PlayedGameOID, curGameOID)}}}}
	} else {
		update = bson.D{{"$set", bson.D{{"recentgametime", t}, {"cntplay", user.CntPlay + 1}, {"playedgameoid", append(user.PlayedGameOID, curGameOID)}}}}
	}
	_, err := mongoDB.Collection(collection).UpdateOne(ctx, filter, update)
	CheckErr(err)
}

func SetUserNick(user *UserData, nick string, mongoDB *mongo.Database, ctx context.Context) {
	filter := bson.D{{"uid", user.UID}}
	update := bson.D{{"$set", bson.D{{"nick", nick}}}}
	_, err := mongoDB.Collection("User").UpdateOne(ctx, filter, update)
	CheckErr(err)
}

// func SetUserTitle(user *UserData, title string, mongoDB *mongo.Database, ctx context.Context) {
// 	filter := bson.D{{"uid", user.UID}}
// 	update := bson.D{{"$set", bson.D{{"title", title}}}}
// 	_, err := mongoDB.Collection("User").UpdateOne(ctx, filter, update)
// 	CheckErr(err)
// }
