// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/things"
)

const (
	//keyPrefix = "room_key"
	roomPrefix  = "room"
)

var _ things.RoomCache = (*roomCache)(nil)

type roomCache struct {
	client *redis.Client
}

// NewThingCache returns redis thing cache implementation.
func NewRoomCache(client *redis.Client) things.RoomCache {
	return &roomCache{
		client: client,
	}
}

func (tc *roomCache) Save(_ context.Context, thingKey string, thingID string) error {
	//tkey := fmt.Sprintf("%s:%s", keyPrefix, thingKey)
	//if err := tc.client.Set(tkey, thingID, 0).Err(); err != nil {
	//	return errors.Wrap(things.ErrCreateEntity, err)
	//}

	tid := fmt.Sprintf("%s:%s", roomPrefix, thingID)
	if err := tc.client.Set(tid, thingKey, 0).Err(); err != nil {
		return errors.Wrap(things.ErrCreateEntity, err)
	}
	return nil
}

func (tc *roomCache) ID(_ context.Context, thingKey string) (string, error) {
	tkey := fmt.Sprintf("%s:%s", roomPrefix, thingKey)
	thingID, err := tc.client.Get(tkey).Result()
	if err != nil {
		return "", errors.Wrap(things.ErrNotFound, err)
	}

	return thingID, nil
}

func (tc *roomCache) Remove(_ context.Context, thingID string) error {
	tid := fmt.Sprintf("%s:%s", roomPrefix, thingID)
	key, err := tc.client.Get(tid).Result()
	// Redis returns Nil Reply when key does not exist.
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return errors.Wrap(things.ErrRemoveEntity, err)
	}

	tkey := fmt.Sprintf("%s:%s", keyPrefix, key)
	if err := tc.client.Del(tkey, tid).Err(); err != nil {
		return errors.Wrap(things.ErrRemoveEntity, err)
	}
	return nil
}
