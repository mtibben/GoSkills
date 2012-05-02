package skills

import (
	"fmt"
)

type Rating struct {
	Mean   float64
	Stddev float64
}

func NewRating(mean, stddev float64) Rating {
	return Rating{mean, stddev}
}

func (r Rating) String() string {
	return fmt.Sprintf("{μ:%.6g σ:%.6g}", r.Mean, r.Stddev)
}
