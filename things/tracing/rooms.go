// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/mainflux/mainflux/things"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveRoomOp               = "save_room"
	saveRoomsOp              = "save_rooms"
	updateRoomOp             = "update_room"
	updateRoomKeyOp          = "update_room_by_key"
	retrieveRoomByIDOp       = "retrieve_room_by_id"
	retrieveRoomByKeyOp      = "retrieve_room_by_key"
	retrieveAllRoomsOp       = "retrieve_all_rooms"
	retrieveRoomsByChannelOp = "retrieve_rooms_by_chan"
	removeRoomOp             = "remove_room"
	retrieveRoomIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ things.RoomRepository = (*roomRepositoryMiddleware)(nil)
	_ things.RoomCache      = (*roomCacheMiddleware)(nil)
)

type roomRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   things.RoomRepository
}

// RoomRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func RoomRepositoryMiddleware(tracer opentracing.Tracer, repo things.RoomRepository) things.RoomRepository {
	return roomRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm roomRepositoryMiddleware) Save(ctx context.Context, ths ...things.Room) ([]things.Room, error) {
	span := createSpan(ctx, trm.tracer, saveRoomsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm roomRepositoryMiddleware) Update(ctx context.Context, th things.Room) error {
	span := createSpan(ctx, trm.tracer, updateRoomOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm roomRepositoryMiddleware) UpdateKey(ctx context.Context, owner, id, key string) error {
	span := createSpan(ctx, trm.tracer, updateRoomKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.UpdateKey(ctx, owner, id, key)
}

func (trm roomRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (things.Room, error) {
	span := createSpan(ctx, trm.tracer, retrieveRoomByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm roomRepositoryMiddleware) RetrieveByKey(ctx context.Context, key string) (string, error) {
	span := createSpan(ctx, trm.tracer, retrieveRoomByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByKey(ctx, key)
}

func (trm roomRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata things.Metadata) (things.RoomsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllRoomsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm roomRepositoryMiddleware) RetrieveAllWithThings(ctx context.Context, owner string, name string, metadata things.Metadata) (things.RoomsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllRoomsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAllWithThings(ctx, owner, name, metadata)
}

//func (trm roomRepositoryMiddleware) RetrieveByChannel(ctx context.Context, owner, channel string, offset, limit uint64, connected bool) (rooms.Page, error) {
//	span := createSpan(ctx, trm.tracer, retrieveRoomsByChannelOp)
//	defer span.Finish()
//	ctx = opentracing.ContextWithSpan(ctx, span)
//
//	return trm.repo.RetrieveByChannel(ctx, owner, channel, offset, limit, connected)
//}

func (trm roomRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeRoomOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type roomCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  things.RoomCache
}

// RoomCacheMiddleware tracks request and their latency, and adds spans
// to context.
func RoomCacheMiddleware(tracer opentracing.Tracer, cache things.RoomCache) things.RoomCache {
	return roomCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm roomCacheMiddleware) Save(ctx context.Context, roomKey string, roomID string) error {
	span := createSpan(ctx, tcm.tracer, saveRoomOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, roomKey, roomID)
}

func (tcm roomCacheMiddleware) ID(ctx context.Context, roomKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveRoomIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, roomKey)
}

func (tcm roomCacheMiddleware) Remove(ctx context.Context, roomID string) error {
	span := createSpan(ctx, tcm.tracer, removeRoomOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, roomID)
}

//func createSpan(ctx context.Context, tracer opentracing.Tracer, opName string) opentracing.Span {
//	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
//		return tracer.StartSpan(
//			opName,
//			opentracing.ChildOf(parentSpan.Context()),
//		)
//	}
//	return tracer.StartSpan(opName)
//}
