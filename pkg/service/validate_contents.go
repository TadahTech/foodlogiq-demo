package service

import (
	"errors"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
)

func validateContents(value interface{}) error {
	c, ok := value.([]*model.Content)

	if !ok {
		return errors.New("contents malformed")
	}

	if len(c) == 0 {
		return errors.New("contents empty")
	}

	for _, content := range c {
		if len(content.Gtin) == 0 {
			return errors.New("gtin is empty")
		}
		if len(content.Lot) == 0 {
			return errors.New("lot is empty")
		}
		if len(content.BestByDate) > 0 {
			_, err := time.Parse(time.RFC3339, content.BestByDate)
			if err != nil {
				return errors.New("bestByDate is not RFC3339")
			}
		}
		if len(content.ExpirationDate) > 0 {
			_, err := time.Parse(time.RFC3339, content.ExpirationDate)
			if err != nil {
				return errors.New("expirationDate is not RFC3339")
			}
		}

	}
	return nil
}
