package shared

import (
	"fmt"
	"os"
)

func FindFreeFileName(base string) (string, error) {
	modifier := 1
	_, existsCheck := os.Stat(base)
	for ; existsCheck == nil && modifier < 27; modifier = modifier + 1 {
		fmt.Printf("BASE %s", base+string(rune(modifier+96)))
		_, existsCheck = os.Stat(base + string(rune(modifier+96)))
	}
	if existsCheck == nil {
		return "", fmt.Errorf(
			"all slots for the path \"%s\" seem to be already taken: %v", base, existsCheck,
		)
	}
	id := ""
	if modifier > 1 {
		id = string(rune(modifier - 1 + 96))
	}
	return base + id, nil
}
