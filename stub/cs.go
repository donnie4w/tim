// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package stub

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	thrift "github.com/donnie4w/gothrift/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = errors.New
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal
// (needed by validator.)
var _ = strings.Contains
var _ = regexp.MatchString

// Attributes:
//  - Rtype
//  - Vnode
//  - FoundNode
//  - Rnode
//  - Body
//  - Dtype
//  - StreamId
type VBean struct {
  Rtype int8 `thrift:"rtype,1,required" db:"rtype" json:"rtype"`
  Vnode string `thrift:"vnode,2,required" db:"vnode" json:"vnode"`
  FoundNode *string `thrift:"foundNode,3" db:"foundNode" json:"foundNode,omitempty"`
  Rnode *string `thrift:"rnode,4" db:"rnode" json:"rnode,omitempty"`
  Body []byte `thrift:"body,5" db:"body" json:"body,omitempty"`
  Dtype *int8 `thrift:"dtype,6" db:"dtype" json:"dtype,omitempty"`
  StreamId *int64 `thrift:"streamId,7" db:"streamId" json:"streamId,omitempty"`
}

func NewVBean() *VBean {
  return &VBean{}
}


func (p *VBean) GetRtype() int8 {
  return p.Rtype
}

func (p *VBean) GetVnode() string {
  return p.Vnode
}
var VBean_FoundNode_DEFAULT string
func (p *VBean) GetFoundNode() string {
  if !p.IsSetFoundNode() {
    return VBean_FoundNode_DEFAULT
  }
return *p.FoundNode
}
var VBean_Rnode_DEFAULT string
func (p *VBean) GetRnode() string {
  if !p.IsSetRnode() {
    return VBean_Rnode_DEFAULT
  }
return *p.Rnode
}
var VBean_Body_DEFAULT []byte

func (p *VBean) GetBody() []byte {
  return p.Body
}
var VBean_Dtype_DEFAULT int8
func (p *VBean) GetDtype() int8 {
  if !p.IsSetDtype() {
    return VBean_Dtype_DEFAULT
  }
return *p.Dtype
}
var VBean_StreamId_DEFAULT int64
func (p *VBean) GetStreamId() int64 {
  if !p.IsSetStreamId() {
    return VBean_StreamId_DEFAULT
  }
return *p.StreamId
}
func (p *VBean) IsSetFoundNode() bool {
  return p.FoundNode != nil
}

func (p *VBean) IsSetRnode() bool {
  return p.Rnode != nil
}

func (p *VBean) IsSetBody() bool {
  return p.Body != nil
}

func (p *VBean) IsSetDtype() bool {
  return p.Dtype != nil
}

func (p *VBean) IsSetStreamId() bool {
  return p.StreamId != nil
}

func (p *VBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetRtype bool = false;
  var issetVnode bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetRtype = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetVnode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 7:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField7(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetRtype{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Rtype is not set"));
  }
  if !issetVnode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Vnode is not set"));
  }
  return nil
}

func (p *VBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.Rtype = temp
}
  return nil
}

func (p *VBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Vnode = v
}
  return nil
}

func (p *VBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.FoundNode = &v
}
  return nil
}

func (p *VBean)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Rnode = &v
}
  return nil
}

func (p *VBean)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.Body = v
}
  return nil
}

func (p *VBean)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  temp := int8(v)
  p.Dtype = &temp
}
  return nil
}

func (p *VBean)  ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 7: ", err)
} else {
  p.StreamId = &v
}
  return nil
}

func (p *VBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "VBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
    if err := p.writeField7(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *VBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "rtype", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:rtype: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Rtype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.rtype (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:rtype: ", p), err) }
  return err
}

func (p *VBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "vnode", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:vnode: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Vnode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.vnode (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:vnode: ", p), err) }
  return err
}

func (p *VBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetFoundNode() {
    if err := oprot.WriteFieldBegin(ctx, "foundNode", thrift.STRING, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:foundNode: ", p), err) }
    if err := oprot.WriteString(ctx, string(*p.FoundNode)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.foundNode (3) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:foundNode: ", p), err) }
  }
  return err
}

func (p *VBean) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetRnode() {
    if err := oprot.WriteFieldBegin(ctx, "rnode", thrift.STRING, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:rnode: ", p), err) }
    if err := oprot.WriteString(ctx, string(*p.Rnode)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.rnode (4) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:rnode: ", p), err) }
  }
  return err
}

func (p *VBean) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBody() {
    if err := oprot.WriteFieldBegin(ctx, "body", thrift.STRING, 5); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:body: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Body); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.body (5) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 5:body: ", p), err) }
  }
  return err
}

func (p *VBean) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetDtype() {
    if err := oprot.WriteFieldBegin(ctx, "dtype", thrift.BYTE, 6); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:dtype: ", p), err) }
    if err := oprot.WriteByte(ctx, int8(*p.Dtype)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.dtype (6) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 6:dtype: ", p), err) }
  }
  return err
}

func (p *VBean) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetStreamId() {
    if err := oprot.WriteFieldBegin(ctx, "streamId", thrift.I64, 7); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:streamId: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.StreamId)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.streamId (7) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 7:streamId: ", p), err) }
  }
  return err
}

func (p *VBean) Equals(other *VBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Rtype != other.Rtype { return false }
  if p.Vnode != other.Vnode { return false }
  if p.FoundNode != other.FoundNode {
    if p.FoundNode == nil || other.FoundNode == nil {
      return false
    }
    if (*p.FoundNode) != (*other.FoundNode) { return false }
  }
  if p.Rnode != other.Rnode {
    if p.Rnode == nil || other.Rnode == nil {
      return false
    }
    if (*p.Rnode) != (*other.Rnode) { return false }
  }
  if bytes.Compare(p.Body, other.Body) != 0 { return false }
  if p.Dtype != other.Dtype {
    if p.Dtype == nil || other.Dtype == nil {
      return false
    }
    if (*p.Dtype) != (*other.Dtype) { return false }
  }
  if p.StreamId != other.StreamId {
    if p.StreamId == nil || other.StreamId == nil {
      return false
    }
    if (*p.StreamId) != (*other.StreamId) { return false }
  }
  return true
}

func (p *VBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("VBean(%+v)", *p)
}

func (p *VBean) Validate() error {
  return nil
}
// Attributes:
//  - ID
//  - VNode
//  - Dtype
//  - Body
//  - FNode
type TimPushStream struct {
  ID int32 `thrift:"id,1,required" db:"id" json:"id"`
  VNode string `thrift:"VNode,2,required" db:"VNode" json:"VNode"`
  Dtype *int8 `thrift:"dtype,3" db:"dtype" json:"dtype,omitempty"`
  Body []byte `thrift:"body,4" db:"body" json:"body,omitempty"`
  FNode string `thrift:"fNode,5,required" db:"fNode" json:"fNode"`
}

func NewTimPushStream() *TimPushStream {
  return &TimPushStream{}
}


func (p *TimPushStream) GetID() int32 {
  return p.ID
}

func (p *TimPushStream) GetVNode() string {
  return p.VNode
}
var TimPushStream_Dtype_DEFAULT int8
func (p *TimPushStream) GetDtype() int8 {
  if !p.IsSetDtype() {
    return TimPushStream_Dtype_DEFAULT
  }
return *p.Dtype
}
var TimPushStream_Body_DEFAULT []byte

func (p *TimPushStream) GetBody() []byte {
  return p.Body
}

func (p *TimPushStream) GetFNode() string {
  return p.FNode
}
func (p *TimPushStream) IsSetDtype() bool {
  return p.Dtype != nil
}

func (p *TimPushStream) IsSetBody() bool {
  return p.Body != nil
}

func (p *TimPushStream) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetID bool = false;
  var issetVNode bool = false;
  var issetFNode bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetID = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetVNode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
        issetFNode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetID{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ID is not set"));
  }
  if !issetVNode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field VNode is not set"));
  }
  if !issetFNode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field FNode is not set"));
  }
  return nil
}

