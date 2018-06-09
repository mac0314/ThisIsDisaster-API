package controllers

import (
	"fmt"
	"net/smtp"

	"github.com/revel/revel"
)

type SMTP struct {
	*revel.Controller
}

func SendEmail(to []string, title string, body string) (bool, string) {
	var err bool = false
	var msg string
	var admin string = "rladudals02@ajou.ac.kr"
	var password string = "password"
	var server string = "smtp.gmail.com"

	// 메일서버 로그인 정보 설정
	auth := smtp.PlainAuth("", admin, password, server)

	from := admin

	message := []byte(title + "\r\n\r\n\r\n" + body)

	// 메일 보내기
	_err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	if _err != nil {
		fmt.Println(_err)
		err = true
		msg = "Fail"
	}

	return err, msg
}

func (c SMTP) Index() revel.Result {
	var code int = RESULT_CODE_SUCCESS
	var msg string = "Success"

	to := []string{"rladudals02@gmail.com"} // 복수 수신자 가능

	title := "Subject : test\r\n"
	body := "hello\r\n\r\n"

	_err, _msg := SendEmail(to, title, body)
	if _err {
		code = RESULT_CODE_FAILURE
		msg = _msg
	}

	// JSON response
	response := make(map[string]interface{})
	response["result_code"] = code
	response["result_msg"] = msg
	response["response_type"] = RESULT_TYPE_SMTP

	return c.RenderJSON(response)
}
