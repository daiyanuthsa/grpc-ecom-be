package utils

import (
	"errors"

	"buf.build/go/protovalidate"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/common"
	"google.golang.org/protobuf/proto"
)

func CheckValidation(req proto.Message) ([]*common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validationError *protovalidate.ValidationError
		if errors.As(err, &validationError){
			var validationErrorResponse []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, violation := range validationError.Violations {
				validationErrorResponse = append(validationErrorResponse, &common.ValidationError{
					Field:   *violation.Proto.Field.Elements[0].FieldName,
					Message: *violation.Proto.Message,
				})
			}
			// If the error is not a ValidationError, return a generic error
		return validationErrorResponse, nil
		}
		return nil, err
	}
	return make([]*common.ValidationError, 0), nil
}
