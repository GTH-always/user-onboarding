package POST

import (
	"context"
	"net/http"
	"user-onboarding/constants"
	helpers "user-onboarding/helpers/userAction"
	structs "user-onboarding/struct"
	response "user-onboarding/struct/response"
	"user-onboarding/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	defer sentry.Recover()
	span := sentry.StartSpan(context.TODO(), "[GIN] UserDetails", sentry.TransactionName("Create a new user"))
	defer span.Finish()

	formRequest := structs.SignUp{}

	if err := c.ShouldBind(&formRequest); err != nil {
		span.Status = sentry.SpanStatusFailedPrecondition
		sentry.CaptureException(err)
		c.JSON(422, utils.SendErrorResponse(err))
		return
	}
	ctx := c.Request.Context()
	resp := response.SuccessResponse{}

	err := helpers.UserLogin(ctx, &formRequest, span.Context()) //the DAO level
	if err != nil {
		resp.Status = "Failed"
		resp.Message = "Incorred password"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] User login")
	token, tokenerror := utils.GenerateToken()
	dbSpan1.Finish()

	if tokenerror != nil {
		resp.Status = constants.API_FAILED_STATUS
		resp.Message = "Unable to login"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Status = "Success"
	resp.Message = "User data verified successfully"
	resp.Token = token
	span.Status = sentry.SpanStatusOK

	c.JSON(http.StatusOK, resp)

}
