// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: api/proto/session/v1/session.proto

package sessionv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SessionStatus int32

const (
	SessionStatus_SESSION_STATUS_UNSPECIFIED SessionStatus = 0
	SessionStatus_SESSION_STATUS_STARTED     SessionStatus = 1
	SessionStatus_SESSION_STATUS_IN_PROGRESS SessionStatus = 2
	SessionStatus_SESSION_STATUS_FINISHED    SessionStatus = 3
	SessionStatus_SESSION_STATUS_TIMEOUT     SessionStatus = 4
)

// Enum value maps for SessionStatus.
var (
	SessionStatus_name = map[int32]string{
		0: "SESSION_STATUS_UNSPECIFIED",
		1: "SESSION_STATUS_STARTED",
		2: "SESSION_STATUS_IN_PROGRESS",
		3: "SESSION_STATUS_FINISHED",
		4: "SESSION_STATUS_TIMEOUT",
	}
	SessionStatus_value = map[string]int32{
		"SESSION_STATUS_UNSPECIFIED": 0,
		"SESSION_STATUS_STARTED":     1,
		"SESSION_STATUS_IN_PROGRESS": 2,
		"SESSION_STATUS_FINISHED":    3,
		"SESSION_STATUS_TIMEOUT":     4,
	}
)

func (x SessionStatus) Enum() *SessionStatus {
	p := new(SessionStatus)
	*p = x
	return p
}

func (x SessionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SessionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_session_v1_session_proto_enumTypes[0].Descriptor()
}

func (SessionStatus) Type() protoreflect.EnumType {
	return &file_api_proto_session_v1_session_proto_enumTypes[0]
}

func (x SessionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SessionStatus.Descriptor instead.
func (SessionStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{0}
}

type ExamSession struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ExamId        string                 `protobuf:"bytes,2,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
	StudentId     string                 `protobuf:"bytes,3,opt,name=student_id,json=studentId,proto3" json:"student_id,omitempty"`
	Status        SessionStatus          `protobuf:"varint,4,opt,name=status,proto3,enum=session.v1.SessionStatus" json:"status,omitempty"`
	StartTime     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime       *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Answers       []*Answer              `protobuf:"bytes,7,rep,name=answers,proto3" json:"answers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExamSession) Reset() {
	*x = ExamSession{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExamSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExamSession) ProtoMessage() {}

func (x *ExamSession) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExamSession.ProtoReflect.Descriptor instead.
func (*ExamSession) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{0}
}

func (x *ExamSession) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ExamSession) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

func (x *ExamSession) GetStudentId() string {
	if x != nil {
		return x.StudentId
	}
	return ""
}

func (x *ExamSession) GetStatus() SessionStatus {
	if x != nil {
		return x.Status
	}
	return SessionStatus_SESSION_STATUS_UNSPECIFIED
}

func (x *ExamSession) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ExamSession) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *ExamSession) GetAnswers() []*Answer {
	if x != nil {
		return x.Answers
	}
	return nil
}

