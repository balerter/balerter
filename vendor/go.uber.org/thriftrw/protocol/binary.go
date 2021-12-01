// Copyright (c) 2021 Uber Technologies, Inc.
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

package protocol

import (
	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/protocol/stream"
)

// Binary implements the Thrift Binary Protocol.
// Binary can be cast up to EnvelopeAgnosticProtocol to support DecodeRequest.
//
// Deprecated: Don't use this directly. Use binary.Default.
var Binary Protocol

// BinaryStreamer implements a streaming version of the Thrift Binary Protocol.
var BinaryStreamer stream.Protocol

// EnvelopeAgnosticBinary implements the Thrift Binary Protocol, using
// DecodeRequest for request bodies that may or may not have an envelope.
// This in turn produces a responder with an EncodeResponse method so a handler
// can reply in-kind.
//
// EnvelopeAgnosticBinary makes some practical assumptions about messages
// to be able to distinguish enveloped, versioned enveloped, and not-enveloped
// messages reliably:
//
//  1.  No message will use an envelope version greater than 0x7fff.  This
//  would flip the bit that makes versioned envelopes recognizable.  The only
//  envelope version we recognize today is version 1.
//
//  2.  No message with an unversioned envelope will have a message name
//  (procedure name) longer than 0x00ffffffff (three bytes of length prefix).
//  This would roll into the byte that distinguishes the type of the first
//  field of an un-enveloped struct.  This would require a 16MB procedure name.
//
// The overlapping grammars are:
//
//  1.  Enveloped (strict, with version)
//
//      versionbits:4 methodname~4 seqid:4 struct
//
//      versionbits := 0x80000000 | (version:2 << 16) | messagetype:2
//
//  2.  Enveloped (non-strict, without version)
//
//      methodname~4 messagetype:1 seqid:4 struct
//
//  3.  Unenveloped
//
//      struct := (typeid:1 fieldid:2 <value>)* typeid:1=0x00
//
// A message can be,
//
//  1.  Enveloped (strict)
//
//      The message must begin with 0x80 -- the first byte of the version number.
//
//  2.  Enveloped (non-strict)
//
//      The message begins with the length of the method name. As long as the
//      method name is not 16 MB long, the first byte of this message is 0x00. If
//      the message contains at least 9 other bytes, it is enveloped using the
//      non-strict format.
//
//          4  bytes  method name length (minimum = 0)
//          n  bytes  method name
//          1  byte   message type
//          4  bytes  sequence ID
//          1  byte   empty struct (0x00)
//          --------
//          >10 bytes
//
//  3.  Unenveloped
//
//      If the message begins with 0x00 but does not contain any more bytes, it's
//      a bare, un-enveloped empty struct.
//
//      If the message begins with any other byte, it's an unenveloped message or
//      an invalid request.
var EnvelopeAgnosticBinary EnvelopeAgnosticProtocol

func init() {
	Binary = binary.Default
	BinaryStreamer = binary.Default
	EnvelopeAgnosticBinary = binary.Default
}

// NoEnvelopeResponder responds to a request without an envelope.
//
// Deprecated: Don't use this directly. Use DecodeRequest.
var NoEnvelopeResponder Responder = binary.NoEnvelopeResponder

// EnvelopeV0Responder responds to requests with a non-strict (unversioned) envelope.
//
// Deprecated: Don't use this directly. Use DecodeRequest.
type EnvelopeV0Responder = binary.EnvelopeV0Responder

// EnvelopeV1Responder responds to requests with a strict, version 1 envelope.
//
// Deprecated: Don't use this directly. Use DecodeRequest.
type EnvelopeV1Responder = binary.EnvelopeV1Responder
