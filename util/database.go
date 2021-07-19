package util

import (
	"context"
	"fmt"
	"log"
	"time"

	//wfGame "github.com/goyohangoyoill/ONUWF/game"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoadDBInfo struct {
	MatchedUserList []*LoadedUser
	LastRoleSeq     []int //User로
}

type SaveDBInfo struct {
	CurUserList []*InputUser
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
dbName := "ONUWF"
colName := "User"

func GetCollection(client *mongo.Client, colName string) *mongdo.Collection {
	return client.Database(dbNmae).Collection(colName)
}
*/
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

func LoadEachUser(uid string, m bool, collection string, mongoDB *mongo.Database, ctx context.Context) (LoadedUser, bool) {
	//var result bson.M
	result := InputUser{}
	filter := bson.D{{"userid", uid}}
	err := mongoDB.Collection(collection).FindOne(ctx, filter).Decode(&result)
	fResult := LoadedUser{}
	if err == nil {
		fResult = LoadedUser{result.UserID, result.Nick, result.Title, result.LastRole}
		if m == false {
			fResult.LastRole = nil
		}
	} else {
		return fResult, false
	}
	return fResult, true
}

/*
func LoadUser(sDB SaveDBInfo, collection string, mongoDB *mongo.Database, ctx context.Context) LoadDBinfo {
	var result bson.M
	uLen = len(sDB.CurUserList)
	lUsers = make([]LoadedUser, uLen)
	lRole := nil
	for i := 0; i < uLen; i++ {
		filter := bson.D{"userid": sDB.CurUserList[i].UserID}
		err := mongoDB.Collection(collection).FindOne(ctx, filter).Decode(&result)
		if err != nil {
			fResult := LoadedUser{result.userid, result.nick, result.title, nil}
			if sDB.CurUserList[i].UserID == sDB.MUserID {
				lRole = result.lastrole
			}
			lUsers.append(fResult)
		}
	}
	return LoadDBinfo{lUsers, lRole}
}
*/

//func LoadUser

func SetStartUser(sDB SaveDBInfo, collection string, mongoDB *mongo.Database, ctx context.Context) string {
	uLen := len(sDB.CurUserList)
	t := time.Now().Format(time.ANSIC)
	for i := 0; i < uLen; i++ {
		Input := InputUser{}
		user := sDB.CurUserList[i]
		filter := bson.D{{"userid", user.UserID}}
		update := bson.D{}
		num, err := mongoDB.Collection(collection).CountDocuments(ctx, filter)
		CheckErr(err)
		//User 정보 없을 시 db에 유저등록
		if num == 0 {
			if user.UserID == sDB.MUserID {
				Input = InputUser{user.UserID, user.Nick, "BetaTester", t, 0, sDB.CurRoleSeq}
			} else {
				Input = InputUser{user.UserID, user.Nick, "BetaTester", t, 0, nil}
			}
			_, err := mongoDB.Collection(collection).InsertOne(ctx, Input)
			CheckErr(err)
		} else if num == 1 {
			//master user 일 경
			if user.UserID == sDB.MUserID {
				update = bson.D{{"$set", bson.D{{"nick", user.Nick}, {"lastplayeddate", t}, {"lastrole", sDB.CurRoleSeq}}}}
			} else {
				update = bson.D{{"$set", bson.D{{"nick", user.Nick}, {"lastplayeddate", t}}}}
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

// DB에 값이 존재하는지 확인

// 새로 넣을 데이터 정의
/*
	newData := bson.M{
		"UserID": user.UserID,
		"nick":   user.Nick(),
	}
*/
// DB값이 존재하지 않으면

//func UpdateUser

/*
func VoteData(collection string, mongoDB *mongo.Database, ctx contxt.Context) string {

}
*/
