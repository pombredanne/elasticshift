package team

import "time"

// Team ..
type Team struct {
	PUUID     string `gorm:"column:puuid"`
	Domain    string
	Name      string
	CreatedDt time.Time
	CreatedBy string
	UpdatedDt time.Time
	UpdatedBy string
}

// Repository provides access a team.
type Repository interface {
	Save(team *Team) error
}
