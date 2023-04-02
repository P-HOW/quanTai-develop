package assets

import (
	"database/sql/driver"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// getDenominator returns 10**precision.
func getDenominator(precision int) *big.Int {
	x := big.NewInt(10)
	return new(big.Int).Exp(x, big.NewInt(int64(precision)), nil)
}

func format(i *big.Int, precision int) string {
	r := big.NewRat(1, 1).SetFrac(i, getDenominator(precision))
	return fmt.Sprintf("%v", r.FloatString(precision))
}

// Qai contains a field to represent the smallest units of QAI
type Qai big.Int

// NewQaiFromJuels returns a new struct to represent QAI from it's smallest unit
func NewQaiFromJuels(w int64) *Qai {
	return (*Qai)(big.NewInt(w))
}

// String returns Qai formatted as a string.
func (q *Qai) String() string {
	if q == nil {
		return "0"
	}
	return fmt.Sprintf("%v", (*big.Int)(q))
}

// Qai returns Qai formatted as a string, in QAI units
func (q *Qai) Qai() string {
	if q == nil {
		return "0"
	}
	return format((*big.Int)(q), 18)
}

// SetInt64 delegates to *big.Int.SetInt64
func (q *Qai) SetInt64(w int64) *Qai {
	return (*Qai)((*big.Int)(q).SetInt64(w))
}

// ToInt returns the Qai value as a *big.Int.
func (q *Qai) ToInt() *big.Int {
	return (*big.Int)(q)
}

// ToHash returns a 32 byte representation of this value
func (q *Qai) ToHash() common.Hash {
	return common.BigToHash((*big.Int)(q))
}

// Set delegates to *big.Int.Set
func (q *Qai) Set(x *Qai) *Qai {
	iq := (*big.Int)(q)
	ix := (*big.Int)(x)

	w := iq.Set(ix)
	return (*Qai)(w)
}

// SetString delegates to *big.Int.SetString
func (q *Qai) SetString(s string, base int) (*Qai, bool) {
	w, ok := (*big.Int)(q).SetString(s, base)
	return (*Qai)(w), ok
}

// Cmp defers to big.Int Cmp
func (q *Qai) Cmp(y *Qai) int {
	return (*big.Int)(q).Cmp((*big.Int)(y))
}

// Add defers to big.Int Add
func (q *Qai) Add(x, y *Qai) *Qai {
	iq := (*big.Int)(q)
	ix := (*big.Int)(x)
	iy := (*big.Int)(y)

	return (*Qai)(iq.Add(ix, iy))
}

// Text defers to big.Int Text
func (q *Qai) Text(base int) string {
	return (*big.Int)(q).Text(base)
}

// IsZero returns true when the value is 0 and false otherwise
func (q *Qai) IsZero() bool {
	zero := big.NewInt(0)
	return (*big.Int)(q).Cmp(zero) == 0
}

// Symbol returns QAI
func (*Qai) Symbol() string {
	return "QAI"
}

// Value returns the Qai value for serialization to database.
func (q Qai) Value() (driver.Value, error) {
	b := (big.Int)(q)
	return b.String(), nil
}

// Scan reads the database value and returns an instance.
func (q *Qai) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		decoded, ok := q.SetString(v, 10)
		if !ok {
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Qai", value, value)
		}
		*q = *decoded
	case []uint8:
		// The SQL library returns numeric() types as []uint8 of the string representation
		decoded, ok := q.SetString(string(v), 10)
		if !ok {
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Qai", value, value)
		}
		*q = *decoded
	case int64:
		return fmt.Errorf("unable to convert %v of %T to Qai, is the sql type set to varchar?", value, value)
	default:
		return fmt.Errorf("unable to convert %v of %T to Qai", value, value)
	}

	return nil
}
