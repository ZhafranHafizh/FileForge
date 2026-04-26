package validation

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePageRange(spec string) ([]int, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("page range is required")
	}

	parts := strings.Split(spec, ",")
	pages := make([]int, 0)
	seen := make(map[int]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("invalid page range %q: empty segment", spec)
		}

		segmentPages, err := parsePageSegment(part)
		if err != nil {
			return nil, fmt.Errorf("invalid page range %q: %w", spec, err)
		}

		for _, page := range segmentPages {
			if _, ok := seen[page]; ok {
				continue
			}
			seen[page] = struct{}{}
			pages = append(pages, page)
		}
	}

	return pages, nil
}

func parsePageSegment(segment string) ([]int, error) {
	if !strings.Contains(segment, "-") {
		page, err := parsePositivePage(segment)
		if err != nil {
			return nil, err
		}
		return []int{page}, nil
	}

	bounds := strings.Split(segment, "-")
	if len(bounds) != 2 || strings.TrimSpace(bounds[0]) == "" || strings.TrimSpace(bounds[1]) == "" {
		return nil, fmt.Errorf("invalid range segment %q", segment)
	}

	start, err := parsePositivePage(bounds[0])
	if err != nil {
		return nil, err
	}

	end, err := parsePositivePage(bounds[1])
	if err != nil {
		return nil, err
	}

	if end < start {
		return nil, fmt.Errorf("range start must be less than or equal to end: %s", segment)
	}

	pages := make([]int, 0, end-start+1)
	for page := start; page <= end; page++ {
		pages = append(pages, page)
	}

	return pages, nil
}

func parsePositivePage(raw string) (int, error) {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("page must be a positive integer: %s", raw)
	}
	if page <= 0 {
		return 0, fmt.Errorf("page must be greater than zero: %d", page)
	}
	return page, nil
}
