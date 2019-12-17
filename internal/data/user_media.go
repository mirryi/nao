package data

import (
	"errors"
	"time"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// UserMedia represents a relationship between a User
// and a Media, containing information about the User's
// opinion on the Media.
type UserMedia struct {
	ID               int
	UserID           int
	MediaID          int
	Status           *WatchStatus
	Priority         *int
	Score            *int
	Recommended      *int
	WatchedInstances []WatchedInstance
	Comments         []Info
	UserMediaListIDs []int
	Version          int
}

// Iden returns the ID.
func (um *UserMedia) Iden() int {
	return um.ID
}

// WatchedInstance contains information about a single
// watch of some Media.
type WatchedInstance struct {
	Episodes  int
	Ongoing   bool
	StartDate *time.Time
	EndDate   *time.Time
	Comments  []Info
}

// WatchStatus is an enum that represents the
// status of a Media's consumption by a User.
type WatchStatus int

const (
	// Completed means that the User has consumed
	// the Media in its entirety at least once.
	Completed WatchStatus = iota

	// Planning means that the User is planning
	// to consume the Media sometime in the future.
	Planning

	// Dropped means that the User has never
	// consumed the Media in its entirety and
	// abandoned it in the middle somewhere.
	Dropped

	// Hold means the User has begun consuming
	// the Media, but has placed it on hold.
	Hold
)

// UnmarshalJSON defines custom JSON deserialization for
// WatchStatus.
func (ws *WatchStatus) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]WatchStatus{"Completed": Completed,
		"Planning": Planning,
		"Dropped":  Dropped,
		"Hold":     Hold,
	}[s]
	if !ok {
		return errors.New("invalid watch status value '" + s + "'")
	}
	*ws = value
	return nil
}

// MarshalJSON defines custom JSON serialization for
// WatchStatus.
func (ws *WatchStatus) MarshalJSON() (v []byte, err error) {
	value, ok := map[WatchStatus]string{Completed: "Completed",
		Planning: "Planning",
		Dropped:  "Dropped",
		Hold:     "Hold",
	}[*ws]
	if !ok {
		return nil, errors.New("Invalid watch status value")
	}
	return json.Marshal(value)
}

// UserMediaBucket is the name of the database bucket
// for UserMedia.
const UserMediaBucket = "UserMedia"

// UserMediaService performs operations on UserMedia.
type UserMediaService struct {
	DB *bolt.DB
	Service
}

// Create persists the given UserMedia.
func (ser *UserMediaService) Create(um *UserMedia) error {
	return Create(um, ser)
}

// Update rumlaces the value of the UserMedia with the
// given ID.
func (ser *UserMediaService) Update(um *UserMedia) error {
	return Update(um, ser)
}

// Delete deletes the UserMedia with the given ID.
func (ser *UserMediaService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of UserMedia.
func (ser *UserMediaService) GetAll() ([]*UserMedia, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of UserMedia that
// pass the filter.
func (ser *UserMediaService) GetFilter(keep func(um *UserMedia) bool) ([]*UserMedia, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		um, err := ser.assertType(m)
		if err != nil {
			return false
		}
		return keep(um)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted UserMedia with the given ID.
func (ser *UserMediaService) GetByID(id int) (*UserMedia, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	um, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}
	return um, nil
}

// GetByUser retrieves the persisted UserMedia with the given
// User ID.
func (ser *UserMediaService) GetByUser(uID int) ([]*UserMedia, error) {
	return ser.GetFilter(func(um *UserMedia) bool {
		return um.UserID == uID
	})
}

// GetByMedia retrieves the persisted UserMedia with the given
// Media ID.
func (ser *UserMediaService) GetByMedia(mID int) ([]*UserMedia, error) {
	return ser.GetFilter(func(um *UserMedia) bool {
		return um.MediaID == mID
	})
}

// Database returns the database reference.
func (ser *UserMediaService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for UserMedia.
func (ser *UserMediaService) Bucket() string {
	return UserMediaBucket
}

// Clean cleans the given UserMedia for storage
func (ser *UserMediaService) Clean(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	if err := infoListClean(e.Comments); err != nil {
		return err
	}
	for _, wi := range e.WatchedInstances {
		if err := infoListClean(wi.Comments); err != nil {
			return err
		}
	}
	return nil
}

// Validate returns an error if the UserMedia is
// not valid for the database.
func (ser *UserMediaService) Validate(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Check if User with ID specified in UserMedia exists
		// Get User bucket, exit if error
		ub, err := Bucket(UserBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.UserID, ub)
		if err != nil {
			return err
		}

		// Check if Media with ID specified in MediaCharacter exists
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return err
		}
		_, err = get(e.MediaID, mb)
		if err != nil {
			return err
		}

		// Check if UserMediaLists with IDs specified in UserMedia exists
		// Get User bucket, exit if error
		umlb, err := Bucket(UserMediaListBucket, tx)
		if err != nil {
			return err
		}
		for _, listID := range e.UserMediaListIDs {
			_, err = get(listID, umlb)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *UserMediaService) Initialize(m Model, id int) error {
	md, err := ser.assertType(m)
	if err != nil {
		return err
	}
	md.ID = id
	md.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing UserMedia in updates.
func (ser *UserMediaService) PersistOldProperties(n Model, o Model) error {
	nm, err := ser.assertType(n)
	if err != nil {
		return err
	}
	om, err := ser.assertType(o)
	if err != nil {
		return err
	}
	nm.Version = om.Version + 1
	return nil
}

// Marshal transforms the given UserMedia into JSON.
func (ser *UserMediaService) Marshal(m Model) ([]byte, error) {
	um, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(um)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into UserMedia.
func (ser *UserMediaService) Unmarshal(buf []byte) (Model, error) {
	var um UserMedia
	err := json.Unmarshal(buf, &um)
	if err != nil {
		return nil, err
	}
	return &um, nil
}

func (ser *UserMediaService) assertType(m Model) (*UserMedia, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	um, ok := m.(*UserMedia)
	if !ok {
		return nil, errors.New("model must be of UserMedia type")
	}
	return um, nil
}

// mapfromModel returns a list of UserMedia type
// asserted from the given list of Model.
func (ser *UserMediaService) mapFromModel(vlist []Model) ([]*UserMedia, error) {
	list := make([]*UserMedia, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.assertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}