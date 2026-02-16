package config

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func (pr *ValueRange) UnmarshalYAML(value *yaml.Node) error {
	var port int
	if err := value.Decode(&port); err == nil {
		*pr = ValueRange{
			Start: port,
			End:   port,
		}
		return nil
	}

	var rangeString string
	if err := value.Decode(&rangeString); err == nil {
		parts := strings.Split(rangeString, "-")
		if len(parts) == 1 {
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}
			*pr = ValueRange{
				Start: v,
				End:   v,
			}
			return nil
		}

		if len(parts) != 2 {
			return fmt.Errorf("invalid value range: %s", rangeString)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		*pr = ValueRange{
			Start: start,
			End:   end,
		}

		return nil
	}

	return fmt.Errorf("value range must be a number or a range like '80-90', got %s", value.Tag)
}

func (pr *ValueRange) MarshalYAML() (interface{}, error) {
	if pr.End == pr.Start {
		return pr.Start, nil
	}
	return fmt.Sprintf("%d-%d", pr.Start, pr.End), nil
}
