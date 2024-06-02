package di

import (
	"first/controllers"
	"first/crons"
	"first/repositories"
	"first/routes"
	"first/utils"

	"go.uber.org/dig"
)

type Injected struct {
    Utils        utils.DependenciesHolder
    Repositories repositories.DependenciesHolder
    Crons crons.DependenciesHolder
    Router routes.DependenciesHolder
    UserController controllers.DependenciesHolder
    // Tasks        tasks.DependenciesHolder
}

func NewInjected(
    ut utils.DependenciesHolder,
    rp repositories.DependenciesHolder,
    cr crons.DependenciesHolder,
    rt routes.DependenciesHolder,
    userCtrl controllers.DependenciesHolder,
    // tk tasks.DependenciesHolder,
) *Injected {
    return &Injected{
        Utils:        ut,
        Repositories: rp,
        Crons: cr,
        Router: rt,
        UserController: userCtrl,
        // Tasks: tk,
    }
}

var Container *dig.Container

func InitDependencies() error {
    Container = dig.New()
    if err := utils.RegisterDependencies(Container); err != nil {
        return err
    }
    if err := repositories.RegisterRepositories(Container); err != nil {
        return err
    }
    if err := crons.RegisterDependencies(Container); err != nil {
        return err
    }
    if err := routes.RegisterDependencies(Container); err != nil {
        return err
    }
    if err := controllers.RegisterDependencies(Container); err != nil {
        return err
    }
    // if err := tasks.RegisterDependencies(Container); err != nil {
    //     return err
    // }
    if err := Container.Provide(NewInjected); err != nil {
        return err
    }
    return nil
}
