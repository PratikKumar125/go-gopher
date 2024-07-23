package users

import (
	"context"
	"first/models"
)

type IUserTransform struct {
	Name string
}

func TransformUser(ctx context.Context, user models.User) (IUserTransform) {
	transformed := IUserTransform{
		Name: user.Name,
	}
	return transformed
}