func (p *TimPushStream)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.ID = v
}
  return nil
}

func (p *TimPushStream)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.VNode = v
}
  return nil
}

func (p *TimPushStream)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.Dtype = &temp
}
  return nil
}

func (p *TimPushStream)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Body = v
}
  return nil
}

func (p *TimPushStream)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.FNode = v
}
  return nil
}

func (p *TimPushStream) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "TimPushStream"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TimPushStream) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "id", thrift.I32, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:id: ", p), err) }
  if err := oprot.WriteI32(ctx, int32(p.ID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:id: ", p), err) }
  return err
}

func (p *TimPushStream) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "VNode", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:VNode: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.VNode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.VNode (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:VNode: ", p), err) }
  return err
}

func (p *TimPushStream) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetDtype() {
    if err := oprot.WriteFieldBegin(ctx, "dtype", thrift.BYTE, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:dtype: ", p), err) }
    if err := oprot.WriteByte(ctx, int8(*p.Dtype)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.dtype (3) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:dtype: ", p), err) }
  }
  return err
}

func (p *TimPushStream) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBody() {
    if err := oprot.WriteFieldBegin(ctx, "body", thrift.STRING, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:body: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Body); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.body (4) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:body: ", p), err) }
  }
  return err
}

func (p *TimPushStream) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "fNode", thrift.STRING, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:fNode: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.FNode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.fNode (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:fNode: ", p), err) }
  return err
}

func (p *TimPushStream) Equals(other *TimPushStream) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.ID != other.ID { return false }
  if p.VNode != other.VNode { return false }
  if p.Dtype != other.Dtype {
    if p.Dtype == nil || other.Dtype == nil {
      return false
    }
    if (*p.Dtype) != (*other.Dtype) { return false }
  }
  if bytes.Compare(p.Body, other.Body) != 0 { return false }
  if p.FNode != other.FNode { return false }
  return true
}

func (p *TimPushStream) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TimPushStream(%+v)", *p)
}

func (p *TimPushStream) Validate() error {
  return nil
}
// Attributes:
//  - FNode
//  - Vnode
//  - Otype
type TimReqVRoom struct {
  FNode *string `thrift:"fNode,1" db:"fNode" json:"fNode,omitempty"`
  Vnode string `thrift:"vnode,2,required" db:"vnode" json:"vnode"`
  Otype int8 `thrift:"otype,3,required" db:"otype" json:"otype"`
}

func NewTimReqVRoom() *TimReqVRoom {
  return &TimReqVRoom{}
}

var TimReqVRoom_FNode_DEFAULT string
func (p *TimReqVRoom) GetFNode() string {
  if !p.IsSetFNode() {
    return TimReqVRoom_FNode_DEFAULT
  }
return *p.FNode
}

func (p *TimReqVRoom) GetVnode() string {
  return p.Vnode
}

func (p *TimReqVRoom) GetOtype() int8 {
  return p.Otype
}
func (p *TimReqVRoom) IsSetFNode() bool {
  return p.FNode != nil
}

func (p *TimReqVRoom) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetVnode bool = false;
  var issetOtype bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetVnode = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetOtype = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetVnode{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Vnode is not set"));
  }
  if !issetOtype{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Otype is not set"));
  }
  return nil
}

func (p *TimReqVRoom)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.FNode = &v
}
  return nil
}

func (p *TimReqVRoom)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Vnode = v
}
  return nil
}

func (p *TimReqVRoom)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.Otype = temp
}
  return nil
}

func (p *TimReqVRoom) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "TimReqVRoom"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TimReqVRoom) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetFNode() {
    if err := oprot.WriteFieldBegin(ctx, "fNode", thrift.STRING, 1); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:fNode: ", p), err) }
    if err := oprot.WriteString(ctx, string(*p.FNode)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.fNode (1) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 1:fNode: ", p), err) }
  }
  return err
}

func (p *TimReqVRoom) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "vnode", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:vnode: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Vnode)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.vnode (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:vnode: ", p), err) }
  return err
}

func (p *TimReqVRoom) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "otype", thrift.BYTE, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:otype: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Otype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.otype (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:otype: ", p), err) }
  return err
}

func (p *TimReqVRoom) Equals(other *TimReqVRoom) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.FNode != other.FNode {
    if p.FNode == nil || other.FNode == nil {
      return false
    }
    if (*p.FNode) != (*other.FNode) { return false }
  }
  if p.Vnode != other.Vnode { return false }
  if p.Otype != other.Otype { return false }
  return true
}

func (p *TimReqVRoom) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TimReqVRoom(%+v)", *p)
}

func (p *TimReqVRoom) Validate() error {
  return nil
}
// Attributes:
//  - Addr
//  - UUID
//  - Nodekv
type Node struct {
  Addr string `thrift:"addr,1,required" db:"addr" json:"addr"`
  UUID int64 `thrift:"uuid,2,required" db:"uuid" json:"uuid"`
  Nodekv map[int64]string `thrift:"nodekv,3" db:"nodekv" json:"nodekv,omitempty"`
}

func NewNode() *Node {
  return &Node{}
}


func (p *Node) GetAddr() string {
  return p.Addr
}

func (p *Node) GetUUID() int64 {
  return p.UUID
}
var Node_Nodekv_DEFAULT map[int64]string

func (p *Node) GetNodekv() map[int64]string {
  return p.Nodekv
}
func (p *Node) IsSetNodekv() bool {
  return p.Nodekv != nil
}

func (p *Node) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetAddr bool = false;
  var issetUUID bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetAddr = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetUUID = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetAddr{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Addr is not set"));
  }
  if !issetUUID{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field UUID is not set"));
  }
  return nil
}

