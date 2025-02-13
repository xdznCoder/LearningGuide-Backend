package api

import (
	"LearningGuide/user_api/forms"
	"LearningGuide/user_api/global"
	"LearningGuide/user_api/validator"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	lancet "github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"net/http"
	"strconv"
	"time"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, ans, err := cp.Generate()
	fmt.Printf("New Captcha Answer: %v\n", ans)
	if err != nil {
		zap.S().Errorw("[captcha] generate captcha fail", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码出错",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}

func SendEmail(c *gin.Context) {
	emailForm := forms.EmailForm{}
	err := c.ShouldBindJSON(&emailForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !lancet.IsEmail(emailForm.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效电子邮箱",
		})
		return
	}

	ctx := context.Background()

	code := strconv.Itoa(random.RandInt(10000, 99999))
	redisResult := global.RDB.SetNX(ctx, emailForm.Email, code, time.Minute*5)
	if !redisResult.Val() {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码已发送",
		})
		return
	}

	err = sendMessage(emailForm.Email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码发送失败",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "验证码发送成功",
	})
}

func sendMessage(email string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", global.ServerConfig.Email.Host)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify Your Email")
	m.SetBody("text/html", fmt.Sprintf(verificationEmailTemplate, code,
		time.Now().Format("2006-01-02 15:04:05")),
	)
	d := gomail.NewDialer("smtp.163.com", 465, global.ServerConfig.Email.Host, global.ServerConfig.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

const verificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>学海导航助手邮箱验证</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            background-color: #007bff;
            color: #fff;
            padding: 10px 0;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .content {
            padding: 20px 0;
        }
        .content p {
            line-height: 1.6;
            margin-bottom: 20px;
        }
        .content strong {
            display: block;
            margin-top: 20px;
            font-size: 24px;
            text-align: center;
            color: #007bff;
        }
        .footer {
            background-color: #333;
            color: #fff;
            text-align: center;
            padding: 10px 0;
            border-radius: 0 0 5px 5px;
        }
        .footer p {
            margin: 0;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>学海导航助手邮箱验证</h1>
        </div>
        <div class="content">
            <p>亲爱的用户：</p>
            <p>感谢您注册学海导航助手！为了确保您的账户安全，我们需要验证您的邮箱。请在下方输入您收到的验证码，以完成邮箱验证。</p>
            <strong>验证码：%s</strong>
            <p>如果您没有注册学海导航助手，请忽略此邮件。此验证码在 5 分钟内有效，请尽快完成验证。</p>
        </div>
        <div class="footer">
            <p>学海导航助手团队</p>
            <p>%s</p>
        </div>
    </div>
</body>
</html>`
