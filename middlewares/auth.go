package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Check if session doesn't exist
	userid := session.Get("user")
	if userid == nil {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	ctx.Set("userid", userid)
}