func (p *Node)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Addr = v
}
  return nil
}

func (p *Node)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.UUID = v
}
  return nil
}

func (p *Node)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[int64]string, size)
  p.Nodekv =  tMap
  for i := 0; i < size; i ++ {
var _key0 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key0 = v
}
var _val1 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val1 = v
}
    p.Nodekv[_key0] = _val1
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *Node) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Node"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *Node) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "addr", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:addr: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Addr)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.addr (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:addr: ", p), err) }
  return err
}

func (p *Node) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "uuid", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:uuid: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.UUID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.uuid (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:uuid: ", p), err) }
  return err
}

func (p *Node) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetNodekv() {
    if err := oprot.WriteFieldBegin(ctx, "nodekv", thrift.MAP, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:nodekv: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.I64, thrift.STRING, len(p.Nodekv)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Nodekv {
      if err := oprot.WriteI64(ctx, int64(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteString(ctx, string(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:nodekv: ", p), err) }
  }
  return err
}

func (p *Node) Equals(other *Node) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Addr != other.Addr { return false }
  if p.UUID != other.UUID { return false }
  if len(p.Nodekv) != len(other.Nodekv) { return false }
  for k, _tgt := range p.Nodekv {
    _src2 := other.Nodekv[k]
    if _tgt != _src2 { return false }
  }
  return true
}

func (p *Node) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Node(%+v)", *p)
}

func (p *Node) Validate() error {
  return nil
}
// Attributes:
//  - SyncNum
//  - OnNum
//  - Bytes
//  - CsNum
type Data struct {
  SyncNum *int64 `thrift:"syncNum,1" db:"syncNum" json:"syncNum,omitempty"`
  OnNum *int64 `thrift:"onNum,2" db:"onNum" json:"onNum,omitempty"`
  Bytes []byte `thrift:"bytes,3" db:"bytes" json:"bytes,omitempty"`
  CsNum *int32 `thrift:"csNum,4" db:"csNum" json:"csNum,omitempty"`
}

func NewData() *Data {
  return &Data{}
}

var Data_SyncNum_DEFAULT int64
func (p *Data) GetSyncNum() int64 {
  if !p.IsSetSyncNum() {
    return Data_SyncNum_DEFAULT
  }
return *p.SyncNum
}
var Data_OnNum_DEFAULT int64
func (p *Data) GetOnNum() int64 {
  if !p.IsSetOnNum() {
    return Data_OnNum_DEFAULT
  }
return *p.OnNum
}
var Data_Bytes_DEFAULT []byte

func (p *Data) GetBytes() []byte {
  return p.Bytes
}
var Data_CsNum_DEFAULT int32
func (p *Data) GetCsNum() int32 {
  if !p.IsSetCsNum() {
    return Data_CsNum_DEFAULT
  }
return *p.CsNum
}
func (p *Data) IsSetSyncNum() bool {
  return p.SyncNum != nil
}

func (p *Data) IsSetOnNum() bool {
  return p.OnNum != nil
}

func (p *Data) IsSetBytes() bool {
  return p.Bytes != nil
}

func (p *Data) IsSetCsNum() bool {
  return p.CsNum != nil
}

func (p *Data) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *Data)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SyncNum = &v
}
  return nil
}

func (p *Data)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.OnNum = &v
}
  return nil
}

func (p *Data)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Bytes = v
}
  return nil
}

func (p *Data)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.CsNum = &v
}
  return nil
}

func (p *Data) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Data"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *Data) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetSyncNum() {
    if err := oprot.WriteFieldBegin(ctx, "syncNum", thrift.I64, 1); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncNum: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.SyncNum)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.syncNum (1) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncNum: ", p), err) }
  }
  return err
}

func (p *Data) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetOnNum() {
    if err := oprot.WriteFieldBegin(ctx, "onNum", thrift.I64, 2); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:onNum: ", p), err) }
    if err := oprot.WriteI64(ctx, int64(*p.OnNum)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.onNum (2) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 2:onNum: ", p), err) }
  }
  return err
}

func (p *Data) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBytes() {
    if err := oprot.WriteFieldBegin(ctx, "bytes", thrift.STRING, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:bytes: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Bytes); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.bytes (3) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:bytes: ", p), err) }
  }
  return err
}

func (p *Data) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetCsNum() {
    if err := oprot.WriteFieldBegin(ctx, "csNum", thrift.I32, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:csNum: ", p), err) }
    if err := oprot.WriteI32(ctx, int32(*p.CsNum)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.csNum (4) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:csNum: ", p), err) }
  }
  return err
}

func (p *Data) Equals(other *Data) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.SyncNum != other.SyncNum {
    if p.SyncNum == nil || other.SyncNum == nil {
      return false
    }
    if (*p.SyncNum) != (*other.SyncNum) { return false }
  }
  if p.OnNum != other.OnNum {
    if p.OnNum == nil || other.OnNum == nil {
      return false
    }
    if (*p.OnNum) != (*other.OnNum) { return false }
  }
  if bytes.Compare(p.Bytes, other.Bytes) != 0 { return false }
  if p.CsNum != other.CsNum {
    if p.CsNum == nil || other.CsNum == nil {
      return false
    }
    if (*p.CsNum) != (*other.CsNum) { return false }
  }
  return true
}

func (p *Data) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Data(%+v)", *p)
}

func (p *Data) Validate() error {
  return nil
}
// Attributes:
//  - Stat
//  - Node
//  - BkNode
type CsUser struct {
  Stat int8 `thrift:"stat,1,required" db:"stat" json:"stat"`
  Node map[string]int8 `thrift:"node,2" db:"node" json:"node,omitempty"`
  BkNode map[string]int64 `thrift:"bkNode,3" db:"bkNode" json:"bkNode,omitempty"`
}

func NewCsUser() *CsUser {
  return &CsUser{}
}


func (p *CsUser) GetStat() int8 {
  return p.Stat
}
var CsUser_Node_DEFAULT map[string]int8

func (p *CsUser) GetNode() map[string]int8 {
  return p.Node
}
var CsUser_BkNode_DEFAULT map[string]int64

func (p *CsUser) GetBkNode() map[string]int64 {
  return p.BkNode
}
func (p *CsUser) IsSetNode() bool {
  return p.Node != nil
}

func (p *CsUser) IsSetBkNode() bool {
  return p.BkNode != nil
}

func (p *CsUser) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetStat bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetStat = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetStat{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Stat is not set"));
  }
  return nil
}

func (p *CsUser)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.Stat = temp
}
  return nil
}

