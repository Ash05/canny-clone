package utils

import (
	"errors"
	"strings"
)

func ValidateBoardName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("Board name cannot be empty")
	}
	if len(name) > 255 {
		return errors.New("Board name cannot exceed 255 characters")
	}
	return nil
}

func ValidateFeedback(title, description string, categoryID int) error {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		return errors.New("Title cannot be empty")
	}
	if len(title) > 255 {
		return errors.New("Title cannot exceed 255 characters")
	}
	if description == "" {
		return errors.New("Description cannot be empty")
	}
	if categoryID <= 0 {
		return errors.New("Invalid category ID")
	}
	return nil
}

// Validate comment content
func ValidateComment(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("Comment cannot be empty")
	}
	if len(content) > 1000 {
		return errors.New("Comment cannot exceed 1000 characters")
	}
	return nil
}
