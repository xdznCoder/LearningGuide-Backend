package api

import (
	"LearningGuide/user_api/forms"
	"LearningGuide/user_api/global"
	proto "LearningGuide/user_api/proto/userProto"
	"LearningGuide/user_api/utils"
	"LearningGuide/user_api/utils/JwtUtil"
	"LearningGuide/user_api/validator"
	"context"
	"errors"
	lancet "github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

func GetUserList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "10"))

	ctx := contextWithSpan(c)

	resp, err := global.UserSrvClient.GetUserList(ctx, &proto.PageInfo{
		PageNum:  uint32(pageNum),
		PageSize: uint32(pageSize),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]map[string]any, 0)
	for _, value := range resp.Data {
		result = append(result, getUserDTO(value))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": resp.Total,
	})
}

func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		validator.HandleValidatorError(err, c)
		return
	}
	if registerForm.Password != registerForm.RePassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "两次密码不一致",
		})
		return
	}
	if !lancet.IsEmail(registerForm.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效电子邮箱",
		})
		return
	}

	ctx := contextWithSpan(c)

	rdbResult := global.RDB.Get(ctx, registerForm.Email)
	if errors.Is(rdbResult.Err(), redis.Nil) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码已过期",
		})
		return
	}
	if rdbResult.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		return
	}
	if rdbResult.Val() != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误，请重试",
		})
		return
	}

	global.RDB.Del(ctx, registerForm.Email)
	result, err := global.UserSrvClient.CreateUser(ctx, &proto.CreateUserInfo{
		NickName: registerForm.Nickname,
		Password: registerForm.Password,
		Email:    registerForm.Email,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": result.Id,
	})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	if !checkIfAuthorized(id, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权更改",
		})
		return
	}

	updateForm := forms.UpdateUserForm{}
	err = c.ShouldBindJSON(&updateForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	ctx := contextWithSpan(c)
	birthday, err := time.Parse("2006-01-02", updateForm.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请正确填写生日",
		})
		return
	}

	if len(updateForm.Nickname) > 20 || len(updateForm.Nickname) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "昵称长度应大于0小于20",
		})
		return
	}

	if !lancet.IsUrl(updateForm.Image) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "image应为图片网址",
		})
		return
	}

	_, err = global.UserSrvClient.UpdateUser(ctx, &proto.UpdateUserInfo{
		Id:       int32(id),
		NickName: updateForm.Nickname,
		Gender:   updateForm.Gender,
		BirthDay: uint64(birthday.Unix()),
		Image:    updateForm.Image,
		Desc:     updateForm.Desc,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	if !checkIfAuthorized(id, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权更改",
		})
		return
	}

	ctx := contextWithSpan(c)

	_, err = global.UserSrvClient.DeleteUser(ctx, &proto.DeleteUserRequest{Id: int32(id)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func PasswordLogin(c *gin.Context) {
	passwordForm := forms.PasswordLoginForm{}
	err := c.ShouldBindJSON(&passwordForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !store.Verify(passwordForm.CaptchaId, passwordForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	ctx := contextWithSpan(c)

	userInfoResp, err := global.UserSrvClient.GetUserByEmail(ctx, &proto.EmailRequest{
		Email: passwordForm.Email,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "电子邮箱 " + passwordForm.Email + " 尚未注册",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登陆失败",
				})
				zap.S().Infof("用户登陆失败: %v", err)
			}
		}
		return
	}

	if passwordValid, _ := global.UserSrvClient.CheckPasswordInfo(ctx, &proto.PasswordCheck{
		Password:          passwordForm.Password,
		EncryptedPassword: userInfoResp.GetPassword(),
	}); passwordValid == nil || !passwordValid.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
		})
		return
	}

	token, err := JwtUtil.CreateJWT(uint(userInfoResp.Id), userInfoResp.NickName, uint(userInfoResp.Role))
	if err != nil {
		zap.S().Debugf("生成JWT失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器出错",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"msg":   "登陆成功",
	})
}

func ChangePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	if !checkIfAuthorized(id, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权更改",
		})
		return
	}

	ctx := contextWithSpan(c)

	passwordForm := forms.ChangePasswordForm{}
	err = c.ShouldBindJSON(&passwordForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if passwordForm.Password != passwordForm.RePassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "两次密码不一致",
		})
		return
	}

	_, err = global.UserSrvClient.ChangePassword(ctx, &proto.ChangePasswordRequest{
		Id:          int32(id),
		OldPassword: passwordForm.OldPassword,
		NewPassword: passwordForm.Password,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更改成功",
	})
}

func GetUserDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	ctx := contextWithSpan(c)

	resp, err := global.UserSrvClient.GetUserById(ctx, &proto.IdRequest{Id: int32(id)})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, getUserDTO(resp))
}

func GetUserByToken(c *gin.Context) {
	id, _ := c.Get("userId")

	ctx := contextWithSpan(c)

	resp, err := global.UserSrvClient.GetUserById(ctx, &proto.IdRequest{Id: int32(id.(uint))})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, getUserDTO(resp))
}

func getUserDTO(resp *proto.UserInfoResponse) map[string]any {
	return map[string]any{
		"id":       resp.Id,
		"email":    resp.Email,
		"nickname": resp.NickName,
		"gender":   resp.Gender,
		"birthday": resp.BirthDay,
		"image":    resp.Image,
		"desc":     resp.Desc,
	}
}

func checkIfAuthorized(id int, c *gin.Context) bool {
	user, _ := c.Get("claims")
	userId := user.(*JwtUtil.CustomClaims).ID
	role := user.(*JwtUtil.CustomClaims).AuthorityId
	if uint(id) == userId {
		return true
	}
	if role != 1 {
		return true
	}

	return false
}

func contextWithSpan(c *gin.Context) context.Context {
	ctx := context.Background()
	span, ok := c.Get("span")

	if !ok {
		zap.S().Info(c.Request.URL.Path + "no tracer injected")
		return ctx
	}

	ctx = opentracing.ContextWithSpan(ctx, span.(opentracing.Span))
	return ctx
}
