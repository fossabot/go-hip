// Copyright 2019 go-hip authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package param

import (
	"encoding/binary"
)

// Puzzle represents a Puzzle parameter.
//
// Spec: 5.2.4.  PUZZLE
type Puzzle struct {
	*Header
	NoOfK    uint8
	Lifetime uint8
	Opaque   uint16
	Random   []byte
}

// NewPuzzle creates a new Puzzle.
func NewPuzzle(bits, lifetime uint8, opaque uint16, random []byte) *Puzzle {
	p := &Puzzle{
		Header:   &Header{Type: ParamTypePuzzle},
		NoOfK:    bits,
		Lifetime: lifetime,
		Opaque:   opaque,
		Random:   random,
	}

	p.Padding = make([]byte, padlen(4+len(random)))
	p.SetLength()
	return p
}

// DecodePuzzle decodes the given bytes as a Puzzle.
func DecodePuzzle(b []byte) (*Puzzle, error) {
	p := &Puzzle{}
	if err := p.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return p, nil
}

// DecodeFromBytes decodes the given bytes as a Puzzle.
func (p *Puzzle) DecodeFromBytes(b []byte) error {
	l := len(b)
	if l < 9 {
		return ErrTooShortToDecode
	}

	var err error
	p.Header, err = DecodeHeader(b)
	if err != nil {
		return err
	}

	p.NoOfK = p.Header.Contents[0]
	p.Lifetime = p.Header.Contents[1]
	p.Opaque = binary.BigEndian.Uint16(p.Header.Contents[2:4])
	p.Random = p.Header.Contents[4:]

	return nil
}

// Serialize serializes a Puzzle into bytes.
func (p *Puzzle) Serialize() ([]byte, error) {
	b := make([]byte, p.Len())
	if err := p.SerializeTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// SerializeTo serializes a Puzzle into bytes.
func (p *Puzzle) SerializeTo(b []byte) error {
	p.Header.Contents = make([]byte, p.Len()-4)
	p.Header.Contents[0] = p.NoOfK
	p.Header.Contents[1] = p.Lifetime
	binary.BigEndian.PutUint16(p.Header.Contents[2:4], p.Opaque)
	copy(p.Header.Contents[4:], p.Random)

	return p.Header.SerializeTo(b)
}

// Len returns the total length of a Puzzle, including Padding.
func (p *Puzzle) Len() int {
	return 4 + 4 + len(p.Random) + len(p.Padding)
}

// SetLength sets the length of Contents in Length field.
func (p *Puzzle) SetLength() {
	p.Length = uint16(4 + len(p.Random))
}
