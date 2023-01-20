// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	lru "github.com/hashicorp/golang-lru/v2"
)

var gridCache *lru.Cache[string, *model.Grid]

func InitializeCaches() {
	size := (configuration.GetConfiguration().GridCacheSize)
	newCache, err := lru.New[string, *model.Grid](size)
	if err != nil {
		configuration.LogError("", "", "Error during cache initialization: %v.", err)
		return
	}
	gridCache = newCache
	configuration.Log("", "", "Cache for grids is initialized with a size of %d.", size)
}

func cacheGrid(grid *model.Grid) {
	if gridCache != nil {
		eviction := gridCache.Add(grid.Uuid, grid)
		configuration.Trace("", "", "Grid %v is added to cache [eviction = %v].", grid, eviction)
	}
}

func getGridFromCache(uuid string) (grid *model.Grid, ok bool) {
	if gridCache != nil {
		grid, ok := gridCache.Get(uuid)
		configuration.Trace("", "", "Get grid %q from cache: ok is %v, grid = %v.", uuid, ok, grid)
		return grid, ok
	}
	return nil, false
}

func removeGridFromCache(uuid string) {
	if gridCache != nil {
		gridCache.Remove(uuid)
		configuration.Trace("", "", "Grid %q id removed from cache.", uuid)
	}
}
