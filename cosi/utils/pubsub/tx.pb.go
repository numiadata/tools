// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosi/v1/tx.proto

package pubsub

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	types "github.com/tendermint/tendermint/abci/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TxResult struct {
	Height int64                   `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Index  uint32                  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Tx     []byte                  `protobuf:"bytes,3,opt,name=tx,proto3" json:"tx,omitempty"`
	Result types.ResponseDeliverTx `protobuf:"bytes,4,opt,name=result,proto3" json:"result"`
	Time   time.Time               `protobuf:"bytes,5,opt,name=time,proto3,stdtime" json:"time"`
}

func (m *TxResult) Reset()         { *m = TxResult{} }
func (m *TxResult) String() string { return proto.CompactTextString(m) }
func (*TxResult) ProtoMessage()    {}
func (*TxResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f4b0e93112b5f1a, []int{0}
}
func (m *TxResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxResult.Merge(m, src)
}
func (m *TxResult) XXX_Size() int {
	return m.Size()
}
func (m *TxResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TxResult.DiscardUnknown(m)
}

var xxx_messageInfo_TxResult proto.InternalMessageInfo

func (m *TxResult) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *TxResult) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *TxResult) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *TxResult) GetResult() types.ResponseDeliverTx {
	if m != nil {
		return m.Result
	}
	return types.ResponseDeliverTx{}
}

func (m *TxResult) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*TxResult)(nil), "pubsub.tx.TxResult")
}

func init() { proto.RegisterFile("cosi/v1/tx.proto", fileDescriptor_5f4b0e93112b5f1a) }

var fileDescriptor_5f4b0e93112b5f1a = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xbd, 0x6a, 0xc3, 0x30,
	0x10, 0xc7, 0xad, 0x7c, 0x91, 0xaa, 0x1f, 0x14, 0x13, 0x8a, 0x49, 0xc1, 0x31, 0x99, 0x3c, 0x14,
	0x89, 0xb6, 0x4b, 0xc7, 0x12, 0x4a, 0x1f, 0x40, 0x64, 0xea, 0x66, 0x27, 0x57, 0x47, 0x60, 0x5b,
	0xc6, 0x3a, 0x05, 0xf5, 0x2d, 0xf2, 0x50, 0x1d, 0x32, 0x66, 0xec, 0xd4, 0x96, 0xe4, 0x45, 0x8a,
	0xe5, 0x84, 0x6e, 0x77, 0x3a, 0xfd, 0xee, 0xcf, 0xef, 0xe8, 0xf5, 0x42, 0x69, 0xc9, 0xd7, 0xf7,
	0x1c, 0x2d, 0xab, 0x6a, 0x85, 0xca, 0x3f, 0xab, 0x4c, 0xaa, 0x4d, 0xca, 0xd0, 0x8e, 0x27, 0x99,
	0x52, 0x59, 0x0e, 0xdc, 0x0d, 0x52, 0xf3, 0xce, 0x51, 0x16, 0xa0, 0x31, 0x29, 0xaa, 0xf6, 0xef,
	0x78, 0x94, 0xa9, 0x4c, 0xb9, 0x92, 0x37, 0xd5, 0xf1, 0xf5, 0x16, 0xa1, 0x5c, 0x42, 0x5d, 0xc8,
	0x12, 0x79, 0x92, 0x2e, 0x24, 0xc7, 0x8f, 0x0a, 0x74, 0x3b, 0x9c, 0x7e, 0x12, 0x3a, 0x9c, 0x5b,
	0x01, 0xda, 0xe4, 0xe8, 0xdf, 0xd0, 0xc1, 0x0a, 0x64, 0xb6, 0xc2, 0x80, 0x44, 0x24, 0xee, 0x8a,
	0x63, 0xe7, 0x8f, 0x68, 0x5f, 0x96, 0x4b, 0xb0, 0x41, 0x27, 0x22, 0xf1, 0xa5, 0x68, 0x1b, 0xff,
	0x8a, 0x76, 0xd0, 0x06, 0xdd, 0x88, 0xc4, 0x17, 0xa2, 0x83, 0xd6, 0x7f, 0xa6, 0x83, 0xda, 0xed,
	0x09, 0x7a, 0x11, 0x89, 0xcf, 0x1f, 0xa6, 0xec, 0x3f, 0x98, 0x35, 0xc1, 0x4c, 0x80, 0xae, 0x54,
	0xa9, 0xe1, 0x05, 0x72, 0xb9, 0x86, 0x7a, 0x6e, 0x67, 0xbd, 0xed, 0xf7, 0xc4, 0x13, 0x47, 0xce,
	0x7f, 0xa2, 0xbd, 0x46, 0x29, 0xe8, 0x3b, 0x7e, 0xcc, 0x5a, 0x5f, 0x76, 0xf2, 0x65, 0xf3, 0x93,
	0xef, 0x6c, 0xd8, 0x70, 0x9b, 0x9f, 0x09, 0x11, 0x8e, 0x98, 0xbd, 0x6e, 0xf7, 0x21, 0xd9, 0xed,
	0x43, 0xf2, 0xbb, 0x0f, 0xc9, 0xe6, 0x10, 0x7a, 0xbb, 0x43, 0xe8, 0x7d, 0x1d, 0x42, 0xef, 0xed,
	0x2e, 0x93, 0xb8, 0x32, 0x29, 0x5b, 0xa8, 0x82, 0x97, 0xa6, 0x90, 0xc9, 0x32, 0xc1, 0x84, 0xa3,
	0x52, 0xb9, 0xe6, 0xee, 0xd8, 0x06, 0x65, 0xae, 0x79, 0x7b, 0xe5, 0x74, 0xe0, 0xb2, 0x1e, 0xff,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xb8, 0x32, 0xc4, 0x88, 0x01, 0x00, 0x00,
}

func (m *TxResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Time):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintTx(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x2a
	{
		size, err := m.Result.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Tx) > 0 {
		i -= len(m.Tx)
		copy(dAtA[i:], m.Tx)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Tx)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Index != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovTx(uint64(m.Height))
	}
	if m.Index != 0 {
		n += 1 + sovTx(uint64(m.Index))
	}
	l = len(m.Tx)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Result.Size()
	n += 1 + l + sovTx(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovTx(uint64(l))
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: TxResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tx", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tx = append(m.Tx[:0], dAtA[iNdEx:postIndex]...)
			if m.Tx == nil {
				m.Tx = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Result.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