type Answer struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	QuestionId     string                 `protobuf:"bytes,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	SelectedChoice string                 `protobuf:"bytes,2,opt,name=selected_choice,json=selectedChoice,proto3" json:"selected_choice,omitempty"`
	AnsweredAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=answered_at,json=answeredAt,proto3" json:"answered_at,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Answer) Reset() {
	*x = Answer{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{1}
}

func (x *Answer) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *Answer) GetSelectedChoice() string {
	if x != nil {
		return x.SelectedChoice
	}
	return ""
}

func (x *Answer) GetAnsweredAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AnsweredAt
	}
	return nil
}

type StartSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ExamId        string                 `protobuf:"bytes,1,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
	StudentId     string                 `protobuf:"bytes,2,opt,name=student_id,json=studentId,proto3" json:"student_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartSessionRequest) Reset() {
	*x = StartSessionRequest{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartSessionRequest) ProtoMessage() {}

func (x *StartSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartSessionRequest.ProtoReflect.Descriptor instead.
func (*StartSessionRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{2}
}

func (x *StartSessionRequest) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

func (x *StartSessionRequest) GetStudentId() string {
	if x != nil {
		return x.StudentId
	}
	return ""
}

type GetSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSessionRequest) Reset() {
	*x = GetSessionRequest{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSessionRequest) ProtoMessage() {}

func (x *GetSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSessionRequest.ProtoReflect.Descriptor instead.
func (*GetSessionRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{3}
}

func (x *GetSessionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SubmitAnswerRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	SessionId      string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	QuestionId     string                 `protobuf:"bytes,2,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	SelectedChoice string                 `protobuf:"bytes,3,opt,name=selected_choice,json=selectedChoice,proto3" json:"selected_choice,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *SubmitAnswerRequest) Reset() {
	*x = SubmitAnswerRequest{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitAnswerRequest) ProtoMessage() {}

func (x *SubmitAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitAnswerRequest.ProtoReflect.Descriptor instead.
func (*SubmitAnswerRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{4}
}

func (x *SubmitAnswerRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *SubmitAnswerRequest) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *SubmitAnswerRequest) GetSelectedChoice() string {
	if x != nil {
		return x.SelectedChoice
	}
	return ""
}

type SubmitAnswerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitAnswerResponse) Reset() {
	*x = SubmitAnswerResponse{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitAnswerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitAnswerResponse) ProtoMessage() {}

func (x *SubmitAnswerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitAnswerResponse.ProtoReflect.Descriptor instead.
func (*SubmitAnswerResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{5}
}

func (x *SubmitAnswerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SubmitAnswerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type FinishSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FinishSessionRequest) Reset() {
	*x = FinishSessionRequest{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FinishSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishSessionRequest) ProtoMessage() {}

func (x *FinishSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishSessionRequest.ProtoReflect.Descriptor instead.
func (*FinishSessionRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{6}
}

func (x *FinishSessionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetRemainingTimeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionId     string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRemainingTimeRequest) Reset() {
	*x = GetRemainingTimeRequest{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRemainingTimeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRemainingTimeRequest) ProtoMessage() {}

func (x *GetRemainingTimeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRemainingTimeRequest.ProtoReflect.Descriptor instead.
func (*GetRemainingTimeRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{7}
}

func (x *GetRemainingTimeRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type GetRemainingTimeResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	RemainingMinutes int32                  `protobuf:"varint,1,opt,name=remaining_minutes,json=remainingMinutes,proto3" json:"remaining_minutes,omitempty"`
	RemainingSeconds int32                  `protobuf:"varint,2,opt,name=remaining_seconds,json=remainingSeconds,proto3" json:"remaining_seconds,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *GetRemainingTimeResponse) Reset() {
	*x = GetRemainingTimeResponse{}
	mi := &file_api_proto_session_v1_session_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRemainingTimeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRemainingTimeResponse) ProtoMessage() {}

func (x *GetRemainingTimeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_session_v1_session_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRemainingTimeResponse.ProtoReflect.Descriptor instead.
func (*GetRemainingTimeResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_session_v1_session_proto_rawDescGZIP(), []int{8}
}

func (x *GetRemainingTimeResponse) GetRemainingMinutes() int32 {
	if x != nil {
		return x.RemainingMinutes
	}
	return 0
}

func (x *GetRemainingTimeResponse) GetRemainingSeconds() int32 {
	if x != nil {
		return x.RemainingSeconds
	}
	return 0
}

var File_api_proto_session_v1_session_proto protoreflect.FileDescriptor

var file_api_proto_session_v1_session_proto_rawDesc = string([]byte{
	0x0a, 0x22, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa8, 0x02, 0x0a, 0x0b, 0x45, 0x78, 0x61, 0x6d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x65, 0x78, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0a,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2c,
	0x0a, 0x07, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x52, 0x07, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x22, 0x8f, 0x01, 0x0a,
	0x06, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x43, 0x68, 0x6f, 0x69, 0x63,
	0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0a, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x65, 0x64, 0x41, 0x74, 0x22, 0x4d,
	0x0a, 0x13, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x78, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x23, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x7e, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x43, 0x68, 0x6f, 0x69,
	0x63, 0x65, 0x22, 0x4a, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x26,
	0x0a, 0x14, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x38, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6d,
	0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x22, 0x74, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67,
	0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x11,
	0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69,
	0x6e, 0x67, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x72, 0x65, 0x6d,
	0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x53,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x2a, 0xa4, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x54,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45,
	0x53, 0x53, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x04, 0x32, 0xa8, 0x03,
	0x0a, 0x0e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4a, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1f, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x78, 0x61, 0x6d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0a,
	0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x0d, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x2e, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x5f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x23, 0x2e, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6d, 0x61,
	0x69, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x70, 0x65, 0x73, 0x4a, 0x73, 0x2f, 0x63, 0x62,
	0x74, 0x2d, 0x65, 0x78, 0x61, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_proto_session_v1_session_proto_rawDescOnce sync.Once
	file_api_proto_session_v1_session_proto_rawDescData []byte
)

func file_api_proto_session_v1_session_proto_rawDescGZIP() []byte {
	file_api_proto_session_v1_session_proto_rawDescOnce.Do(func() {
		file_api_proto_session_v1_session_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_session_v1_session_proto_rawDesc), len(file_api_proto_session_v1_session_proto_rawDesc)))
	})
	return file_api_proto_session_v1_session_proto_rawDescData
}

var file_api_proto_session_v1_session_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_session_v1_session_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_proto_session_v1_session_proto_goTypes = []any{
	(SessionStatus)(0),               // 0: session.v1.SessionStatus
	(*ExamSession)(nil),              // 1: session.v1.ExamSession
	(*Answer)(nil),                   // 2: session.v1.Answer
	(*StartSessionRequest)(nil),      // 3: session.v1.StartSessionRequest
	(*GetSessionRequest)(nil),        // 4: session.v1.GetSessionRequest
	(*SubmitAnswerRequest)(nil),      // 5: session.v1.SubmitAnswerRequest
	(*SubmitAnswerResponse)(nil),     // 6: session.v1.SubmitAnswerResponse
	(*FinishSessionRequest)(nil),     // 7: session.v1.FinishSessionRequest
	(*GetRemainingTimeRequest)(nil),  // 8: session.v1.GetRemainingTimeRequest
	(*GetRemainingTimeResponse)(nil), // 9: session.v1.GetRemainingTimeResponse
	(*timestamppb.Timestamp)(nil),    // 10: google.protobuf.Timestamp
}
var file_api_proto_session_v1_session_proto_depIdxs = []int32{
	0,  // 0: session.v1.ExamSession.status:type_name -> session.v1.SessionStatus
	10, // 1: session.v1.ExamSession.start_time:type_name -> google.protobuf.Timestamp
	10, // 2: session.v1.ExamSession.end_time:type_name -> google.protobuf.Timestamp
	2,  // 3: session.v1.ExamSession.answers:type_name -> session.v1.Answer
	10, // 4: session.v1.Answer.answered_at:type_name -> google.protobuf.Timestamp
	3,  // 5: session.v1.SessionService.StartSession:input_type -> session.v1.StartSessionRequest
	4,  // 6: session.v1.SessionService.GetSession:input_type -> session.v1.GetSessionRequest
	5,  // 7: session.v1.SessionService.SubmitAnswer:input_type -> session.v1.SubmitAnswerRequest
	7,  // 8: session.v1.SessionService.FinishSession:input_type -> session.v1.FinishSessionRequest
	8,  // 9: session.v1.SessionService.GetRemainingTime:input_type -> session.v1.GetRemainingTimeRequest
	1,  // 10: session.v1.SessionService.StartSession:output_type -> session.v1.ExamSession
	1,  // 11: session.v1.SessionService.GetSession:output_type -> session.v1.ExamSession
	6,  // 12: session.v1.SessionService.SubmitAnswer:output_type -> session.v1.SubmitAnswerResponse
	1,  // 13: session.v1.SessionService.FinishSession:output_type -> session.v1.ExamSession
	9,  // 14: session.v1.SessionService.GetRemainingTime:output_type -> session.v1.GetRemainingTimeResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_proto_session_v1_session_proto_init() }
func file_api_proto_session_v1_session_proto_init() {
	if File_api_proto_session_v1_session_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_session_v1_session_proto_rawDesc), len(file_api_proto_session_v1_session_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_session_v1_session_proto_goTypes,
		DependencyIndexes: file_api_proto_session_v1_session_proto_depIdxs,
		EnumInfos:         file_api_proto_session_v1_session_proto_enumTypes,
		MessageInfos:      file_api_proto_session_v1_session_proto_msgTypes,
	}.Build()
	File_api_proto_session_v1_session_proto = out.File
	file_api_proto_session_v1_session_proto_goTypes = nil
	file_api_proto_session_v1_session_proto_depIdxs = nil
}
