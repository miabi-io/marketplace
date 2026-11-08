/*
 * Copyright 2026 Jonas Kaninda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package registry

import (
	"regexp"
	"sort"
	"strconv"
)

var (
	digitsRe = regexp.MustCompile(`\d+`)
	// A short git SHA ("4511a22") carries no ordering, so such tags are never compared.
	commitRe = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
)

type Candidates struct {
	// Comparable is false for tags with no version order (latest, a commit SHA).
	Comparable bool
	// LatestInMajor treats major.minor as the major for 0.x versions.
	LatestInMajor string
	Latest        string
}

// Compare only matches tags shaped like the current one, so a variant or pre-release is never offered as an upgrade.
func Compare(current string, tags []string) Candidates {
	cur := numbers(current)
	if len(cur) == 0 || (commitRe.MatchString(current) && hasLetter(current)) {
		return Candidates{}
	}
	shape := shapeOf(current)
	majorLen := 1
	if cur[0] == 0 && len(cur) > 1 {
		majorLen = 2
	}
	var same []string
	for _, t := range tags {
		if shapeOf(t) == shape {
			same = append(same, t)
		}
	}
	sort.SliceStable(same, func(i, j int) bool { return less(numbers(same[i]), numbers(same[j])) })
	c := Candidates{Comparable: true, LatestInMajor: current, Latest: current}
	for _, t := range same {
		n := numbers(t)
		if less(numbers(c.Latest), n) {
			c.Latest = t
		}
		if prefixEqual(n, cur, majorLen) && less(numbers(c.LatestInMajor), n) {
			c.LatestInMajor = t
		}
	}
	return c
}

func shapeOf(t string) string { return digitsRe.ReplaceAllString(t, "N") }

func numbers(t string) []int {
	var out []int
	for _, s := range digitsRe.FindAllString(t, -1) {
		n, _ := strconv.Atoi(s)
		out = append(out, n)
	}
	return out
}

func less(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func prefixEqual(a, b []int, n int) bool {
	if len(a) < n || len(b) < n {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func hasLetter(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'f' {
			return true
		}
	}
	return false
}
