package controllers

import (
	"ThisIsDisaster-API/app/models"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

func makeClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "13.124.166.242:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func CreateMatchingRoom() string {
	client := makeClient()

	room := RandStringBytesMaskImprSrc(6)

	err := client.SAdd("room/list", room).Err()
	if err != nil {
		panic(err)
	}

	err = client.SAdd("room/list/available", room).Err()
	if err != nil {
		panic(err)
	}

	return room
}

func GetAllMatchingRoom() []string {
	client := makeClient()

	val, err := client.SMembers("room/list").Result()
	if err != nil {
		panic(err)
	}

	return val
}

func CheckAvailableRoom() bool {
	var result bool
	client := makeClient()

	num, err := client.SCard("room/list/available").Result()
	if err != nil {
		panic(err)
	}

	if num > 0 {
		result = true
	} else {
		result = false
	}

	return result
}

func GetAvailableRoom() string {
	client := makeClient()

	room, err := client.SRandMember("room/list/available").Result()
	if err != nil {
		panic(err)
	}

	return room
}

func Matching(user models.UserLocal) {
	var room string

	flag := CheckAvailableRoom()

	if flag {
		room = GetAvailableRoom()

	} else {
		room = CreateMatchingRoom()
	}
	JoinMatchingRoom(room, user)
}

func JoinMatchingRoom(room string, user models.UserLocal) {
	client := makeClient()

	err := client.Set("user/"+user.Email+"/room", room, 0).Err()
	if err != nil {
		panic(err)
	}

	err = client.SAdd("room/"+room, user.Email).Err()
	if err != nil {
		panic(err)
	}

	roomNum, err := client.SCard("room/" + room).Result()
	if err != nil {
		panic(err)
	}

	if roomNum >= 4 {
		client.SRem("room/list/available", room)
	}
}

func GetMatchingRoom(room string) []interface{} {
	client := makeClient()

	val, err := client.SMembers("room/" + room).Result()
	if err != nil {
		panic(err)
	}

	var userData []interface{}
	for _, userString := range val {
		fmt.Println(userString)
		var user models.UserLocal
		json.Unmarshal([]byte(userString), &user)

		userData = append(userData, user)
	}

	return userData
}

func GetMyMatchingRoom(user models.UserLocal) []interface{} {
	client := makeClient()

	room, err := client.Get("user/" + user.Email + "/room").Result()
	if err != nil {
		panic(err)
	}

	userData := GetMatchingRoom(room)

	return userData
}

func LeaveMatchingRoom(email string) {
	client := makeClient()

	room, err := client.Get("user/" + email + "/room").Result()
	if err != nil {
		panic(err)
	}

	client.Del("user/" + email + "/room")
	client.SRem("room/"+room, email)
	client.SAdd("room/list/available", room)

}

func ClearMatchingRoom(room string) {
	client := makeClient()

	users, err := client.SMembers("room/" + room).Result()
	if err != nil {
		panic(err)
	}

	for _, name := range users {
		client.Del("user/" + name + "/room")
	}

	client.SRem("room/list", room)
	client.SRem("room/list/available", room)
	client.Del("room/" + room)

	fmt.Println("ClearMatchingRoom")
}
