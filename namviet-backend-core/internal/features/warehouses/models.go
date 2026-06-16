package warehouses

// Warehouse represents the warehouses table
type Warehouse struct {
	ID        int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string   `gorm:"type:text;not null" json:"key"`
	Name      string   `gorm:"type:text;not null" json:"name"`
	Unit      string   `gorm:"type:text;default:'Hộp'" json:"unit"`
	Address   *string  `gorm:"type:text" json:"address,omitempty"`
	Type      string   `gorm:"type:text;default:'retail'" json:"type"`
	Latitude  *float64 `gorm:"type:numeric" json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:numeric" json:"longitude,omitempty"`
	Code      *string  `gorm:"type:text" json:"code,omitempty"`
	Manager   *string  `gorm:"type:text" json:"manager,omitempty"`
	Phone     *string  `gorm:"type:text" json:"phone,omitempty"`
}
