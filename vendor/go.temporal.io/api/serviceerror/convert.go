// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package serviceerror

import (
	"context"
	"errors"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"

	"go.temporal.io/api/errordetails/v1"
)

// ToStatus converts service error to gogo gRPC Status.
// If error is not a service error it returns status with code Unknown.
func ToStatus(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "")
	}

	if svcerr, ok := err.(ServiceError); ok {
		return svcerr.Status()
	}

	// Special case for context.DeadlineExceeded and context.Canceled because they can happen in unpredictable places.
	if errors.Is(err, context.DeadlineExceeded) {
		return status.New(codes.DeadlineExceeded, err.Error())
	}
	if errors.Is(err, context.Canceled) {
		return status.New(codes.Canceled, err.Error())
	}

	// Internal logic of status.Convert is:
	//   - if err is already gogo Status or gRPC Status, then just return it (this should never happen though).
	//   - otherwise returns codes.Unknown with message from err.Error() (this might happen if some generic go error reach to this point).
	return status.Convert(err)
}

// FromStatus converts gogo gRPC Status to service error.
func FromStatus(st *status.Status) error {
	if st == nil || st.Code() == codes.OK {
		return nil
	}

	// Simple case. Code to serviceerror is one to one mapping and there are no error details.
	switch st.Code() {
	case codes.Internal:
		return newInternal(st)
	case codes.DataLoss:
		return newDataLoss(st)
	case codes.ResourceExhausted:
		return newResourceExhausted(st)
	case codes.DeadlineExceeded:
		return newDeadlineExceeded(st)
	case codes.Canceled:
		return newCanceled(st)
	case codes.Unavailable:
		return newUnavailable(st)
	case codes.Unimplemented:
		return newUnimplemented(st)
	case codes.Unknown:
		// Unwrap error message from unknown error.
		return errors.New(st.Message())
	// Unsupported codes.
	case codes.OutOfRange,
		codes.Unauthenticated:
		// Use standard gRPC error representation for unsupported codes ("rpc error: code = %s desc = %s").
		return st.Err()
	}

	errDetails := extractErrorDetails(st)

	// If there was an error during details extraction, it will go to errDetails.
	if err, ok := errDetails.(error); ok {
		return NewInvalidArgument(err.Error())
	}

	switch st.Code() {
	case codes.NotFound:
		if errDetails == nil {
			return newNotFound(st, nil)
		}
		switch errDetails := errDetails.(type) {
		case *errordetails.NotFoundFailure:
			return newNotFound(st, errDetails)
		}
	case codes.InvalidArgument:
		if errDetails == nil {
			return newInvalidArgument(st)
		}
		switch errDetails.(type) {
		case *errordetails.QueryFailedFailure:
			return newQueryFailed(st)
		}
	case codes.AlreadyExists:
		switch errDetails := errDetails.(type) {
		case *errordetails.NamespaceAlreadyExistsFailure:
			return newNamespaceAlreadyExists(st)
		case *errordetails.WorkflowExecutionAlreadyStartedFailure:
			return newWorkflowExecutionAlreadyStarted(st, errDetails)
		case *errordetails.CancellationAlreadyRequestedFailure:
			return newCancellationAlreadyRequested(st)
		}
	case codes.FailedPrecondition:
		switch errDetails := errDetails.(type) {
		case *errordetails.NamespaceNotActiveFailure:
			return newNamespaceNotActive(st, errDetails)
		case *errordetails.ClientVersionNotSupportedFailure:
			return newClientVersionNotSupported(st, errDetails)
		case *errordetails.ServerVersionNotSupportedFailure:
			return newServerVersionNotSupported(st, errDetails)
		}
	case codes.PermissionDenied:
		switch errDetails := errDetails.(type) {
		case *errordetails.PermissionDeniedFailure:
			return newPermissionDenied(st, errDetails)
		}
		return newPermissionDenied(st, nil)
	}

	// st.Code() should have error details but it didn't (or error details are of a wrong type).
	// Then use standard gRPC error representation ("rpc error: code = %s desc = %s").
	return st.Err()
}

func extractErrorDetails(st *status.Status) interface{} {
	details := st.Details()
	if len(details) > 0 {
		return details[0]
	}

	return nil
}
