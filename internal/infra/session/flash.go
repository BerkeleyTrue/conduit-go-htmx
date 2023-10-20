package session

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type (
	class string

	Flash struct {
		Class   string
		Message string
	}

	flashMap map[class][]string
	flashes  []Flash
)

const (
	Danger  class = "danger"
	Success class = "success"
	Warning class = "warning"
	Info    class = "info"
)

// AddFlash adds a flash message to the session
func AddFlash(fc *fiber.Ctx, _class class, message string) error {
	session, ok := fc.Locals("session").(*session.Session)

	if !ok {
		return fmt.Errorf("session not found")
	}

	_flashMap, ok := session.Get("flash").(flashMap)

	if !ok {
		_flashMap = make(flashMap)
	}

	_flashMap[_class] = append(_flashMap[_class], message)

	session.Set("flash", _flashMap)
	return nil
}

// GetFlashes returns all flashes in the session and clears them
func GetFlashes(fc *fiber.Ctx) ([]Flash, error) {
	session, ok := fc.Locals("session").(*session.Session)

	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	flashMap, ok := session.Get("flash").(flashMap)

	if !ok {
		return nil, nil
	}

	flashes := []Flash{}

	for class, messages := range flashMap {
		for _, message := range messages {
			flashes = append(flashes, Flash{
				Class:   string(class),
				Message: message,
			})
		}
	}

	session.Delete("flash")
	return flashes, nil
}