func (p *CsUser)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]int8, size)
  p.Node =  tMap
  for i := 0; i < size; i ++ {
var _key3 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key3 = v
}
var _val4 int8
    if v, err := iprot.ReadByte(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    temp := int8(v)
    _val4 = temp
}
    p.Node[_key3] = _val4
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *CsUser)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]int64, size)
  p.BkNode =  tMap
  for i := 0; i < size; i ++ {
var _key5 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key5 = v
}
var _val6 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val6 = v
}
    p.BkNode[_key5] = _val6
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *CsUser) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsUser"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *CsUser) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "stat", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:stat: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.Stat)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.stat (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:stat: ", p), err) }
  return err
}

func (p *CsUser) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetNode() {
    if err := oprot.WriteFieldBegin(ctx, "node", thrift.MAP, 2); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:node: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.BYTE, len(p.Node)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Node {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteByte(ctx, int8(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 2:node: ", p), err) }
  }
  return err
}

func (p *CsUser) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBkNode() {
    if err := oprot.WriteFieldBegin(ctx, "bkNode", thrift.MAP, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:bkNode: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.I64, len(p.BkNode)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.BkNode {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteI64(ctx, int64(v)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:bkNode: ", p), err) }
  }
  return err
}

func (p *CsUser) Equals(other *CsUser) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Stat != other.Stat { return false }
  if len(p.Node) != len(other.Node) { return false }
  for k, _tgt := range p.Node {
    _src7 := other.Node[k]
    if _tgt != _src7 { return false }
  }
  if len(p.BkNode) != len(other.BkNode) { return false }
  for k, _tgt := range p.BkNode {
    _src8 := other.BkNode[k]
    if _tgt != _src8 { return false }
  }
  return true
}

func (p *CsUser) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("CsUser(%+v)", *p)
}

func (p *CsUser) Validate() error {
  return nil
}
// Attributes:
//  - TransType
//  - BsType
//  - MType
//  - Bs
//  - Node
//  - Cache
type CsBs struct {
  TransType int8 `thrift:"transType,1,required" db:"transType" json:"transType"`
  BsType int8 `thrift:"bsType,2,required" db:"bsType" json:"bsType"`
  MType int8 `thrift:"mType,3,required" db:"mType" json:"mType"`
  Bs []byte `thrift:"bs,4" db:"bs" json:"bs,omitempty"`
  Node []byte `thrift:"node,5" db:"node" json:"node,omitempty"`
  Cache *bool `thrift:"cache,6" db:"cache" json:"cache,omitempty"`
}

func NewCsBs() *CsBs {
  return &CsBs{}
}


func (p *CsBs) GetTransType() int8 {
  return p.TransType
}

func (p *CsBs) GetBsType() int8 {
  return p.BsType
}

func (p *CsBs) GetMType() int8 {
  return p.MType
}
var CsBs_Bs_DEFAULT []byte

func (p *CsBs) GetBs() []byte {
  return p.Bs
}
var CsBs_Node_DEFAULT []byte

func (p *CsBs) GetNode() []byte {
  return p.Node
}
var CsBs_Cache_DEFAULT bool
func (p *CsBs) GetCache() bool {
  if !p.IsSetCache() {
    return CsBs_Cache_DEFAULT
  }
return *p.Cache
}
func (p *CsBs) IsSetBs() bool {
  return p.Bs != nil
}

func (p *CsBs) IsSetNode() bool {
  return p.Node != nil
}

func (p *CsBs) IsSetCache() bool {
  return p.Cache != nil
}

func (p *CsBs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetTransType bool = false;
  var issetBsType bool = false;
  var issetMType bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetTransType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
        issetBsType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.BYTE {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
        issetMType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField4(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField5(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField6(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetTransType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TransType is not set"));
  }
  if !issetBsType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field BsType is not set"));
  }
  if !issetMType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field MType is not set"));
  }
  return nil
}

func (p *CsBs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  temp := int8(v)
  p.TransType = temp
}
  return nil
}

func (p *CsBs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  temp := int8(v)
  p.BsType = temp
}
  return nil
}

func (p *CsBs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadByte(ctx); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  temp := int8(v)
  p.MType = temp
}
  return nil
}

func (p *CsBs)  ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Bs = v
}
  return nil
}

func (p *CsBs)  ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.Node = v
}
  return nil
}

func (p *CsBs)  ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.Cache = &v
}
  return nil
}

func (p *CsBs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsBs"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
    if err := p.writeField4(ctx, oprot); err != nil { return err }
    if err := p.writeField5(ctx, oprot); err != nil { return err }
    if err := p.writeField6(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *CsBs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "transType", thrift.BYTE, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:transType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.TransType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.transType (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:transType: ", p), err) }
  return err
}

func (p *CsBs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "bsType", thrift.BYTE, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:bsType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.BsType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.bsType (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:bsType: ", p), err) }
  return err
}

func (p *CsBs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "mType", thrift.BYTE, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:mType: ", p), err) }
  if err := oprot.WriteByte(ctx, int8(p.MType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.mType (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:mType: ", p), err) }
  return err
}

func (p *CsBs) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBs() {
    if err := oprot.WriteFieldBegin(ctx, "bs", thrift.STRING, 4); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:bs: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Bs); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.bs (4) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 4:bs: ", p), err) }
  }
  return err
}

func (p *CsBs) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetNode() {
    if err := oprot.WriteFieldBegin(ctx, "node", thrift.STRING, 5); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:node: ", p), err) }
    if err := oprot.WriteBinary(ctx, p.Node); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.node (5) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 5:node: ", p), err) }
  }
  return err
}

func (p *CsBs) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetCache() {
    if err := oprot.WriteFieldBegin(ctx, "cache", thrift.BOOL, 6); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:cache: ", p), err) }
    if err := oprot.WriteBool(ctx, bool(*p.Cache)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.cache (6) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 6:cache: ", p), err) }
  }
  return err
}

func (p *CsBs) Equals(other *CsBs) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.TransType != other.TransType { return false }
  if p.BsType != other.BsType { return false }
  if p.MType != other.MType { return false }
  if bytes.Compare(p.Bs, other.Bs) != 0 { return false }
  if bytes.Compare(p.Node, other.Node) != 0 { return false }
  if p.Cache != other.Cache {
    if p.Cache == nil || other.Cache == nil {
      return false
    }
    if (*p.Cache) != (*other.Cache) { return false }
  }
  return true
}

func (p *CsBs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("CsBs(%+v)", *p)
}

func (p *CsBs) Validate() error {
  return nil
}
// Attributes:
//  - RType
//  - Bsm
//  - Bsm2
type CsBean struct {
  RType int16 `thrift:"rType,1,required" db:"rType" json:"rType"`
  Bsm map[string][]byte `thrift:"bsm,2" db:"bsm" json:"bsm,omitempty"`
  Bsm2 map[string][]int64 `thrift:"bsm2,3" db:"bsm2" json:"bsm2,omitempty"`
}

func NewCsBean() *CsBean {
  return &CsBean{}
}


func (p *CsBean) GetRType() int16 {
  return p.RType
}
var CsBean_Bsm_DEFAULT map[string][]byte

func (p *CsBean) GetBsm() map[string][]byte {
  return p.Bsm
}
var CsBean_Bsm2_DEFAULT map[string][]int64

func (p *CsBean) GetBsm2() map[string][]int64 {
  return p.Bsm2
}
func (p *CsBean) IsSetBsm() bool {
  return p.Bsm != nil
}

func (p *CsBean) IsSetBsm2() bool {
  return p.Bsm2 != nil
}

func (p *CsBean) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }

  var issetRType bool = false;

  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I16 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
        issetRType = true
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  if !issetRType{
    return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RType is not set"));
  }
  return nil
}

