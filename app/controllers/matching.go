package controllers

import (
	"ThisIsDisaster-API/app/models"
	"fmt"

	"github.com/go-redis/redis"
)

type Matching struct {
	GorpController
}

func makeClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "13.124.166.242:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func (c Matching) SelectMatchingUsers(emails []string) []models.User {
	//	var msg string
	var query string
	fmt.Println(len(emails))
	if len(emails) > 0 {
		query = "SELECT * FROM user"
	}

	for i, email := range emails {
		if i == 0 {
			query += " WHERE email_mn = '" + email + "'"
		} else {
			query += " OR email_mn = '" + email + "'"
		}
	}

	var list []models.User
	_, _err := c.Txn.Select(&list, query)
	if _err != nil {
		fmt.Println(_err)
	}
	fmt.Println(query)

	return list
}

func (c Matching) UpdateIP(email string, ip string) (bool, string) {
	var err bool
	var msg string

	_, _err := c.Txn.Exec("UPDATE user SET ip_sn = ? WHERE email_mn = ?", ip, email)

	if _err != nil {
		err = true
	} else {
		err = false
	}

	return err, msg
}

func CreateMatchingRoom(user models.User) string {
	client := makeClient()

	room := RandStringBytesMaskImprSrc(6)

	err := client.SAdd("room/list", room).Err()
	if err != nil {
		fmt.Println(err)
	}

	err = client.SAdd("room/list/available", room).Err()
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set("room/"+room+"/host", user.Email, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	return room
}

func GetAllMatchingRoom() []string {
	client := makeClient()

	val, err := client.SMembers("room/list").Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}

func CheckAvailableRoom() bool {
	var result bool
	client := makeClient()

	num, err := client.SCard("room/list/available").Result()
	fmt.Println(err)

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
	fmt.Println(err)

	return room
}

func GoMatching(user models.User) {
	var room string

	flag := CheckAvailableRoom()

	if flag {
		room = GetAvailableRoom()

	} else {
		room = CreateMatchingRoom(user)
	}
	JoinMatchingRoom(room, user)
}

func JoinMatchingRoom(room string, user models.User) {
	client := makeClient()

	err := client.Set("user/"+user.Email+"/room", room, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	err = client.SAdd("room/"+room, user.Email).Err()
	if err != nil {
		fmt.Println(err)
	}

	roomNum, err := client.SCard("room/" + room).Result()
	if err != nil {
		fmt.Println(err)
	}

	if roomNum >= 4 {
		client.SRem("room/list/available", room)
	}
}

func (c Matching) GetMatchingRoom(room string) []models.User {
	client := makeClient()

	val, err := client.SMembers("room/" + room).Result()
	if err != nil {
		fmt.Println(err)
	}

	users := c.SelectMatchingUsers(val)

	return users
}

func (c Matching) GetMyMatchingRoom(email string) (string, []models.User) {
	client := makeClient()

	room, err := client.Get("user/" + email + "/room").Result()
	if err != nil {
		fmt.Println(err)
	}

	users := c.GetMatchingRoom(room)

	return room, users
}

func LoadHost(room string) string {
	client := makeClient()

	host, err := client.Get("room/" + room + "/host").Result()

	if err != nil {
		fmt.Println(err)
	}

	return host
}

func LeaveMatchingRoom(email string) {
	client := makeClient()

	room, err := client.Get("user/" + email + "/room").Result()
	if err != nil {
		fmt.Println(err)
	}

	host := LoadHost(room)

	users, err := client.SMembers("room/" + room).Result()
	if err != nil {
		fmt.Println(err)
	}
	if host == email {

		for _, value := range users {
			if value != email {
				err = client.Set("room/"+room+"/host", value, 0).Err()
				if err != nil {
					fmt.Println(err)
				}

				break
			}
		}

	}

	fmt.Println(len(users))

	if len(users) == 1 {
		client.Del("room/" + room + "/host")
		client.SRem("room/list", room)
		client.SRem("room/list/available", room)
	} else {
		client.SAdd("room/list/available", room)
	}
	client.Del("user/" + email + "/room")
	client.SRem("room/"+room, email)

}

func ClearMatchingRoom(room string) {
	client := makeClient()

	users, err := client.SMembers("room/" + room).Result()
	if err != nil {
		fmt.Println(err)
	}

	for _, name := range users {
		client.Del("user/" + name + "/room")
	}

	client.SRem("room/list", room)
	client.SRem("room/list/available", room)
	client.Del("room/" + room + "/host")
	client.Del("room/" + room)
}
