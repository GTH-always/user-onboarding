package POST

import (
	"context"
	"net/http"
	helpers "user-onboarding/helpers/userAction"
	structs "user-onboarding/struct"
	response "user-onboarding/struct/response"
	"user-onboarding/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func FetchUser(c *gin.Context) {
	defer sentry.Recover()
	span := sentry.StartSpan(context.TODO(), "[GIN] FetchUser", sentry.TransactionName("Fetch a new user"))
	defer span.Finish()

	formRequest := structs.UserDetails{}

	if err := c.ShouldBind(&formRequest); err != nil {
		span.Status = sentry.SpanStatusFailedPrecondition
		sentry.CaptureException(err)
		c.JSON(422, utils.SendErrorResponse(err))
		return
	}
	ctx := c.Request.Context()
	resp := response.GetUserResponse{}

	response, err := helpers.FetchUser(ctx, &formRequest, span.Context()) //the DAO level
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Status = "Success"
	resp.Message = "User fecthed successfully"
	resp.Response = response
	span.Status = sentry.SpanStatusOK

	c.JSON(http.StatusOK, resp)

}