func (p *CsBean)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI16(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.RType = v
}
  return nil
}

func (p *CsBean)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string][]byte, size)
  p.Bsm =  tMap
  for i := 0; i < size; i ++ {
var _key9 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key9 = v
}
var _val10 []byte
    if v, err := iprot.ReadBinary(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val10 = v
}
    p.Bsm[_key9] = _val10
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *CsBean)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string][]int64, size)
  p.Bsm2 =  tMap
  for i := 0; i < size; i ++ {
var _key11 string
    if v, err := iprot.ReadString(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key11 = v
}
    _, size, err := iprot.ReadListBegin(ctx)
    if err != nil {
      return thrift.PrependError("error reading list begin: ", err)
    }
    tSlice := make([]int64, 0, size)
    _val12 :=  tSlice
    for i := 0; i < size; i ++ {
var _elem13 int64
      if v, err := iprot.ReadI64(ctx); err != nil {
      return thrift.PrependError("error reading field 0: ", err)
} else {
      _elem13 = v
}
      _val12 = append(_val12, _elem13)
    }
    if err := iprot.ReadListEnd(ctx); err != nil {
      return thrift.PrependError("error reading list end: ", err)
    }
    p.Bsm2[_key11] = _val12
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *CsBean) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsBean"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *CsBean) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "rType", thrift.I16, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:rType: ", p), err) }
  if err := oprot.WriteI16(ctx, int16(p.RType)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.rType (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:rType: ", p), err) }
  return err
}

func (p *CsBean) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBsm() {
    if err := oprot.WriteFieldBegin(ctx, "bsm", thrift.MAP, 2); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:bsm: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Bsm)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Bsm {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteBinary(ctx, v); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 2:bsm: ", p), err) }
  }
  return err
}

func (p *CsBean) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if p.IsSetBsm2() {
    if err := oprot.WriteFieldBegin(ctx, "bsm2", thrift.MAP, 3); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:bsm2: ", p), err) }
    if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.LIST, len(p.Bsm2)); err != nil {
      return thrift.PrependError("error writing map begin: ", err)
    }
    for k, v := range p.Bsm2 {
      if err := oprot.WriteString(ctx, string(k)); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      if err := oprot.WriteListBegin(ctx, thrift.I64, len(v)); err != nil {
        return thrift.PrependError("error writing list begin: ", err)
      }
      for _, v := range v {
        if err := oprot.WriteI64(ctx, int64(v)); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
      }
      if err := oprot.WriteListEnd(ctx); err != nil {
        return thrift.PrependError("error writing list end: ", err)
      }
    }
    if err := oprot.WriteMapEnd(ctx); err != nil {
      return thrift.PrependError("error writing map end: ", err)
    }
    if err := oprot.WriteFieldEnd(ctx); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 3:bsm2: ", p), err) }
  }
  return err
}

func (p *CsBean) Equals(other *CsBean) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.RType != other.RType { return false }
  if len(p.Bsm) != len(other.Bsm) { return false }
  for k, _tgt := range p.Bsm {
    _src14 := other.Bsm[k]
    if bytes.Compare(_tgt, _src14) != 0 { return false }
  }
  if len(p.Bsm2) != len(other.Bsm2) { return false }
  for k, _tgt := range p.Bsm2 {
    _src15 := other.Bsm2[k]
    if len(_tgt) != len(_src15) { return false }
    for i, _tgt := range _tgt {
      _src16 := _src15[i]
      if _tgt != _src16 { return false }
    }
  }
  return true
}

func (p *CsBean) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("CsBean(%+v)", *p)
}

func (p *CsBean) Validate() error {
  return nil
}
type Itnet interface {
  // Parameters:
  //  - PingBs
  Ping(ctx context.Context, pingBs []byte) (_err error)
  // Parameters:
  //  - PongBs
  Pong(ctx context.Context, pongBs []byte) (_err error)
  // Parameters:
  //  - Bss
  Chap(ctx context.Context, bss []byte) (_err error)
  // Parameters:
  //  - AuthKey
  Auth2(ctx context.Context, authKey []byte) (_err error)
  // Parameters:
  //  - Node
  //  - Ir
  SyncNode(ctx context.Context, node *Node, ir bool) (_err error)
  // Parameters:
  //  - Node
  //  - Ir
  SyncAddr(ctx context.Context, node string, ir bool) (_err error)
  // Parameters:
  //  - SyncList
  SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error)
  // Parameters:
  //  - SendId
  //  - Cu
  CsUser(ctx context.Context, sendId int64, cu *CsUser) (_err error)
  // Parameters:
  //  - SendId
  //  - Cb
  CsBs(ctx context.Context, sendId int64, cb *CsBs) (_err error)
  // Parameters:
  //  - SendId
  //  - Ack
  //  - Cb
  CsReq(ctx context.Context, sendId int64, ack bool, cb *CsBean) (_err error)
  // Parameters:
  //  - SendId
  //  - Vrb
  CsVr(ctx context.Context, sendId int64, vrb *VBean) (_err error)
}

type ItnetClient struct {
  c thrift.TClient
  meta thrift.ResponseMeta
}

func NewItnetClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ItnetClient {
  return &ItnetClient{
    c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
  }
}

func NewItnetClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ItnetClient {
  return &ItnetClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewItnetClient(c thrift.TClient) *ItnetClient {
  return &ItnetClient{
    c: c,
  }
}

func (p *ItnetClient) Client_() thrift.TClient {
  return p.c
}

func (p *ItnetClient) LastResponseMeta_() thrift.ResponseMeta {
  return p.meta
}

func (p *ItnetClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
  p.meta = meta
}

// Parameters:
//  - PingBs
func (p *ItnetClient) Ping(ctx context.Context, pingBs []byte) (_err error) {
  var _args17 ItnetPingArgs
  _args17.PingBs = pingBs
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Ping", &_args17, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - PongBs
func (p *ItnetClient) Pong(ctx context.Context, pongBs []byte) (_err error) {
  var _args18 ItnetPongArgs
  _args18.PongBs = pongBs
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Pong", &_args18, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Bss
func (p *ItnetClient) Chap(ctx context.Context, bss []byte) (_err error) {
  var _args19 ItnetChapArgs
  _args19.Bss = bss
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Chap", &_args19, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - AuthKey
func (p *ItnetClient) Auth2(ctx context.Context, authKey []byte) (_err error) {
  var _args20 ItnetAuth2Args
  _args20.AuthKey = authKey
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "Auth2", &_args20, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Node
//  - Ir
func (p *ItnetClient) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
  var _args21 ItnetSyncNodeArgs
  _args21.Node = node
  _args21.Ir = ir
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncNode", &_args21, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - Node
//  - Ir
func (p *ItnetClient) SyncAddr(ctx context.Context, node string, ir bool) (_err error) {
  var _args22 ItnetSyncAddrArgs
  _args22.Node = node
  _args22.Ir = ir
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncAddr", &_args22, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SyncList
func (p *ItnetClient) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
  var _args23 ItnetSyncTxMergeArgs
  _args23.SyncList = syncList
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "SyncTxMerge", &_args23, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SendId
//  - Cu
func (p *ItnetClient) CsUser(ctx context.Context, sendId int64, cu *CsUser) (_err error) {
  var _args24 ItnetCsUserArgs
  _args24.SendId = sendId
  _args24.Cu = cu
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CsUser", &_args24, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SendId
//  - Cb
func (p *ItnetClient) CsBs(ctx context.Context, sendId int64, cb *CsBs) (_err error) {
  var _args25 ItnetCsBsArgs
  _args25.SendId = sendId
  _args25.Cb = cb
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CsBs", &_args25, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SendId
//  - Ack
//  - Cb
func (p *ItnetClient) CsReq(ctx context.Context, sendId int64, ack bool, cb *CsBean) (_err error) {
  var _args26 ItnetCsReqArgs
  _args26.SendId = sendId
  _args26.Ack = ack
  _args26.Cb = cb
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CsReq", &_args26, nil); err != nil {
    return err
  }
  return nil
}

// Parameters:
//  - SendId
//  - Vrb
func (p *ItnetClient) CsVr(ctx context.Context, sendId int64, vrb *VBean) (_err error) {
  var _args27 ItnetCsVrArgs
  _args27.SendId = sendId
  _args27.Vrb = vrb
  p.SetLastResponseMeta_(thrift.ResponseMeta{})
  if _, err := p.Client_().Call(ctx, "CsVr", &_args27, nil); err != nil {
    return err
  }
  return nil
}

type ItnetProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler Itnet
}

func (p *ItnetProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *ItnetProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *ItnetProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewItnetProcessor(handler Itnet) *ItnetProcessor {

  self28 := &ItnetProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
  self28.processorMap["Ping"] = &itnetProcessorPing{handler:handler}
  self28.processorMap["Pong"] = &itnetProcessorPong{handler:handler}
  self28.processorMap["Chap"] = &itnetProcessorChap{handler:handler}
  self28.processorMap["Auth2"] = &itnetProcessorAuth2{handler:handler}
  self28.processorMap["SyncNode"] = &itnetProcessorSyncNode{handler:handler}
  self28.processorMap["SyncAddr"] = &itnetProcessorSyncAddr{handler:handler}
  self28.processorMap["SyncTxMerge"] = &itnetProcessorSyncTxMerge{handler:handler}
  self28.processorMap["CsUser"] = &itnetProcessorCsUser{handler:handler}
  self28.processorMap["CsBs"] = &itnetProcessorCsBs{handler:handler}
  self28.processorMap["CsReq"] = &itnetProcessorCsReq{handler:handler}
  self28.processorMap["CsVr"] = &itnetProcessorCsVr{handler:handler}
return self28
}

func (p *ItnetProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
  if err2 != nil { return false, thrift.WrapTException(err2) }
  if processor, ok := p.GetProcessorFunction(name); ok {
    return processor.Process(ctx, seqId, iprot, oprot)
  }
  iprot.Skip(ctx, thrift.STRUCT)
  iprot.ReadMessageEnd(ctx)
  x29 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
  oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
  x29.Write(ctx, oprot)
  oprot.WriteMessageEnd(ctx)
  oprot.Flush(ctx)
  return false, x29

}

type itnetProcessorPing struct {
  handler Itnet
}

func (p *itnetProcessorPing) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPingArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Ping(ctx, args.PingBs); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorPong struct {
  handler Itnet
}

func (p *itnetProcessorPong) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetPongArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Pong(ctx, args.PongBs); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorChap struct {
  handler Itnet
}

func (p *itnetProcessorChap) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetChapArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Chap(ctx, args.Bss); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorAuth2 struct {
  handler Itnet
}

func (p *itnetProcessorAuth2) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetAuth2Args{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.Auth2(ctx, args.AuthKey); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncNode struct {
  handler Itnet
}

func (p *itnetProcessorSyncNode) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncNodeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncNode(ctx, args.Node, args.Ir); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncAddr struct {
  handler Itnet
}

func (p *itnetProcessorSyncAddr) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncAddrArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncAddr(ctx, args.Node, args.Ir); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorSyncTxMerge struct {
  handler Itnet
}

func (p *itnetProcessorSyncTxMerge) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetSyncTxMergeArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.SyncTxMerge(ctx, args.SyncList); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCsUser struct {
  handler Itnet
}

func (p *itnetProcessorCsUser) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCsUserArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CsUser(ctx, args.SendId, args.Cu); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCsBs struct {
  handler Itnet
}

func (p *itnetProcessorCsBs) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCsBsArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CsBs(ctx, args.SendId, args.Cb); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCsReq struct {
  handler Itnet
}

func (p *itnetProcessorCsReq) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCsReqArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CsReq(ctx, args.SendId, args.Ack, args.Cb); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}

type itnetProcessorCsVr struct {
  handler Itnet
}

func (p *itnetProcessorCsVr) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := ItnetCsVrArgs{}
  if err2 := args.Read(ctx, iprot); err2 != nil {
    iprot.ReadMessageEnd(ctx)
    return false, thrift.WrapTException(err2)
  }
  iprot.ReadMessageEnd(ctx)

  tickerCancel := func() {}
  _ = tickerCancel

  if err2 := p.handler.CsVr(ctx, args.SendId, args.Vrb); err2 != nil {
    tickerCancel()
    err = thrift.WrapTException(err2)
  }
  tickerCancel()
  return true, err
}


// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - PingBs
type ItnetPingArgs struct {
  PingBs []byte `thrift:"pingBs,1" db:"pingBs" json:"pingBs"`
}

func NewItnetPingArgs() *ItnetPingArgs {
  return &ItnetPingArgs{}
}


func (p *ItnetPingArgs) GetPingBs() []byte {
  return p.PingBs
}
func (p *ItnetPingArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPingArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.PingBs = v
}
  return nil
}

func (p *ItnetPingArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Ping_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPingArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pingBs", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pingBs: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.PingBs); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pingBs (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pingBs: ", p), err) }
  return err
}

func (p *ItnetPingArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPingArgs(%+v)", *p)
}

// Attributes:
//  - PongBs
type ItnetPongArgs struct {
  PongBs []byte `thrift:"pongBs,1" db:"pongBs" json:"pongBs"`
}

func NewItnetPongArgs() *ItnetPongArgs {
  return &ItnetPongArgs{}
}


func (p *ItnetPongArgs) GetPongBs() []byte {
  return p.PongBs
}
func (p *ItnetPongArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetPongArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.PongBs = v
}
  return nil
}

func (p *ItnetPongArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Pong_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetPongArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "pongBs", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:pongBs: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.PongBs); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.pongBs (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:pongBs: ", p), err) }
  return err
}

func (p *ItnetPongArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetPongArgs(%+v)", *p)
}

// Attributes:
//  - Bss
type ItnetChapArgs struct {
  Bss []byte `thrift:"bss,1" db:"bss" json:"bss"`
}

func NewItnetChapArgs() *ItnetChapArgs {
  return &ItnetChapArgs{}
}


func (p *ItnetChapArgs) GetBss() []byte {
  return p.Bss
}
func (p *ItnetChapArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetChapArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Bss = v
}
  return nil
}

