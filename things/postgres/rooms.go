// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/things"
)

var _ things.RoomRepository = (*roomRepository)(nil)

type roomRepository struct {
	db Database
}

// NewChannelRepository instantiates a PostgreSQL implementation of channel
// repository.
func NewRoomRepository(db Database) things.RoomRepository {
	return &roomRepository{
		db: db,
	}
}


func (tr roomRepository) Save(ctx context.Context, ths ...things.Room) ([]things.Room, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return []things.Room{}, errors.Wrap(things.ErrCreateEntity, err)
	}
	q := `INSERT INTO rooms (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, thing := range ths {
		dbth, err := toDBRoom(thing)
		if err != nil {
			return []things.Room{}, errors.Wrap(things.ErrCreateEntity, err)
		}

		if _, err := tx.NamedExecContext(ctx, q, dbth); err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []things.Room{}, errors.Wrap(things.ErrMalformedEntity, err)
				case errDuplicate:
					return []things.Room{}, errors.Wrap(things.ErrConflict, err)
				}
			}

			return []things.Room{}, errors.Wrap(things.ErrCreateEntity, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return []things.Room{}, errors.Wrap(things.ErrCreateEntity, err)
	}

	return ths, nil
}

func (tr roomRepository) Update(ctx context.Context, t things.Room) error {
	q := `UPDATE things SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBRoom(t)
	if err != nil {
		return errors.Wrap(things.ErrUpdateEntity, err)
	}

	res, errdb := tr.db.NamedExecContext(ctx, q, dbth)
	if errdb != nil {
		pqErr, ok := errdb.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.Wrap(things.ErrMalformedEntity, errdb)
			}
		}

		return errors.Wrap(things.ErrUpdateEntity, errdb)
	}

	cnt, errdb := res.RowsAffected()
	if err != nil {
		return errors.Wrap(things.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return things.ErrNotFound
	}

	return nil
}

func (tr roomRepository) UpdateKey(ctx context.Context, owner, id, key string) error {
	q := `UPDATE things SET key = :key WHERE owner = :owner AND id = :id;`

	dbth := dbRoom{
		ID:    id,
		Owner: owner,
		//Key:   key,
	}

	res, err := tr.db.NamedExecContext(ctx, q, dbth)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid:
				return errors.Wrap(things.ErrMalformedEntity, err)
			case errDuplicate:
				return errors.Wrap(things.ErrConflict, err)
			}
		}

		return errors.Wrap(things.ErrUpdateEntity, err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(things.ErrUpdateEntity, err)
	}

	if cnt == 0 {
		return things.ErrNotFound
	}

	return nil
}


func (tr roomRepository) RetrieveByID(ctx context.Context, owner, id string) (things.Room, error) {
	q := `SELECT name, key, metadata FROM things WHERE id = $1 AND owner = $2;`

	dbth := dbRoom{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return things.Room{}, errors.Wrap(things.ErrNotFound, err)
		}
		return things.Room{}, errors.Wrap(things.ErrSelectEntity, err)
	}

	return toRoom(dbth)
}

func (tr roomRepository) RetrieveByKey(ctx context.Context, key string) (string, error) {
	q := `SELECT id FROM things WHERE key = $1;`

	var id string
	if err := tr.db.QueryRowxContext(ctx, q, key).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Wrap(things.ErrNotFound, err)
		}
		return "", errors.Wrap(things.ErrSelectEntity, err)
	}

	return id, nil
}

func (tr roomRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, tm things.Metadata) (things.RoomsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(tm)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM rooms
		  WHERE owner = :owner %s%s ORDER BY id LIMIT :limit OFFSET :offset;`, mq, nq)

	params := map[string]interface{}{
		"owner":    owner,
		"limit":    limit,
		"offset":   offset,
		"name":     name,
		"metadata": m,
	}

	rows, err := tr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []things.Room
	for rows.Next() {
		dbth := dbRoom{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
		}

		th, err := toRoom(dbth)
		if err != nil {
			return things.RoomsPage{}, errors.Wrap(things.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM things WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}

	page := things.RoomsPage{
		Rooms: items,
		PageMetadata: things.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (tr roomRepository) RetrieveAllWithThings(ctx context.Context, owner string, name string, tm things.Metadata) (things.RoomsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(tm)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`
		SELECT id, name, metadata FROM rooms
	  	WHERE owner = :owner %s%s ORDER BY id;
	`, mq, nq)

	//tq := fmt.Sprintf(`
	//	SELECT id, name, key, metadata FROM things
	//  	WHERE room_id = :id ORDER BY id;
	//`)

	params := map[string]interface{}{
		"owner":    owner,
		"name":     name,
		"metadata": m,
	}

	rows, err := tr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []things.Room
	for rows.Next() {
		dbro := dbRoom{Owner: owner}
		if err := rows.StructScan(&dbro); err != nil {
			return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
		}

		ro, err := toRoom(dbro)
		if err != nil {
			return things.RoomsPage{}, errors.Wrap(things.ErrViewEntity, err)
		}

		trows, err := tr.db.NamedQueryContext(ctx, `
			SELECT id, name, key, metadata FROM things
		 	WHERE room_id = :id ORDER BY id;
		`, map[string]interface{}{
			"id":    ro.ID,
		})
		if err != nil {
			return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
		}
		defer trows.Close()

		var ths []things.Thing
		for trows.Next() {
			dbth := dbThing{Owner: owner}
			if err := trows.StructScan(&dbth); err != nil {
				return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
			}

			th, err := toThing(dbth)

			if err != nil {
				return things.RoomsPage{}, errors.Wrap(things.ErrViewEntity, err)
			}

			ths = append(ths, th)
		}

		ro.Things = ths
		items = append(items, ro)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM things WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return things.RoomsPage{}, errors.Wrap(things.ErrSelectEntity, err)
	}

	page := things.RoomsPage{
		Rooms: items,
		PageMetadata: things.PageMetadata{
			Total:  total,
		},
	}

	return page, nil
}

func (tr roomRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbThing{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM things WHERE id = :id AND owner = :owner;`
	if _, err := tr.db.NamedExecContext(ctx, q, dbth); err != nil {
		return errors.Wrap(things.ErrRemoveEntity, err)
	}
	return nil
}

type dbRoom struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Metadata []byte `db:"metadata"`
}

func toDBRoom(th things.Room) (dbRoom, error) {
	data := []byte("{}")
	if len(th.Metadata) > 0 {
		b, err := json.Marshal(th.Metadata)
		if err != nil {
			return dbRoom{}, errors.Wrap(things.ErrMalformedEntity, err)
		}
		data = b
	}

	return dbRoom{
		ID:       th.ID,
		Owner:    th.Owner,
		Name:     th.Name,
		Metadata: data,
	}, nil
}

func toRoom(dbth dbRoom) (things.Room, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return things.Room{}, errors.Wrap(things.ErrMalformedEntity, err)
	}

	return things.Room{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Metadata: metadata,
	}, nil
}
