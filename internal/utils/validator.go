package utils

import (
	"errors"
	"mime/multipart"
	"strings"

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

type ImageValidator struct {
    MaxSizeBytes int64
    AllowedMimeTypes map[string]string
}

func (v *ImageValidator) Validate(file *multipart.FileHeader) error {
    // 1. Cek Ukuran
    if file.Size > v.MaxSizeBytes {
        return errors.New("image size exceeds limit")
    }

    // 2. Cek MIME/Content-Type
    contentType := strings.ToLower(file.Header.Get("Content-Type"))
    _, mimeOk := v.AllowedMimeTypes[contentType]
    if !mimeOk {
        return errors.New("content type is not allowed")
    }

    
    return nil
}