func (p *ItnetChapArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Chap_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetChapArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "bss", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:bss: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.Bss); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.bss (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:bss: ", p), err) }
  return err
}

func (p *ItnetChapArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetChapArgs(%+v)", *p)
}

// Attributes:
//  - AuthKey
type ItnetAuth2Args struct {
  AuthKey []byte `thrift:"authKey,1" db:"authKey" json:"authKey"`
}

func NewItnetAuth2Args() *ItnetAuth2Args {
  return &ItnetAuth2Args{}
}


func (p *ItnetAuth2Args) GetAuthKey() []byte {
  return p.AuthKey
}
func (p *ItnetAuth2Args) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetAuth2Args)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.AuthKey = v
}
  return nil
}

func (p *ItnetAuth2Args) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "Auth2_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetAuth2Args) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "authKey", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:authKey: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.AuthKey); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.authKey (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:authKey: ", p), err) }
  return err
}

func (p *ItnetAuth2Args) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetAuth2Args(%+v)", *p)
}

// Attributes:
//  - Node
//  - Ir
type ItnetSyncNodeArgs struct {
  Node *Node `thrift:"node,1" db:"node" json:"node"`
  Ir bool `thrift:"ir,2" db:"ir" json:"ir"`
}

func NewItnetSyncNodeArgs() *ItnetSyncNodeArgs {
  return &ItnetSyncNodeArgs{}
}

var ItnetSyncNodeArgs_Node_DEFAULT *Node
func (p *ItnetSyncNodeArgs) GetNode() *Node {
  if !p.IsSetNode() {
    return ItnetSyncNodeArgs_Node_DEFAULT
  }
return p.Node
}

func (p *ItnetSyncNodeArgs) GetIr() bool {
  return p.Ir
}
func (p *ItnetSyncNodeArgs) IsSetNode() bool {
  return p.Node != nil
}

func (p *ItnetSyncNodeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncNodeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  p.Node = &Node{}
  if err := p.Node.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Node), err)
  }
  return nil
}

func (p *ItnetSyncNodeArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Ir = v
}
  return nil
}

func (p *ItnetSyncNodeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncNode_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncNodeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "node", thrift.STRUCT, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:node: ", p), err) }
  if err := p.Node.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Node), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:node: ", p), err) }
  return err
}

func (p *ItnetSyncNodeArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ir", thrift.BOOL, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ir: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ir)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ir (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ir: ", p), err) }
  return err
}

func (p *ItnetSyncNodeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncNodeArgs(%+v)", *p)
}

// Attributes:
//  - Node
//  - Ir
type ItnetSyncAddrArgs struct {
  Node string `thrift:"node,1" db:"node" json:"node"`
  Ir bool `thrift:"ir,2" db:"ir" json:"ir"`
}

func NewItnetSyncAddrArgs() *ItnetSyncAddrArgs {
  return &ItnetSyncAddrArgs{}
}


func (p *ItnetSyncAddrArgs) GetNode() string {
  return p.Node
}

func (p *ItnetSyncAddrArgs) GetIr() bool {
  return p.Ir
}
func (p *ItnetSyncAddrArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncAddrArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Node = v
}
  return nil
}

func (p *ItnetSyncAddrArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Ir = v
}
  return nil
}

func (p *ItnetSyncAddrArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncAddr_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncAddrArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "node", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:node: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Node)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.node (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:node: ", p), err) }
  return err
}

func (p *ItnetSyncAddrArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ir", thrift.BOOL, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ir: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ir)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ir (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ir: ", p), err) }
  return err
}

func (p *ItnetSyncAddrArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncAddrArgs(%+v)", *p)
}

// Attributes:
//  - SyncList
type ItnetSyncTxMergeArgs struct {
  SyncList map[int64]int8 `thrift:"syncList,1" db:"syncList" json:"syncList"`
}

func NewItnetSyncTxMergeArgs() *ItnetSyncTxMergeArgs {
  return &ItnetSyncTxMergeArgs{}
}


func (p *ItnetSyncTxMergeArgs) GetSyncList() map[int64]int8 {
  return p.SyncList
}
func (p *ItnetSyncTxMergeArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetSyncTxMergeArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin(ctx)
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[int64]int8, size)
  p.SyncList =  tMap
  for i := 0; i < size; i ++ {
var _key30 int64
    if v, err := iprot.ReadI64(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key30 = v
}
var _val31 int8
    if v, err := iprot.ReadByte(ctx); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    temp := int8(v)
    _val31 = temp
}
    p.SyncList[_key30] = _val31
  }
  if err := iprot.ReadMapEnd(ctx); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *ItnetSyncTxMergeArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "SyncTxMerge_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetSyncTxMergeArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "syncList", thrift.MAP, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:syncList: ", p), err) }
  if err := oprot.WriteMapBegin(ctx, thrift.I64, thrift.BYTE, len(p.SyncList)); err != nil {
    return thrift.PrependError("error writing map begin: ", err)
  }
  for k, v := range p.SyncList {
    if err := oprot.WriteI64(ctx, int64(k)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    if err := oprot.WriteByte(ctx, int8(v)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
  }
  if err := oprot.WriteMapEnd(ctx); err != nil {
    return thrift.PrependError("error writing map end: ", err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:syncList: ", p), err) }
  return err
}

func (p *ItnetSyncTxMergeArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetSyncTxMergeArgs(%+v)", *p)
}

// Attributes:
//  - SendId
//  - Cu
type ItnetCsUserArgs struct {
  SendId int64 `thrift:"sendId,1" db:"sendId" json:"sendId"`
  Cu *CsUser `thrift:"cu,2" db:"cu" json:"cu"`
}

func NewItnetCsUserArgs() *ItnetCsUserArgs {
  return &ItnetCsUserArgs{}
}


func (p *ItnetCsUserArgs) GetSendId() int64 {
  return p.SendId
}
var ItnetCsUserArgs_Cu_DEFAULT *CsUser
func (p *ItnetCsUserArgs) GetCu() *CsUser {
  if !p.IsSetCu() {
    return ItnetCsUserArgs_Cu_DEFAULT
  }
return p.Cu
}
func (p *ItnetCsUserArgs) IsSetCu() bool {
  return p.Cu != nil
}

func (p *ItnetCsUserArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCsUserArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SendId = v
}
  return nil
}

func (p *ItnetCsUserArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Cu = &CsUser{}
  if err := p.Cu.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Cu), err)
  }
  return nil
}

