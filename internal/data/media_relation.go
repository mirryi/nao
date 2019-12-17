package data

import (
	"errors"
	"strings"

	json "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

// MediaRelation represents a relationship between single
// instances of Media and Producer.
type MediaRelation struct {
	ID           int
	OwnerID      int
	RelatedID    int
	Relationship string
	Version      int
}

// Iden returns the ID.
func (mr *MediaRelation) Iden() int {
	return mr.ID
}

// MediaRelationBucket is the name of the database bucket for
// MediaRelation.
const MediaRelationBucket = "MediaRelation"

// MediaRelationService performs operations on MediaRelation.
type MediaRelationService struct {
	DB *bolt.DB
	Service
}

// Create persists the given MediaRelation.
func (ser *MediaRelationService) Create(mr *MediaRelation) error {
	return Create(mr, ser)
}

// Update rmrlaces the value of the MediaRelation with the
// given ID.
func (ser *MediaRelationService) Update(mr *MediaRelation) error {
	return Update(mr, ser)
}

// Delete deletes the MediaRelation with the given ID.
func (ser *MediaRelationService) Delete(id int) error {
	return Delete(id, ser)
}

// GetAll retrieves all persisted values of MediaRelation.
func (ser *MediaRelationService) GetAll() ([]*MediaRelation, error) {
	vlist, err := GetAll(ser)
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetFilter retrieves all persisted values of MediaRelation that
// pass the filter.
func (ser *MediaRelationService) GetFilter(keep func(mr *MediaRelation) bool) ([]*MediaRelation, error) {
	vlist, err := GetFilter(ser, func(m Model) bool {
		mr, err := ser.assertType(m)
		if err != nil {
			return false
		}
		return keep(mr)
	})
	if err != nil {
		return nil, err
	}

	return ser.mapFromModel(vlist)
}

// GetByID retrieves the persisted MediaRelation with the given ID.
func (ser *MediaRelationService) GetByID(id int) (*MediaRelation, error) {
	m, err := GetByID(id, ser)
	if err != nil {
		return nil, err
	}

	mr, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}
	return mr, nil
}

// GetByOwner retrieves a list of instances of MediaRelation
// with the given owning Media ID.
func (ser *MediaRelationService) GetByOwner(mID int) ([]*MediaRelation, error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.OwnerID == mID
	})
}

// GetByRelated retrieves a list of instances of MediaRelation
// with the given related Media ID.
func (ser *MediaRelationService) GetByRelated(mID int) ([]*MediaRelation, error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.RelatedID == mID
	})
}

// GetByRelationship retrieves a list of instances of Media Relation
// with the given relationship.
func (ser *MediaRelationService) GetByRelationship(relationship string) ([]*MediaRelation, error) {
	return ser.GetFilter(func(mr *MediaRelation) bool {
		return mr.Relationship == relationship
	})
}

// Database returns the database reference.
func (ser *MediaRelationService) Database() *bolt.DB {
	return ser.DB
}

// Bucket returns the name of the bucket for MediaRelation.
func (ser *MediaRelationService) Bucket() string {
	return MediaRelationBucket
}

// Clean cleans the given MediaRelation for storage
func (ser *MediaRelationService) Clean(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}
	e.Relationship = strings.Trim(e.Relationship, " ")
	return nil
}

// Validate returns an error if the MediaRelation is
// not valid for the database.
func (ser *MediaRelationService) Validate(m Model) error {
	e, err := ser.assertType(m)
	if err != nil {
		return err
	}

	return ser.DB.View(func(tx *bolt.Tx) error {
		// Get Media bucket, exit if error
		mb, err := Bucket(MediaBucket, tx)
		if err != nil {
			return err
		}

		// Check if owning Media with ID specified in new MediaRelation exists
		_, err = get(e.OwnerID, mb)
		if err != nil {
			return err
		}

		// Check if related Media with ID specified in new MediaRelation exists
		_, err = get(e.RelatedID, mb)
		if err != nil {
			return err
		}

		return nil
	})
}

// Initialize sets initial values for some properties.
func (ser *MediaRelationService) Initialize(m Model, id int) error {
	mr, err := ser.assertType(m)
	if err != nil {
		return err
	}
	mr.ID = id
	mr.Version = 0
	return nil
}

// PersistOldProperties maintains certain properties
// of the existing MediaRelation in updates.
func (ser *MediaRelationService) PersistOldProperties(n Model, o Model) error {
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

// Marshal transforms the given MediaRelation into JSON.
func (ser *MediaRelationService) Marshal(m Model) ([]byte, error) {
	mr, err := ser.assertType(m)
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(mr)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Unmarshal parses the given JSON into MediaRelation.
func (ser *MediaRelationService) Unmarshal(buf []byte) (Model, error) {
	var mr MediaRelation
	err := json.Unmarshal(buf, &mr)
	if err != nil {
		return nil, err
	}
	return &mr, nil
}

func (ser *MediaRelationService) assertType(m Model) (*MediaRelation, error) {
	if m == nil {
		return nil, errors.New("model must not be nil")
	}

	mr, ok := m.(*MediaRelation)
	if !ok {
		return nil, errors.New("model must be of MediaRelation type")
	}
	return mr, nil
}

// mapfromModel returns a list of MediaRelation type
// asserted from the given list of Model.
func (ser *MediaRelationService) mapFromModel(vlist []Model) ([]*MediaRelation, error) {
	list := make([]*MediaRelation, len(vlist))
	var err error
	for i, v := range vlist {
		list[i], err = ser.assertType(v)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}