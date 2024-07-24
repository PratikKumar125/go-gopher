package users

import (
	"context"
	"first/models"
)

type IUserTransform struct {
	Name string `json:"name"`
}

type ICreateUserTranform struct {
	Token string `json:"token"`
}

func TransformUser(ctx context.Context, user models.User) (IUserTransform) {
	transformed := IUserTransform{
		Name: user.Name,
	}
	return transformed
}

func TranformCreateUser(ctx context.Context, token string) (ICreateUserTranform) {
	return ICreateUserTranform{
		Token: token,
	}
}