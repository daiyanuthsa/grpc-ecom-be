package utils

import (
	"github.com/daiyanuthsa/grpc-ecom-be/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:   true,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		ValidationErrors: validationErrors,
		Message:          "Validation Error",
		IsError:         true,
	}
}
func UnauthenticatedResponse()error{
	return  status.Error(codes.Unauthenticated,"Unauthenticated")
}