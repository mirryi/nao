package data

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Episode represents a single episode or chapter
// for some media
type Episode struct {
	ID       int
	MediaID  int
	Titles   []Info
	Date     *time.Time
	Synopses []Info
	Duration uint
	Filler   bool
	Recap    bool
}

const episodeBucketName = "Episode"

// EpisodeGet retrieves a single instance of Episode with
// the given ID
func EpisodeGet(ID int, db *bolt.DB) (ep Episode, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Episode bucket, exit if error
		b, err := bucket(episodeBucketName, tx)
		if err != nil {
			return err
		}

		// Get Episode by ID, exit if error
		v, err := get(ID, b)
		if err != nil {
			return err
		}

		return json.Unmarshal(v, &ep)
	})

	return
}

// EpisodeGetAll retrieves all persisted Episode values
func EpisodeGetAll(db *bolt.DB) (list []Episode, err error) {
	return EpisodeGetFilter(db, func(ep *Episode) bool { return true })
}

// EpisodeGetFilter retrieves all persisted Episode values
// that pass the filter
func EpisodeGetFilter(db *bolt.DB, filter func(ep *Episode) bool) (list []Episode, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Episode bucket, exit if error
		b, err := bucket(episodeBucketName, tx)

		if err != nil {
			return err
		}

		// Unmarshal and add all Episode to slice,
		// exit if error
		return b.ForEach(func(k, v []byte) error {
			ep := Episode{}
			err = json.Unmarshal(v, &ep)
			if err != nil {
				return err
			}

			if filter(&ep) {
				list = append(list, ep)
			}
			return err
		})
	})

	return
}

// EpisodeGetByMedia retrieves a list of instances of Episode
// with the given Media ID
func EpisodeGetByMedia(mID int, db *bolt.DB) (list []Episode, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		// Get Episode bucket, exit if error
		b, err := bucket(episodeBucketName, tx)
		if err != nil {
			return err
		}

		// Get Episode by Media ID
		return b.ForEach(func(k, v []byte) (err error) {
			ep := Episode{}
			err = json.Unmarshal(v, &ep)
			if err != nil {
				return err
			}

			if ep.MediaID == mID {
				list = append(list, ep)
			}
			return nil
		})
	})

	return
}

// EpisodeCreate persists a new instance of Episode
func EpisodeCreate(ep *Episode, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Episode bucket, exit if error
		b, err := bucket(episodeBucketName, tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and
		// assign to Episode
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		ep.ID = int(id)

		// Save Episode in bucket
		buf, err := json.Marshal(ep)
		if err != nil {
			return err
		}

		return b.Put(itob(ep.ID), buf)
	})
}

// EpisodeUpdate updates the properties of an existing
// persisted Producer instance
func EpisodeUpdate(ep *Episode, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Get Episode bucket, exit if error
		b, err := bucket(episodeBucketName, tx)
		if err != nil {
			return err
		}

		// Check if Episode with ID exists
		_, err = get(ep.ID, b)
		if err != nil {
			return err
		}

		// Replace properties of new with immutable
		// ones of old (none yet)

		// Save Episode
		buf, err := json.Marshal(ep)
		if err != nil {
			return err
		}

		return b.Put(itob(ep.ID), buf)
	})
}
