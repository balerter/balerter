// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: google/rpc/status.proto

package rpc

import (
	bytes "bytes"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// The `Status` type defines a logical error model that is suitable for
// different programming environments, including REST APIs and RPC APIs. It is
// used by [gRPC](https://github.com/grpc). Each `Status` message contains
// three pieces of data: error code, error message, and error details.
//
// You can find out more about this error model and how to work with it in the
// [API Design Guide](https://cloud.google.com/apis/design/errors).
type Status struct {
	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// A developer-facing error message, which should be in English. Any
	// user-facing error message should be localized and sent in the
	// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// A list of messages that carry the error details.  There is a common set of
	// message types for APIs to use.
	Details              []*types.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Status) Reset()      { *m = Status{} }
func (*Status) ProtoMessage() {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d244abaf643bfe, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Status) GetDetails() []*types.Any {
	if m != nil {
		return m.Details
	}
	return nil
}

func (*Status) XXX_MessageName() string {
	return "google.rpc.Status"
}
func init() {
	proto.RegisterType((*Status)(nil), "google.rpc.Status")
}

func init() { proto.RegisterFile("google/rpc/status.proto", fileDescriptor_24d244abaf643bfe) }

var fileDescriptor_24d244abaf643bfe = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x2f, 0x2a, 0x48, 0xd6, 0x2f, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0xd6, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x82, 0x48, 0xe8, 0x15, 0x15, 0x24, 0x4b, 0x49, 0x42, 0x15, 0x81,
	0x65, 0x92, 0x4a, 0xd3, 0xf4, 0x13, 0xf3, 0x2a, 0x21, 0xca, 0x94, 0xd2, 0xb8, 0xd8, 0x82, 0xc1,
	0xda, 0x84, 0x84, 0xb8, 0x58, 0x92, 0xf3, 0x53, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83,
	0xc0, 0x6c, 0x21, 0x09, 0x2e, 0xf6, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x09, 0x26, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0x48, 0x8f, 0x8b, 0x3d, 0x25, 0xb5, 0x24, 0x31, 0x33, 0xa7,
	0x58, 0x82, 0x59, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x44, 0x0f, 0x6a, 0x21, 0xcc, 0x12, 0x3d, 0xc7,
	0xbc, 0xca, 0x20, 0x98, 0x22, 0xa7, 0xc4, 0x0b, 0x0f, 0xe5, 0x18, 0x6e, 0x3c, 0x94, 0x63, 0xf8,
	0xf0, 0x50, 0x8e, 0xf1, 0xc7, 0x43, 0x39, 0xc6, 0x86, 0x47, 0x72, 0x8c, 0x2b, 0x1e, 0xc9, 0x31,
	0x9e, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x2f, 0x1e, 0xc9,
	0x31, 0x7c, 0x00, 0x89, 0x3f, 0x96, 0x63, 0x3c, 0xf1, 0x58, 0x8e, 0x91, 0x8b, 0x2f, 0x39, 0x3f,
	0x57, 0x0f, 0xe1, 0x11, 0x27, 0x6e, 0x88, 0x5b, 0x03, 0x40, 0x56, 0x04, 0x30, 0x46, 0x31, 0x17,
	0x15, 0x24, 0xff, 0x60, 0x64, 0x5c, 0xc4, 0xc4, 0x1c, 0x14, 0xe0, 0x9c, 0xc4, 0x06, 0xb6, 0xd9,
	0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x16, 0xcd, 0x7b, 0x60, 0x13, 0x01, 0x00, 0x00,
}

func (this *Status) Compare(that interface{}) int {
	if that == nil {
		if this == nil {
			return 0
		}
		return 1
	}

	that1, ok := that.(*Status)
	if !ok {
		that2, ok := that.(Status)
		if ok {
			that1 = &that2
		} else {
			return 1
		}
	}
	if that1 == nil {
		if this == nil {
			return 0
		}
		return 1
	} else if this == nil {
		return -1
	}
	if this.Code != that1.Code {
		if this.Code < that1.Code {
			return -1
		}
		return 1
	}
	if this.Message != that1.Message {
		if this.Message < that1.Message {
			return -1
		}
		return 1
	}
	if len(this.Details) != len(that1.Details) {
		if len(this.Details) < len(that1.Details) {
			return -1
		}
		return 1
	}
	for i := range this.Details {
		if c := this.Details[i].Compare(that1.Details[i]); c != 0 {
			return c
		}
	}
	if c := bytes.Compare(this.XXX_unrecognized, that1.XXX_unrecognized); c != 0 {
		return c
	}
	return 0
}
func (this *Status) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Status)
	if !ok {
		that2, ok := that.(Status)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	if len(this.Details) != len(that1.Details) {
		return false
	}
	for i := range this.Details {
		if !this.Details[i].Equal(that1.Details[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Status) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&rpc.Status{")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	if this.Details != nil {
		s = append(s, "Details: "+fmt.Sprintf("%#v", this.Details)+",\n")
	}
	if this.XXX_unrecognized != nil {
		s = append(s, "XXX_unrecognized:"+fmt.Sprintf("%#v", this.XXX_unrecognized)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringStatus(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Details) > 0 {
		for iNdEx := len(m.Details) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Details[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintStatus(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintStatus(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintStatus(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintStatus(dAtA []byte, offset int, v uint64) int {
	offset -= sovStatus(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func NewPopulatedStatus(r randyStatus, easy bool) *Status {
	this := &Status{}
	this.Code = int32(r.Int31())
	if r.Intn(2) == 0 {
		this.Code *= -1
	}
	this.Message = string(randStringStatus(r))
	if r.Intn(5) != 0 {
		v1 := r.Intn(5)
		this.Details = make([]*types.Any, v1)
		for i := 0; i < v1; i++ {
			this.Details[i] = types.NewPopulatedAny(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedStatus(r, 4)
	}
	return this
}

type randyStatus interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneStatus(r randyStatus) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringStatus(r randyStatus) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneStatus(r)
	}
	return string(tmps)
}
func randUnrecognizedStatus(r randyStatus, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldStatus(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldStatus(dAtA []byte, r randyStatus, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateStatus(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateStatus(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovStatus(uint64(m.Code))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovStatus(uint64(l))
	}
	if len(m.Details) > 0 {
		for _, e := range m.Details {
			l = e.Size()
			n += 1 + l + sovStatus(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStatus(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStatus(x uint64) (n int) {
	return sovStatus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Status) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForDetails := "[]*Any{"
	for _, f := range this.Details {
		repeatedStringForDetails += strings.Replace(fmt.Sprintf("%v", f), "Any", "types.Any", 1) + ","
	}
	repeatedStringForDetails += "}"
	s := strings.Join([]string{`&Status{`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`Details:` + repeatedStringForDetails + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringStatus(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStatus
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStatus
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStatus
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Details", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStatus
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStatus
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Details = append(m.Details, &types.Any{})
			if err := m.Details[len(m.Details)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStatus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStatus
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStatus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStatus
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStatus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStatus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthStatus
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStatus
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStatus
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStatus        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStatus          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStatus = fmt.Errorf("proto: unexpected end of group")
)
