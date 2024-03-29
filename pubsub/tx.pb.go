// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosi/v1/tx.proto

package pubsub

import (
	fmt "fmt"
	types "github.com/cometbft/cometbft/abci/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
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
	Height         int64                   `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Index          uint32                  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Tx             []byte                  `protobuf:"bytes,3,opt,name=tx,proto3" json:"tx,omitempty"`
	Result         types.ResponseDeliverTx `protobuf:"bytes,4,opt,name=result,proto3" json:"result"`
	BlockTimestamp time.Time               `protobuf:"bytes,5,opt,name=block_timestamp,json=blockTimestamp,proto3,stdtime" json:"block_timestamp"`
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

func (m *TxResult) GetBlockTimestamp() time.Time {
	if m != nil {
		return m.BlockTimestamp
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*TxResult)(nil), "pubsub.tx.TxResult")
}

func init() { proto.RegisterFile("cosi/v1/tx.proto", fileDescriptor_5f4b0e93112b5f1a) }

var fileDescriptor_5f4b0e93112b5f1a = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x51, 0xbd, 0x4e, 0xeb, 0x30,
	0x14, 0x8e, 0xfb, 0xa7, 0x5e, 0xdf, 0x7b, 0x0b, 0x8a, 0x2a, 0x14, 0x15, 0x29, 0x8d, 0x3a, 0x65,
	0x40, 0xb6, 0x80, 0x17, 0x40, 0x15, 0x62, 0x63, 0xb1, 0x3a, 0xb1, 0xa0, 0x38, 0x35, 0xa9, 0x45,
	0x12, 0x47, 0xf1, 0x71, 0x65, 0xde, 0xa2, 0x8f, 0xd5, 0xb1, 0x23, 0x2c, 0x80, 0xda, 0x17, 0x41,
	0x75, 0x5a, 0xba, 0x9d, 0xe3, 0xef, 0x4f, 0x9f, 0x0f, 0x3e, 0x4f, 0x95, 0x96, 0x74, 0x79, 0x4d,
	0xc1, 0x92, 0xaa, 0x56, 0xa0, 0xfc, 0x3f, 0x95, 0xe1, 0xda, 0x70, 0x02, 0x76, 0x34, 0xce, 0x94,
	0xca, 0x72, 0x41, 0x1d, 0xc0, 0xcd, 0x0b, 0x05, 0x59, 0x08, 0x0d, 0x49, 0x51, 0x35, 0xdc, 0xd1,
	0x30, 0x53, 0x99, 0x72, 0x23, 0xdd, 0x4f, 0x87, 0xd7, 0x4b, 0x10, 0xe5, 0x5c, 0xd4, 0x85, 0x2c,
	0x81, 0x26, 0x3c, 0x95, 0x14, 0xde, 0x2a, 0xa1, 0x1b, 0x70, 0xf2, 0x81, 0x70, 0x7f, 0x66, 0x99,
	0xd0, 0x26, 0x07, 0xff, 0x02, 0xf7, 0x16, 0x42, 0x66, 0x0b, 0x08, 0x50, 0x84, 0xe2, 0x36, 0x3b,
	0x6c, 0xfe, 0x10, 0x77, 0x65, 0x39, 0x17, 0x36, 0x68, 0x45, 0x28, 0xfe, 0xcf, 0x9a, 0xc5, 0x1f,
	0xe0, 0x16, 0xd8, 0xa0, 0x1d, 0xa1, 0xf8, 0x1f, 0x6b, 0x81, 0xf5, 0xef, 0x70, 0xaf, 0x76, 0x3e,
	0x41, 0x27, 0x42, 0xf1, 0xdf, 0x9b, 0x09, 0x39, 0x05, 0x93, 0x7d, 0x30, 0x61, 0x42, 0x57, 0xaa,
	0xd4, 0xe2, 0x5e, 0xe4, 0x72, 0x29, 0xea, 0x99, 0x9d, 0x76, 0xd6, 0x9f, 0x63, 0x8f, 0x1d, 0x74,
	0xfe, 0x23, 0x3e, 0xe3, 0xb9, 0x4a, 0x5f, 0x9f, 0x7f, 0x8b, 0x05, 0x5d, 0x67, 0x35, 0x22, 0x4d,
	0x75, 0x72, 0xac, 0x4e, 0x66, 0x47, 0xc6, 0xb4, 0xbf, 0xb7, 0x58, 0x7d, 0x8d, 0x11, 0x1b, 0x38,
	0xf1, 0x09, 0x79, 0x58, 0x6f, 0x43, 0xb4, 0xd9, 0x86, 0xe8, 0x7b, 0x1b, 0xa2, 0xd5, 0x2e, 0xf4,
	0x36, 0xbb, 0xd0, 0x7b, 0xdf, 0x85, 0xde, 0xd3, 0x55, 0x26, 0x61, 0x61, 0x38, 0x49, 0x55, 0x41,
	0x4b, 0x53, 0xc8, 0x64, 0x9e, 0x40, 0x42, 0x41, 0xa9, 0x5c, 0x53, 0x77, 0x01, 0x03, 0x32, 0xd7,
	0xb4, 0xf9, 0x7a, 0xde, 0x73, 0xa9, 0xb7, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x77, 0xe1, 0x45,
	0xef, 0x9d, 0x01, 0x00, 0x00,
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
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.BlockTimestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTimestamp):])
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
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTimestamp)
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
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTimestamp", wireType)
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.BlockTimestamp, dAtA[iNdEx:postIndex]); err != nil {
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
