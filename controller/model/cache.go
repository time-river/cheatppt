package model

import (
	"fmt"
	"sync"

	"cheatppt/model/sql"
)

type cacheModel struct {
	DisplayName string
	ModelName   string
	Provider    string
	InputCoins  int // virtual coins
	OutputCoins int // virtual coins
	Activated   bool
}

type cacheModelManager struct {
	models map[string]*cacheModel

	mu *sync.RWMutex
}

var cacheModelMgr *cacheModelManager

func (cm *cacheModelManager) add(key string, val *cacheModel) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.models[key] = val
}

func (cm *cacheModelManager) del(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.models, key)
}

func (cm *cacheModelManager) find(key string) (*cacheModel, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	model, ok := cm.models[key]
	return model, ok
}

// key: `<provider>-<model>`
func BuildCacheKey(provider, model string) string {
	return fmt.Sprintf("%s-%s", provider, model)
}

func CacheAdd(model *sql.Model) {
	key := BuildCacheKey(model.Provider, model.ModelName)
	val := cacheModel{
		DisplayName: model.DisplayName,
		ModelName:   model.ModelName,
		Provider:    model.Provider,
		InputCoins:  model.InputCoins,
		OutputCoins: model.OutputCoins,
		Activated:   model.Activated,
	}
	cacheModelMgr.add(key, &val)
}

func CacheDel(key string) {
	cacheModelMgr.del(key)
}

func CacheFind(key string) (*cacheModel, bool) {
	return cacheModelMgr.find(key)
}
