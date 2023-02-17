package PATCH

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

func UpdateUser(c *gin.Context) {
	defer sentry.Recover()
	span := sentry.StartSpan(context.TODO(), "[GIN] UserDetails", sentry.TransactionName("Create a new user"))
	defer span.Finish()

	formRequest := structs.UserDetails{}

	if err := c.ShouldBind(&formRequest); err != nil {
		span.Status = sentry.SpanStatusFailedPrecondition
		sentry.CaptureException(err)
		c.JSON(422, utils.SendErrorResponse(err))
		return
	}
	ctx := c.Request.Context()
	resp := response.SuccessResponse{}

	err := helpers.UpdateUserDetails(ctx, &formRequest, span.Context()) //the DAO level
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Status = "Success"
	resp.Message = "User details updated successfully"
	span.Status = sentry.SpanStatusOK

	c.JSON(http.StatusOK, resp)

}
