//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package histogram

import (
	"math"
	"sort"
)

var (
	negativeInfinity = math.Inf(-1)
	infinity         = math.Inf(1)
	nan              = math.NaN()
)

type Instance []float64

func (h Instance) Min() float64 {
	return h.sorted().min()
}

func (h Instance) Max() float64 {
	return h.sorted().max()
}

func (h Instance) Mean() float64 {
	total := float64(0)
	for _, v := range h {
		total += v
	}
	return total / float64(len(h))
}

func (h Instance) Quantile(phi float64) float64 {
	return h.sorted().quantile(phi)
}

func (h Instance) Quantiles(phis ...float64) []float64 {
	return h.sorted().quantiles(phis...)
}

func (h Instance) sorted() sorted {
	out := make(sorted, 0, len(h))
	out = append(out, h...)

	sort.Float64s(out)
	return out
}

type sorted []float64

func (s sorted) min() float64 {
	if len(s) == 0 {
		return negativeInfinity
	}

	return s[0]
}

func (s sorted) max() float64 {
	if len(s) == 0 {
		return infinity
	}

	return s[len(s)-1]
}

func (s sorted) quantile(phi float64) float64 {
	if len(s) == 0 || math.IsNaN(phi) {
		return nan
	}
	if phi <= 0 {
		return s.min()
	}
	if phi >= 1 {
		return s.max()
	}
	idx := uint(phi*float64(len(s)-1) + 0.5)
	if idx >= uint(len(s)) {
		idx = uint(len(s) - 1)
	}
	return s[idx]
}

func (s sorted) quantiles(phis ...float64) []float64 {
	out := make([]float64, 0, len(phis))
	for _, phi := range phis {
		out = append(out, s.quantile(phi))
	}
	return out
}
