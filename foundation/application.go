package foundation

import (
	"fmt"
	SysLog "log"
	"os"
	"sync"

	"github.com/fwidjaya20/symphonic/config"
	"github.com/fwidjaya20/symphonic/console"
	ContractConfig "github.com/fwidjaya20/symphonic/contracts/config"
	ContractConsole "github.com/fwidjaya20/symphonic/contracts/console"
	ContractDatabase "github.com/fwidjaya20/symphonic/contracts/database"
	ContractEvent "github.com/fwidjaya20/symphonic/contracts/event"
	"github.com/fwidjaya20/symphonic/contracts/foundation"
	ContractLog "github.com/fwidjaya20/symphonic/contracts/log"
	ContractSchedule "github.com/fwidjaya20/symphonic/contracts/schedule"
	"github.com/fwidjaya20/symphonic/database"
	"github.com/fwidjaya20/symphonic/event"
	"github.com/fwidjaya20/symphonic/log"
	"github.com/fwidjaya20/symphonic/schedule"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

var (
	App foundation.Application
)

type _Application struct {
	bindings  sync.Map
	instances sync.Map
}

func init() {
	app := &_Application{}

	baseServiceProvider := app.getBaseServiceProvider()

	app.registerServiceProviders(baseServiceProvider)
	app.bootServiceProviders(baseServiceProvider)

	App = app
}

func NewApplication() foundation.Application {
	return App
}

func (app *_Application) Boot() {
	configuredServiceProviders := app.GetConfig().Get("app.providers").([]foundation.ServiceProvider)
	configuredTz := app.GetConfig().GetString("app.timezone", carbon.UTC)

	app.registerServiceProviders(configuredServiceProviders)
	app.bootServiceProviders(configuredServiceProviders)

	app.GetConsole().Run(os.Args, true)
	carbon.SetTimezone(configuredTz)
}

func (app *_Application) Get(key any) (any, error) {
	binding, ok := app.bindings.Load(key)
	if !ok {
		return nil, fmt.Errorf("binding was not found: %+v", key)
	}

	if instance, ok := app.instances.Load(key); ok {
		return instance, nil
	}

	bindingImpl, err := binding.(func(app foundation.Application) (any, error))(app)
	if nil != err {
		return nil, err
	}

	app.instances.Store(key, bindingImpl)

	return bindingImpl, nil
}

func (app *_Application) GetConfig() ContractConfig.Config {
	instance, err := app.Get(config.Binding)
	if nil != err {
		SysLog.Fatalln(err.Error())
		return nil
	}
	return instance.(ContractConfig.Config)
}

func (app *_Application) GetConsole() ContractConsole.Console {
	instance, err := app.Get(console.Binding)
	if nil != err {
		SysLog.Fatalln(err.Error())
		return nil
	}
	return instance.(ContractConsole.Console)
}

func (app *_Application) GetDatabase() *gorm.DB {
	instance, err := app.Get(database.Binding)
	if err != nil {
		SysLog.Fatalln(err.Error())
		return nil
	}

	dbInstance, ok := instance.(ContractDatabase.Database)
	if !ok {
		SysLog.Fatalln("Instance is not of type ContractDatabase.Database")
		return nil
	}

	return dbInstance.GetSession()
}

func (app *_Application) GetEvent() ContractEvent.Event {
	instance, err := app.Get(event.Binding)
	if nil != err {
		SysLog.Fatalln(err.Error())
		return nil
	}
	return instance.(ContractEvent.Event)
}

func (app *_Application) GetLogger() ContractLog.Logger {
	instance, err := app.Get(log.Binding)
	if nil != err {
		SysLog.Fatalln(err.Error())
		return nil
	}
	return instance.(ContractLog.Logger)
}

func (app *_Application) GetSchedule() ContractSchedule.Schedule {
	instance, err := app.Get(schedule.Binding)
	if nil != err {
		SysLog.Fatalln(err.Error())
		return nil
	}
	return instance.(ContractSchedule.Schedule)
}

func (app *_Application) Instance(key any, callback func(app foundation.Application) (any, error)) {
	app.bindings.Store(key, callback)
}

func (app *_Application) bootServiceProviders(providers []foundation.ServiceProvider) {
	for _, it := range providers {
		it.Boot(app)
	}
}

func (app *_Application) getBaseServiceProvider() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&config.ServiceProvider{},
		&console.ServiceProvider{},
	}
}

func (app *_Application) registerServiceProviders(providers []foundation.ServiceProvider) {
	for _, it := range providers {
		it.Register(app)
	}
}
