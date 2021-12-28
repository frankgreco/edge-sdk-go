package firewall

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func port(from, to int) string {
	if from == to {
		return strconv.Itoa(from)
	}
	return fmt.Sprintf("%d-%d", from, to)
}

func ports(p string) (int, int, error) {
	if !strings.Contains(p, "-") {
		i, err := strconv.Atoi(p)
		if err != nil {
			return 0, 0, errors.New("Malformed port.")
		}
		return i, i, nil
	}

	arr := strings.Split(p, "-")
	from, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, 0, errors.New("Malformed port.")
	}
	to, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, 0, errors.New("Malformed port.")
	}

	return from, to, nil
}
