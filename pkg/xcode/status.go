package xcode

import (
	"context"
	"fmt"
	types "github.com/lustresix/beifeng/pkg/xcode/type"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// XCode 定义了自定义错误码的接口。
var _ XCode = (*Status)(nil)

// Status 结构体包装了 protobuf 中的 Status 消息，实现了 XCode 接口。
type Status struct {
	sts *types.Status
}

// Error 创建一个包含指定错误码的 Status 实例。
func Error(code Code) *Status {
	return &Status{sts: &types.Status{Code: int32(code.Code()), Message: code.Message()}}
}

// Errorf 创建一个包含指定格式化错误码的 Status 实例。
func Errorf(code Code, format string, args ...interface{}) *Status {
	code.msg = fmt.Sprintf(format, args...)
	return Error(code)
}

// Error 返回错误消息。
func (s *Status) Error() string {
	return s.Message()
}

// Code 返回错误码。
func (s *Status) Code() int {
	return int(s.sts.Code)
}

// Message 返回错误消息，如果消息为空，则返回错误码的字符串形式。
func (s *Status) Message() string {
	if s.sts.Message == "" {
		return strconv.Itoa(int(s.sts.Code))
	}
	return s.sts.Message
}

// Details 返回错误的详细信息。
func (s *Status) Details() []interface{} {
	if s == nil || s.sts == nil {
		return nil
	}
	details := make([]interface{}, 0, len(s.sts.Details))
	for _, d := range s.sts.Details {
		detail := &ptypes.DynamicAny{}
		if err := d.UnmarshalTo(detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

// WithDetails 将 proto.Message 类型的详细信息附加到错误中。
func (s *Status) WithDetails(msgs ...proto.Message) (*Status, error) {
	for _, msg := range msgs {
		anyMsg, err := anypb.New(msg)
		if err != nil {
			return s, err
		}
		s.sts.Details = append(s.sts.Details, anyMsg)
	}
	return s, nil
}

// Proto 返回包装的 protobuf Status 消息。
func (s *Status) Proto() *types.Status {
	return s.sts
}

// FromCode 根据错误码创建 Status 实例。
func FromCode(code Code) *Status {
	return &Status{sts: &types.Status{Code: int32(code.Code()), Message: code.Message()}}
}

// FromProto 根据 proto.Message 创建 XCode 实例。
func FromProto(pbMsg proto.Message) XCode {
	msg, ok := pbMsg.(*types.Status)
	if ok {
		if len(msg.Message) == 0 || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			return Code{code: int(msg.Code)}
		}
		return &Status{sts: msg}
	}
	return Errorf(ServerErr, "invalid proto message get %v", pbMsg)
}

// toXCode 将 gRPC Status 转换为自定义错误码。
func toXCode(grpcStatus *status.Status) Code {
	grpcCode := grpcStatus.Code()
	switch grpcCode {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return RequestErr
	case codes.NotFound:
		return NotFound
	case codes.PermissionDenied:
		return AccessDenied
	case codes.Unauthenticated:
		return Unauthorized
	case codes.ResourceExhausted:
		return LimitExceed
	case codes.Unimplemented:
		return MethodNotAllowed
	case codes.DeadlineExceeded:
		return Deadline
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.Unknown:
		return String(grpcStatus.Message())
	}
	return ServerErr
}

// CodeFromError 将错误转换为自定义错误码。
func CodeFromError(err error) XCode {
	err = errors.Cause(err)
	var code XCode
	if errors.As(err, &code) {
		return code
	}
	switch {
	case errors.Is(err, context.Canceled):
		return Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return Deadline
	}
	return ServerErr
}

// FromError 将错误转换为 gRPC Status。
func FromError(err error) *status.Status {
	err = errors.Cause(err)
	var code XCode
	if errors.As(err, &code) {
		grpcStatus, e := gRPCStatusFromXCode(code)
		if e == nil {
			return grpcStatus
		}
	}
	var grpcStatus *status.Status
	switch {
	case errors.Is(err, context.Canceled):
		grpcStatus, _ = gRPCStatusFromXCode(Canceled)
	case errors.Is(err, context.DeadlineExceeded):
		grpcStatus, _ = gRPCStatusFromXCode(Deadline)
	default:
		grpcStatus, _ = status.FromError(err)
	}
	return grpcStatus
}

// gRPCStatusFromXCode 将自定义错误码转换为 gRPC Status。
func gRPCStatusFromXCode(code XCode) (*status.Status, error) {
	var sts *Status
	switch v := code.(type) {
	case *Status:
		sts = v
	case Code:
		sts = FromCode(v)
	default:
		sts = Error(Code{code.Code(), code.Message()})
		for _, detail := range code.Details() {
			if msg, ok := detail.(proto.Message); ok {
				_, _ = sts.WithDetails(msg)
			}
		}
	}
	st := status.New(codes.Unknown, strconv.Itoa(sts.Code()))
	return st.WithDetails(sts.Proto())
}

func GrpcStatusToXCode(rpcStatus *status.Status) XCode {
	details := rpcStatus.Details()
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		if pb, ok := detail.(proto.Message); ok {
			return FromProto(pb)
		}
	}

	return toXCode(rpcStatus)
}
