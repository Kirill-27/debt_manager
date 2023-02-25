package data

type Friends struct {
	Friend1Id int `json:"friend1_id" db:"friend1_id" binding:"required"`
	Friend2Id int `json:"friend2_id" db:"friend2_id" binding:"required"`
}
