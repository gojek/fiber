// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: testdata/proto/upi.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PredictValuesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PredictionRows []*PredictionRow `protobuf:"bytes,1,rep,name=prediction_rows,json=predictionRows,proto3" json:"prediction_rows,omitempty"`
	Metadata       *RequestMetadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *PredictValuesRequest) Reset() {
	*x = PredictValuesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictValuesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictValuesRequest) ProtoMessage() {}

func (x *PredictValuesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictValuesRequest.ProtoReflect.Descriptor instead.
func (*PredictValuesRequest) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{0}
}

func (x *PredictValuesRequest) GetPredictionRows() []*PredictionRow {
	if x != nil {
		return x.PredictionRows
	}
	return nil
}

func (x *PredictValuesRequest) GetMetadata() *RequestMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type RequestMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestTimestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=request_timestamp,json=requestTimestamp,proto3" json:"request_timestamp,omitempty"`
	TargetName       string                 `protobuf:"bytes,2,opt,name=target_name,json=targetName,proto3" json:"target_name,omitempty"`
	Others           []*NamedValue          `protobuf:"bytes,3,rep,name=others,proto3" json:"others,omitempty"`
}

func (x *RequestMetadata) Reset() {
	*x = RequestMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMetadata) ProtoMessage() {}

func (x *RequestMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMetadata.ProtoReflect.Descriptor instead.
func (*RequestMetadata) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{1}
}

func (x *RequestMetadata) GetRequestTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.RequestTimestamp
	}
	return nil
}

func (x *RequestMetadata) GetTargetName() string {
	if x != nil {
		return x.TargetName
	}
	return ""
}

func (x *RequestMetadata) GetOthers() []*NamedValue {
	if x != nil {
		return x.Others
	}
	return nil
}

type PredictionRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowId              string        `protobuf:"bytes,1,opt,name=row_id,json=rowId,proto3" json:"row_id,omitempty"`
	ContextualFeatures []*NamedValue `protobuf:"bytes,2,rep,name=contextual_features,json=contextualFeatures,proto3" json:"contextual_features,omitempty"`
	RowEntities        []*NamedValue `protobuf:"bytes,3,rep,name=row_entities,json=rowEntities,proto3" json:"row_entities,omitempty"`
}

func (x *PredictionRow) Reset() {
	*x = PredictionRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictionRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictionRow) ProtoMessage() {}

func (x *PredictionRow) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictionRow.ProtoReflect.Descriptor instead.
func (*PredictionRow) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{2}
}

func (x *PredictionRow) GetRowId() string {
	if x != nil {
		return x.RowId
	}
	return ""
}

func (x *PredictionRow) GetContextualFeatures() []*NamedValue {
	if x != nil {
		return x.ContextualFeatures
	}
	return nil
}

func (x *PredictionRow) GetRowEntities() []*NamedValue {
	if x != nil {
		return x.RowEntities
	}
	return nil
}

type PredictValuesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Predictions []*PredictionResult `protobuf:"bytes,1,rep,name=predictions,proto3" json:"predictions,omitempty"`
	Metadata    *ResponseMetadata   `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *PredictValuesResponse) Reset() {
	*x = PredictValuesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictValuesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictValuesResponse) ProtoMessage() {}

func (x *PredictValuesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictValuesResponse.ProtoReflect.Descriptor instead.
func (*PredictValuesResponse) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{3}
}

func (x *PredictValuesResponse) GetPredictions() []*PredictionResult {
	if x != nil {
		return x.Predictions
	}
	return nil
}

func (x *PredictValuesResponse) GetMetadata() *ResponseMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type PredictionResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowId string      `protobuf:"bytes,1,opt,name=row_id,json=rowId,proto3" json:"row_id,omitempty"`
	Value *NamedValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *PredictionResult) Reset() {
	*x = PredictionResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictionResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictionResult) ProtoMessage() {}

func (x *PredictionResult) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictionResult.ProtoReflect.Descriptor instead.
func (*PredictionResult) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{4}
}

func (x *PredictionResult) GetRowId() string {
	if x != nil {
		return x.RowId
	}
	return ""
}

func (x *PredictionResult) GetValue() *NamedValue {
	if x != nil {
		return x.Value
	}
	return nil
}

type ResponseMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PredictionId string        `protobuf:"bytes,1,opt,name=prediction_id,json=predictionId,proto3" json:"prediction_id,omitempty"`
	TargetName   string        `protobuf:"bytes,2,opt,name=target_name,json=targetName,proto3" json:"target_name,omitempty"`
	ModelName    string        `protobuf:"bytes,3,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"`
	ModelVersion string        `protobuf:"bytes,4,opt,name=model_version,json=modelVersion,proto3" json:"model_version,omitempty"`
	ExperimentId string        `protobuf:"bytes,5,opt,name=experiment_id,json=experimentId,proto3" json:"experiment_id,omitempty"`
	TreatmentId  string        `protobuf:"bytes,6,opt,name=treatment_id,json=treatmentId,proto3" json:"treatment_id,omitempty"`
	Others       []*NamedValue `protobuf:"bytes,7,rep,name=others,proto3" json:"others,omitempty"`
}

func (x *ResponseMetadata) Reset() {
	*x = ResponseMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_proto_upi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseMetadata) ProtoMessage() {}

func (x *ResponseMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_proto_upi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseMetadata.ProtoReflect.Descriptor instead.
func (*ResponseMetadata) Descriptor() ([]byte, []int) {
	return file_testdata_proto_upi_proto_rawDescGZIP(), []int{5}
}

func (x *ResponseMetadata) GetPredictionId() string {
	if x != nil {
		return x.PredictionId
	}
	return ""
}

func (x *ResponseMetadata) GetTargetName() string {
	if x != nil {
		return x.TargetName
	}
	return ""
}

func (x *ResponseMetadata) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *ResponseMetadata) GetModelVersion() string {
	if x != nil {
		return x.ModelVersion
	}
	return ""
}

func (x *ResponseMetadata) GetExperimentId() string {
	if x != nil {
		return x.ExperimentId
	}
	return ""
}

func (x *ResponseMetadata) GetTreatmentId() string {
	if x != nil {
		return x.TreatmentId
	}
	return ""
}

func (x *ResponseMetadata) GetOthers() []*NamedValue {
	if x != nil {
		return x.Others
	}
	return nil
}

var File_testdata_proto_upi_proto protoreflect.FileDescriptor

var file_testdata_proto_upi_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x75, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x65, 0x73, 0x74,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x14, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x0f, 0x70,
	0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x77, 0x52, 0x0e,
	0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x77, 0x73, 0x12, 0x36,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0xaa, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x47, 0x0a, 0x11, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x10, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x6f, 0x74, 0x68,
	0x65, 0x72, 0x73, 0x22, 0xa8, 0x01, 0x0a, 0x0d, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x6f, 0x77, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x77, 0x49, 0x64, 0x12, 0x46, 0x0a, 0x13,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x12, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x75, 0x61, 0x6c, 0x46, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x0c, 0x72, 0x6f, 0x77, 0x5f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x0b, 0x72, 0x6f, 0x77, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x22, 0x8f,
	0x01, 0x0a, 0x15, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x37, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x56, 0x0a, 0x10, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x77, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x93, 0x02, 0x0a, 0x10, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a,
	0x0d, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x78, 0x70, 0x65, 0x72,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x74, 0x72, 0x65, 0x61, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x74, 0x72, 0x65, 0x61, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x2d, 0x0a, 0x06, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65,
	0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x32, 0x70,
	0x0a, 0x1a, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x61, 0x6c, 0x50, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x0d,
	0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x1f, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63,
	0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x84, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x42, 0x08, 0x55, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x25,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6a, 0x65, 0x6b,
	0x2f, 0x66, 0x69, 0x62, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0xa2, 0x02, 0x03, 0x54, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x54, 0x65,
	0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0xca, 0x02, 0x09, 0x54, 0x65, 0x73, 0x74, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0xe2, 0x02, 0x15, 0x54, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09, 0x54, 0x65,
	0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testdata_proto_upi_proto_rawDescOnce sync.Once
	file_testdata_proto_upi_proto_rawDescData = file_testdata_proto_upi_proto_rawDesc
)

func file_testdata_proto_upi_proto_rawDescGZIP() []byte {
	file_testdata_proto_upi_proto_rawDescOnce.Do(func() {
		file_testdata_proto_upi_proto_rawDescData = protoimpl.X.CompressGZIP(file_testdata_proto_upi_proto_rawDescData)
	})
	return file_testdata_proto_upi_proto_rawDescData
}

var file_testdata_proto_upi_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_testdata_proto_upi_proto_goTypes = []interface{}{
	(*PredictValuesRequest)(nil),  // 0: testproto.PredictValuesRequest
	(*RequestMetadata)(nil),       // 1: testproto.RequestMetadata
	(*PredictionRow)(nil),         // 2: testproto.PredictionRow
	(*PredictValuesResponse)(nil), // 3: testproto.PredictValuesResponse
	(*PredictionResult)(nil),      // 4: testproto.PredictionResult
	(*ResponseMetadata)(nil),      // 5: testproto.ResponseMetadata
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*NamedValue)(nil),            // 7: testproto.NamedValue
}
var file_testdata_proto_upi_proto_depIdxs = []int32{
	2,  // 0: testproto.PredictValuesRequest.prediction_rows:type_name -> testproto.PredictionRow
	1,  // 1: testproto.PredictValuesRequest.metadata:type_name -> testproto.RequestMetadata
	6,  // 2: testproto.RequestMetadata.request_timestamp:type_name -> google.protobuf.Timestamp
	7,  // 3: testproto.RequestMetadata.others:type_name -> testproto.NamedValue
	7,  // 4: testproto.PredictionRow.contextual_features:type_name -> testproto.NamedValue
	7,  // 5: testproto.PredictionRow.row_entities:type_name -> testproto.NamedValue
	4,  // 6: testproto.PredictValuesResponse.predictions:type_name -> testproto.PredictionResult
	5,  // 7: testproto.PredictValuesResponse.metadata:type_name -> testproto.ResponseMetadata
	7,  // 8: testproto.PredictionResult.value:type_name -> testproto.NamedValue
	7,  // 9: testproto.ResponseMetadata.others:type_name -> testproto.NamedValue
	0,  // 10: testproto.UniversalPredictionService.PredictValues:input_type -> testproto.PredictValuesRequest
	3,  // 11: testproto.UniversalPredictionService.PredictValues:output_type -> testproto.PredictValuesResponse
	11, // [11:12] is the sub-list for method output_type
	10, // [10:11] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_testdata_proto_upi_proto_init() }
func file_testdata_proto_upi_proto_init() {
	if File_testdata_proto_upi_proto != nil {
		return
	}
	file_testdata_proto_value_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_testdata_proto_upi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictValuesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testdata_proto_upi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testdata_proto_upi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictionRow); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testdata_proto_upi_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictValuesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testdata_proto_upi_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictionResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_testdata_proto_upi_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_testdata_proto_upi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_testdata_proto_upi_proto_goTypes,
		DependencyIndexes: file_testdata_proto_upi_proto_depIdxs,
		MessageInfos:      file_testdata_proto_upi_proto_msgTypes,
	}.Build()
	File_testdata_proto_upi_proto = out.File
	file_testdata_proto_upi_proto_rawDesc = nil
	file_testdata_proto_upi_proto_goTypes = nil
	file_testdata_proto_upi_proto_depIdxs = nil
}
