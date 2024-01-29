package repository

import (
	"context"
	"payment/models"
)

type Payment interface {
	Create(context.Context, *models.Payment) error
}