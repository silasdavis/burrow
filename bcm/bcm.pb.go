// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bcm.proto

package bcm

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/duration"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	github_com_hyperledger_burrow_binary "github.com/hyperledger/burrow/binary"
	io "io"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type SyncInfo struct {
	LatestBlockHeight uint64                                        `protobuf:"varint,1,opt,name=LatestBlockHeight,proto3" json:""`
	LatestBlockHash   github_com_hyperledger_burrow_binary.HexBytes `protobuf:"bytes,2,opt,name=LatestBlockHash,proto3,customtype=github.com/hyperledger/burrow/binary.HexBytes" json:"LatestBlockHash"`
	LatestAppHash     github_com_hyperledger_burrow_binary.HexBytes `protobuf:"bytes,3,opt,name=LatestAppHash,proto3,customtype=github.com/hyperledger/burrow/binary.HexBytes" json:"LatestAppHash"`
	// Timestamp of block as set by the block proposer
	LatestBlockTime time.Time `protobuf:"bytes,4,opt,name=LatestBlockTime,proto3,stdtime" json:"LatestBlockTime"`
	// Time at which we committed the last block
	LatestBlockSeenTime time.Time `protobuf:"bytes,5,opt,name=LatestBlockSeenTime,proto3,stdtime" json:"LatestBlockSeenTime"`
	// Time elapsed since last commit
	LatestBlockDuration  time.Duration `protobuf:"bytes,6,opt,name=LatestBlockDuration,proto3,stdduration" json:"LatestBlockDuration"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SyncInfo) Reset()         { *m = SyncInfo{} }
func (m *SyncInfo) String() string { return proto.CompactTextString(m) }
func (*SyncInfo) ProtoMessage()    {}
func (*SyncInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c9ff3e1ca1cc0f1, []int{0}
}
func (m *SyncInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *SyncInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncInfo.Merge(m, src)
}
func (m *SyncInfo) XXX_Size() int {
	return m.Size()
}
func (m *SyncInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SyncInfo proto.InternalMessageInfo

func (m *SyncInfo) GetLatestBlockHeight() uint64 {
	if m != nil {
		return m.LatestBlockHeight
	}
	return 0
}

func (m *SyncInfo) GetLatestBlockTime() time.Time {
	if m != nil {
		return m.LatestBlockTime
	}
	return time.Time{}
}

func (m *SyncInfo) GetLatestBlockSeenTime() time.Time {
	if m != nil {
		return m.LatestBlockSeenTime
	}
	return time.Time{}
}

func (m *SyncInfo) GetLatestBlockDuration() time.Duration {
	if m != nil {
		return m.LatestBlockDuration
	}
	return 0
}

func (*SyncInfo) XXX_MessageName() string {
	return "bcm.SyncInfo"
}

type PersistedState struct {
	AppHashAfterLastBlock github_com_hyperledger_burrow_binary.HexBytes `protobuf:"bytes,1,opt,name=AppHashAfterLastBlock,proto3,customtype=github.com/hyperledger/burrow/binary.HexBytes" json:"AppHashAfterLastBlock"`
	LastBlockTime         time.Time                                     `protobuf:"bytes,2,opt,name=LastBlockTime,proto3,stdtime" json:"LastBlockTime"`
	LastBlockHeight       uint64                                        `protobuf:"varint,3,opt,name=LastBlockHeight,proto3" json:"LastBlockHeight,omitempty"`
	GenesisHash           github_com_hyperledger_burrow_binary.HexBytes `protobuf:"bytes,4,opt,name=GenesisHash,proto3,customtype=github.com/hyperledger/burrow/binary.HexBytes" json:"GenesisHash"`
	XXX_NoUnkeyedLiteral  struct{}                                      `json:"-"`
	XXX_unrecognized      []byte                                        `json:"-"`
	XXX_sizecache         int32                                         `json:"-"`
}

func (m *PersistedState) Reset()         { *m = PersistedState{} }
func (m *PersistedState) String() string { return proto.CompactTextString(m) }
func (*PersistedState) ProtoMessage()    {}
func (*PersistedState) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c9ff3e1ca1cc0f1, []int{1}
}
func (m *PersistedState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PersistedState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *PersistedState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PersistedState.Merge(m, src)
}
func (m *PersistedState) XXX_Size() int {
	return m.Size()
}
func (m *PersistedState) XXX_DiscardUnknown() {
	xxx_messageInfo_PersistedState.DiscardUnknown(m)
}

var xxx_messageInfo_PersistedState proto.InternalMessageInfo

func (m *PersistedState) GetLastBlockTime() time.Time {
	if m != nil {
		return m.LastBlockTime
	}
	return time.Time{}
}

func (m *PersistedState) GetLastBlockHeight() uint64 {
	if m != nil {
		return m.LastBlockHeight
	}
	return 0
}

func (*PersistedState) XXX_MessageName() string {
	return "bcm.PersistedState"
}
func init() {
	proto.RegisterType((*SyncInfo)(nil), "bcm.SyncInfo")
	golang_proto.RegisterType((*SyncInfo)(nil), "bcm.SyncInfo")
	proto.RegisterType((*PersistedState)(nil), "bcm.PersistedState")
	golang_proto.RegisterType((*PersistedState)(nil), "bcm.PersistedState")
}

func init() { proto.RegisterFile("bcm.proto", fileDescriptor_0c9ff3e1ca1cc0f1) }
func init() { golang_proto.RegisterFile("bcm.proto", fileDescriptor_0c9ff3e1ca1cc0f1) }

var fileDescriptor_0c9ff3e1ca1cc0f1 = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x3d, 0x8f, 0x94, 0x40,
	0x18, 0xc7, 0x6f, 0x76, 0xf1, 0x72, 0xce, 0xf9, 0x12, 0xc7, 0x98, 0xe0, 0x16, 0xb0, 0x5e, 0x45,
	0x23, 0x24, 0x67, 0xac, 0xac, 0x8e, 0x98, 0x78, 0x9a, 0x8b, 0x31, 0xec, 0xa9, 0x89, 0x16, 0x06,
	0xd8, 0x67, 0x61, 0x72, 0x0b, 0x43, 0x66, 0x86, 0x28, 0xdf, 0xc2, 0xd2, 0x8f, 0x63, 0xb9, 0x85,
	0x85, 0xa5, 0xb1, 0x58, 0x0d, 0xd7, 0xf9, 0x15, 0x6c, 0x0c, 0x03, 0x44, 0xe0, 0x2e, 0x31, 0xeb,
	0x76, 0x3c, 0x6f, 0x3f, 0xe6, 0xf9, 0x3f, 0x7f, 0x7c, 0x35, 0x08, 0x13, 0x3b, 0xe3, 0x4c, 0x32,
	0x32, 0x0e, 0xc2, 0x64, 0x72, 0x3f, 0xa2, 0x32, 0xce, 0x03, 0x3b, 0x64, 0x89, 0x13, 0xb1, 0x88,
	0x39, 0xaa, 0x16, 0xe4, 0x0b, 0x15, 0xa9, 0x40, 0x7d, 0xd5, 0x33, 0x13, 0x33, 0x62, 0x2c, 0x5a,
	0xc2, 0xdf, 0x2e, 0x49, 0x13, 0x10, 0xd2, 0x4f, 0xb2, 0xa6, 0xc1, 0x18, 0x36, 0xcc, 0x73, 0xee,
	0x4b, 0xca, 0xd2, 0xba, 0x7e, 0xf0, 0x7b, 0x8c, 0xf7, 0x66, 0x45, 0x1a, 0x3e, 0x4d, 0x17, 0x8c,
	0x1c, 0xe2, 0x5b, 0x27, 0xbe, 0x04, 0x21, 0xdd, 0x25, 0x0b, 0xcf, 0x8e, 0x81, 0x46, 0xb1, 0xd4,
	0xd1, 0x14, 0x59, 0x9a, 0xab, 0xfd, 0x5a, 0x9b, 0x3b, 0xde, 0xc5, 0x32, 0x79, 0x87, 0x6f, 0x76,
	0x93, 0xbe, 0x88, 0xf5, 0xd1, 0x14, 0x59, 0xd7, 0xdc, 0x87, 0xab, 0xb5, 0xb9, 0xf3, 0x7d, 0x6d,
	0x76, 0x37, 0x8a, 0x8b, 0x0c, 0xf8, 0x12, 0xe6, 0x11, 0x70, 0x27, 0xc8, 0x39, 0x67, 0xef, 0x9d,
	0x80, 0xa6, 0x3e, 0x2f, 0xec, 0x63, 0xf8, 0xe0, 0x16, 0x12, 0x84, 0x37, 0xa4, 0x91, 0xb7, 0xf8,
	0x7a, 0x9d, 0x3a, 0xca, 0x32, 0x85, 0x1f, 0x6f, 0x83, 0xef, 0xb3, 0xc8, 0xf3, 0xde, 0xeb, 0x4f,
	0x69, 0x02, 0xba, 0x36, 0x45, 0xd6, 0xfe, 0xe1, 0xc4, 0xae, 0x85, 0xb3, 0x5b, 0xe1, 0xec, 0xd3,
	0x56, 0x59, 0x77, 0xaf, 0xfa, 0xf5, 0xc7, 0x1f, 0x26, 0xf2, 0x86, 0xc3, 0xe4, 0x15, 0xbe, 0xdd,
	0x49, 0xcd, 0x00, 0x52, 0xc5, 0xbc, 0xb2, 0x01, 0xf3, 0x32, 0x00, 0x79, 0xd9, 0xe3, 0x3e, 0x6e,
	0x6e, 0xa8, 0xef, 0x2a, 0xee, 0xdd, 0x0b, 0xdc, 0xb6, 0xa1, 0xc6, 0x7e, 0x1a, 0x62, 0xdb, 0xf2,
	0xc1, 0x97, 0x11, 0xbe, 0xf1, 0x02, 0xb8, 0xa0, 0x42, 0xc2, 0x7c, 0x26, 0x7d, 0x09, 0xe4, 0x0c,
	0xdf, 0x69, 0xc4, 0x39, 0x5a, 0x48, 0xe0, 0x27, 0x7e, 0x33, 0xa3, 0x7c, 0xf0, 0xdf, 0xb2, 0x5f,
	0xce, 0x24, 0xcf, 0xaa, 0xdb, 0x76, 0xc5, 0x1f, 0x6d, 0x20, 0x54, 0x7f, 0x94, 0x58, 0xd5, 0x29,
	0xfb, 0xd6, 0xad, 0x9c, 0xa2, 0x79, 0xc3, 0x34, 0x79, 0x8d, 0xf7, 0x9f, 0x40, 0x0a, 0x82, 0x0a,
	0xe5, 0x27, 0x6d, 0x9b, 0xc5, 0xba, 0x24, 0xf7, 0xd1, 0xaa, 0x34, 0xd0, 0xd7, 0xd2, 0x40, 0xdf,
	0x4a, 0x03, 0xfd, 0x2c, 0x0d, 0xf4, 0xf9, 0xdc, 0x40, 0xab, 0x73, 0x03, 0xbd, 0xb9, 0xf7, 0x0f,
	0x6a, 0x98, 0x04, 0xbb, 0x6a, 0xd9, 0x07, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x25, 0xfd, 0xe5,
	0xbc, 0x12, 0x04, 0x00, 0x00,
}

func (m *SyncInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.LatestBlockHeight != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintBcm(dAtA, i, uint64(m.LatestBlockHeight))
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintBcm(dAtA, i, uint64(m.LatestBlockHash.Size()))
	n1, err := m.LatestBlockHash.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x1a
	i++
	i = encodeVarintBcm(dAtA, i, uint64(m.LatestAppHash.Size()))
	n2, err := m.LatestAppHash.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x22
	i++
	i = encodeVarintBcm(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.LatestBlockTime)))
	n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.LatestBlockTime, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	dAtA[i] = 0x2a
	i++
	i = encodeVarintBcm(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.LatestBlockSeenTime)))
	n4, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.LatestBlockSeenTime, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	dAtA[i] = 0x32
	i++
	i = encodeVarintBcm(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdDuration(m.LatestBlockDuration)))
	n5, err := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.LatestBlockDuration, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *PersistedState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PersistedState) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintBcm(dAtA, i, uint64(m.AppHashAfterLastBlock.Size()))
	n6, err := m.AppHashAfterLastBlock.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	dAtA[i] = 0x12
	i++
	i = encodeVarintBcm(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.LastBlockTime)))
	n7, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.LastBlockTime, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n7
	if m.LastBlockHeight != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintBcm(dAtA, i, uint64(m.LastBlockHeight))
	}
	dAtA[i] = 0x22
	i++
	i = encodeVarintBcm(dAtA, i, uint64(m.GenesisHash.Size()))
	n8, err := m.GenesisHash.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n8
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintBcm(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *SyncInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.LatestBlockHeight != 0 {
		n += 1 + sovBcm(uint64(m.LatestBlockHeight))
	}
	l = m.LatestBlockHash.Size()
	n += 1 + l + sovBcm(uint64(l))
	l = m.LatestAppHash.Size()
	n += 1 + l + sovBcm(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.LatestBlockTime)
	n += 1 + l + sovBcm(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.LatestBlockSeenTime)
	n += 1 + l + sovBcm(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.LatestBlockDuration)
	n += 1 + l + sovBcm(uint64(l))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PersistedState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.AppHashAfterLastBlock.Size()
	n += 1 + l + sovBcm(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.LastBlockTime)
	n += 1 + l + sovBcm(uint64(l))
	if m.LastBlockHeight != 0 {
		n += 1 + sovBcm(uint64(m.LastBlockHeight))
	}
	l = m.GenesisHash.Size()
	n += 1 + l + sovBcm(uint64(l))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBcm(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBcm(x uint64) (n int) {
	return sovBcm(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SyncInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBcm
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
			return fmt.Errorf("proto: SyncInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestBlockHeight", wireType)
			}
			m.LatestBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LatestBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestBlockHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LatestBlockHash.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestAppHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LatestAppHash.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestBlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.LatestBlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestBlockSeenTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.LatestBlockSeenTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestBlockDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.LatestBlockDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBcm(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBcm
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBcm
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
func (m *PersistedState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBcm
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
			return fmt.Errorf("proto: PersistedState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PersistedState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppHashAfterLastBlock", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AppHashAfterLastBlock.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastBlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.LastBlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastBlockHeight", wireType)
			}
			m.LastBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBcm
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
				return ErrInvalidLengthBcm
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBcm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GenesisHash.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBcm(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBcm
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBcm
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
func skipBcm(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBcm
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
					return 0, ErrIntOverflowBcm
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBcm
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
				return 0, ErrInvalidLengthBcm
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthBcm
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBcm
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipBcm(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthBcm
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthBcm = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBcm   = fmt.Errorf("proto: integer overflow")
)