func (p *ItnetCsUserArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsUser_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCsUserArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "sendId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:sendId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SendId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.sendId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:sendId: ", p), err) }
  return err
}

func (p *ItnetCsUserArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "cu", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:cu: ", p), err) }
  if err := p.Cu.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Cu), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:cu: ", p), err) }
  return err
}

func (p *ItnetCsUserArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCsUserArgs(%+v)", *p)
}

// Attributes:
//  - SendId
//  - Cb
type ItnetCsBsArgs struct {
  SendId int64 `thrift:"sendId,1" db:"sendId" json:"sendId"`
  Cb *CsBs `thrift:"cb,2" db:"cb" json:"cb"`
}

func NewItnetCsBsArgs() *ItnetCsBsArgs {
  return &ItnetCsBsArgs{}
}


func (p *ItnetCsBsArgs) GetSendId() int64 {
  return p.SendId
}
var ItnetCsBsArgs_Cb_DEFAULT *CsBs
func (p *ItnetCsBsArgs) GetCb() *CsBs {
  if !p.IsSetCb() {
    return ItnetCsBsArgs_Cb_DEFAULT
  }
return p.Cb
}
func (p *ItnetCsBsArgs) IsSetCb() bool {
  return p.Cb != nil
}

func (p *ItnetCsBsArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCsBsArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SendId = v
}
  return nil
}

func (p *ItnetCsBsArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Cb = &CsBs{}
  if err := p.Cb.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Cb), err)
  }
  return nil
}

func (p *ItnetCsBsArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsBs_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCsBsArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "sendId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:sendId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SendId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.sendId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:sendId: ", p), err) }
  return err
}

func (p *ItnetCsBsArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "cb", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:cb: ", p), err) }
  if err := p.Cb.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Cb), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:cb: ", p), err) }
  return err
}

func (p *ItnetCsBsArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCsBsArgs(%+v)", *p)
}

// Attributes:
//  - SendId
//  - Ack
//  - Cb
type ItnetCsReqArgs struct {
  SendId int64 `thrift:"sendId,1" db:"sendId" json:"sendId"`
  Ack bool `thrift:"ack,2" db:"ack" json:"ack"`
  Cb *CsBean `thrift:"cb,3" db:"cb" json:"cb"`
}

func NewItnetCsReqArgs() *ItnetCsReqArgs {
  return &ItnetCsReqArgs{}
}


func (p *ItnetCsReqArgs) GetSendId() int64 {
  return p.SendId
}

func (p *ItnetCsReqArgs) GetAck() bool {
  return p.Ack
}
var ItnetCsReqArgs_Cb_DEFAULT *CsBean
func (p *ItnetCsReqArgs) GetCb() *CsBean {
  if !p.IsSetCb() {
    return ItnetCsReqArgs_Cb_DEFAULT
  }
return p.Cb
}
func (p *ItnetCsReqArgs) IsSetCb() bool {
  return p.Cb != nil
}

func (p *ItnetCsReqArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField3(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCsReqArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SendId = v
}
  return nil
}

func (p *ItnetCsReqArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Ack = v
}
  return nil
}

func (p *ItnetCsReqArgs)  ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
  p.Cb = &CsBean{}
  if err := p.Cb.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Cb), err)
  }
  return nil
}

func (p *ItnetCsReqArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsReq_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
    if err := p.writeField3(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCsReqArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "sendId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:sendId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SendId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.sendId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:sendId: ", p), err) }
  return err
}

func (p *ItnetCsReqArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "ack", thrift.BOOL, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ack: ", p), err) }
  if err := oprot.WriteBool(ctx, bool(p.Ack)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ack (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ack: ", p), err) }
  return err
}

func (p *ItnetCsReqArgs) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "cb", thrift.STRUCT, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:cb: ", p), err) }
  if err := p.Cb.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Cb), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:cb: ", p), err) }
  return err
}

func (p *ItnetCsReqArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCsReqArgs(%+v)", *p)
}

// Attributes:
//  - SendId
//  - Vrb
type ItnetCsVrArgs struct {
  SendId int64 `thrift:"sendId,1" db:"sendId" json:"sendId"`
  Vrb *VBean `thrift:"vrb,2" db:"vrb" json:"vrb"`
}

func NewItnetCsVrArgs() *ItnetCsVrArgs {
  return &ItnetCsVrArgs{}
}


func (p *ItnetCsVrArgs) GetSendId() int64 {
  return p.SendId
}
var ItnetCsVrArgs_Vrb_DEFAULT *VBean
func (p *ItnetCsVrArgs) GetVrb() *VBean {
  if !p.IsSetVrb() {
    return ItnetCsVrArgs_Vrb_DEFAULT
  }
return p.Vrb
}
func (p *ItnetCsVrArgs) IsSetVrb() bool {
  return p.Vrb != nil
}

func (p *ItnetCsVrArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *ItnetCsVrArgs)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SendId = v
}
  return nil
}

func (p *ItnetCsVrArgs)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  p.Vrb = &VBean{}
  if err := p.Vrb.Read(ctx, iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Vrb), err)
  }
  return nil
}

func (p *ItnetCsVrArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "CsVr_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *ItnetCsVrArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "sendId", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:sendId: ", p), err) }
  if err := oprot.WriteI64(ctx, int64(p.SendId)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.sendId (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:sendId: ", p), err) }
  return err
}

func (p *ItnetCsVrArgs) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "vrb", thrift.STRUCT, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:vrb: ", p), err) }
  if err := p.Vrb.Write(ctx, oprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Vrb), err)
  }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:vrb: ", p), err) }
  return err
}

func (p *ItnetCsVrArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("ItnetCsVrArgs(%+v)", *p)
}


