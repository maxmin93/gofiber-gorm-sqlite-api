package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"gorm.io/gorm"
)

// Dog type with sql.NullInt16
type Dog struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required,min=3,max=32"`
	Breed     string    `json:"breed" validate:"required"`
	Age       NullInt16 `json:"age" validate:"number"`
	IsGoodBoy bool      `json:"isGoodBoy" gorm:"default:false"`
	// IsGoodBoy bool      `json:"isGoodBoy" gorm:"default:false"`
}

// BeforeDelete prevent delete sample data which ID < 4
// **NOTE: 같은 모듈 안에서만 정의할 수 있음
func (d *Dog) BeforeDelete(tx *gorm.DB) (err error) {
	// log.Printf("BeforeDelete: ID=%d (%t)", d.ID, d.ID < 4)
	if d.ID < 4 {
		log.Printf("cancel: ID=%d", d.ID)
		return errors.New("Sample Data (ID<4) not allowed to delete")
	}
	return
}

// NullInt16 is wrapper for sql.NullInt16
// 참고 https://stackoverflow.com/a/33072822/6811653
type NullInt16 struct {
	sql.NullInt16
}

// ToNullInt16 convert int to sql.NullInt16
func ToNullInt16(v int) NullInt16 {
	return NullInt16{sql.NullInt16{Int16: int16(v), Valid: true}}
}

// MarshalJSON marshal json of NullInt16
func (v NullInt16) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int16)
	} else {
		return json.Marshal(nil)
	}
}

// UnmarshalJSON unmarshal json of NullInt16
func (v *NullInt16) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int16
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int16 = *x
	} else {
		v.Valid = false
	}
	return nil
}
