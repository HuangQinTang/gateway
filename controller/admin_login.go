package controller

import (
	"encoding/json"
	"gateway/dao"
	"gateway/dto"
	"gateway/middleware"
	"gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)    //登陆
	group.GET("/logout", adminLogin.AdminLoginOut) //退出登陆
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.VerifyErrorCode, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, middleware.GROUPALL_SAVE_FLOWERROR, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, middleware.MysqlQueryErrorCode, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, middleware.JsonTransformErrorCode, err)
		return
	}
	sess := sessions.Default(c)                                    //获取一开始通过中间件设置到上下文到sessino配置
	sess.Options(sessions.Options{MaxAge: public.SessionInfoTime}) //session过期时间为三天，客户端，服务端一致
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	if err = sess.Save(); err != nil {
		middleware.ResponseError(c, middleware.SessionErrorCode, err)
	}

	out := dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
	return
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	if err := sess.Save(); err != nil {
		middleware.ResponseError(c, middleware.SessionErrorCode, err)
	}
	middleware.ResponseSuccess(c, "")
	return
}
