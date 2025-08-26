package helper

import "github.com/jackc/pgx/v5/pgtype"

func ToText(p *string) pgtype.Text {
	if p == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *p, Valid: true}
}
func ToBool(p *bool) pgtype.Bool {
	if p == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *p, Valid: true}
}
