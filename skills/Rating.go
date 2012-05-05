package skills

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills/numerics"
)

type Rating struct {
	Mean   float64
	Stddev float64
}

func NewRating(mean, stddev float64) Rating {
	return Rating{mean, stddev}
}

func (r Rating) Variance() float64 {
	return numerics.Sqr(r.Stddev)
}

func (r Rating) String() string {
	return fmt.Sprintf("{μ:%.6g σ:%.6g}", r.Mean, r.Stddev)
}

func MeanSum(r Rating, a float64) float64 {
	return a + r.Mean
}

func VarianceSum(r Rating, a float64) float64 {
	return a + r.Variance()
}
