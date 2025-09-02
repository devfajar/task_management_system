package helper

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ========== INT4 (nullable INT) ==========

// ToInt4  mengubah *int32 → pgtype.Int4 (NULL jika nil)
func ToInt4(p *int32) pgtype.Int4 {
	if p == nil {
		return pgtype.Int4{} // Valid=false → NULL
	}
	return pgtype.Int4{Int32: *p, Valid: true}
}

// ToInt4Ptr  mengubah *int32 → *pgtype.Int4 (nil jika nil)
// dipakai kalau sqlc generate param bertipe *pgtype.Int4
func ToInt4Ptr(p *int32) *pgtype.Int4 {
	if p == nil {
		return nil // nil → NULL
	}
	v := pgtype.Int4{Int32: *p, Valid: true}
	return &v
}

// Int4ToPtr  mengubah pgtype.Int4 → *int32 (nil jika NULL)
func Int4ToPtr(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}
	x := v.Int32
	return &x
}

func Int4PtrToPtr(p *pgtype.Int4) *int32 {
	if p == nil || !p.Valid {
		return nil
	}
	x := p.Int32
	return &x
}

// ========== FLOAT8 (nullable double precision) ==========

func ToFloat8(p *float64) pgtype.Float8 {
	if p == nil {
		return pgtype.Float8{} // NULL
	}
	return pgtype.Float8{Float64: *p, Valid: true}
}

func ToFloat8Ptr(p *float64) *pgtype.Float8 {
	if p == nil {
		return nil
	}
	v := pgtype.Float8{Float64: *p, Valid: true}
	return &v
}

func Float8ToPtr(v pgtype.Float8) *float64 {
	if !v.Valid {
		return nil
	}
	x := v.Float64
	return &x
}

func Float64PtrFromPg(f pgtype.Float8) *float64 {
	if !f.Valid {
		return nil
	}
	v := f.Float64
	return &v
}

// ========== TEXT (nullable) ==========

// ToNullableText mengubah string → pgtype.Text (empty string dianggap NULL untuk header/UA/IP opsional)
func toNullableText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{} // NULL
	}
	return pgtype.Text{String: s, Valid: true}
}

func TextToPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}

// ========== UUID (nullable) ==========

func ToUUID(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: v, Valid: true}
}

func UuidToPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	s := uuid.UUID(u.Bytes).String()
	return &s
}

func ToNullableUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{} // Valid=false → NULL
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func ToNullableTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{} // NULL
	}
	return pgtype.Text{String: *s, Valid: true}
}

func OwnerToPtr(o pgtype.UUID) *string {
	if o.Valid {
		s := o.String() // sqlc biasanya expose .UUID
		return &s
	}
	return nil
}

// ========== TIMESTAMPTZ (nullable) ==========

func ToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{} // NULL
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func TimestamptzToPtr(v pgtype.Timestamptz) *time.Time {
	if !v.Valid {
		return nil
	}
	x := v.Time
	return &x
}
