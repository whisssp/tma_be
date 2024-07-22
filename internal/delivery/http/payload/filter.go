package payload

import "time"

type TaskFilter struct {
	Keyword       string     `form:"keyword"`
	ID            int64      `form:"id"`
	Title         string     `form:"title"`
	Description   string     `form:"description"`
	Status        int        `form:"status"`
	CreatedAtFrom *time.Time `form:"createdAtFrom"`
	CreatedAtTo   *time.Time `form:"createdAtTo"`
	UpdatedAtFrom *time.Time `form:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `form:"updatedAtTo"`
}