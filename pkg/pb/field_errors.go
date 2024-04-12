package pb

import (
	"errors"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_errors"
)

// GetFieldValidationErrors parse the details from the error to obtain field errors
func GetFieldValidationErrors(
	returnedError error,
) ([]*pb_errors.FieldError, error) {
	errorDetails, ok := status.FromError(returnedError)
	if !ok {
		return nil, errors.New("error is not a status error")
	}
	fieldErrors := []*pb_errors.FieldError{}
	for _, detail := range errorDetails.Details() {
		if errDetail, ok := detail.(*anypb.Any); ok {
			fieldError := &pb_errors.FieldError{}
			if err := anypb.UnmarshalTo(errDetail, fieldError, proto.UnmarshalOptions{}); err == nil {
				fieldErrors = append(fieldErrors, fieldError)
			} else {
				return nil, err
			}
		}
	}
	return fieldErrors, nil
}
