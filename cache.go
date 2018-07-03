package main

import "sync"

var (
	SavedUserCache	sync.Map
)

func SavedUserStoreCache(U SavedUser) {
	SavedUserCache.Store(U.TgID, U)
}

func SavedUserLoadCache(tgID int64) (SavedUser, bool) {
	if i, ok := SavedUserCache.Load(tgID); ok {
		return i.(SavedUser), ok
	}
	return SavedUser{}, false
}
