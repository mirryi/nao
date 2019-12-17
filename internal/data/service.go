package data

import (
	bolt "go.etcd.io/bbolt"
)

// Model encompasses all data models.
type Model interface {
	Iden() int
}

// Service provides various functions to operate on Models.
// All implementations should use type assertions to guarantee
// prevention of runtime errors.
type Service interface {
	Database() *bolt.DB
	Bucket() string

	Clean(m Model) error
	Validate(m Model) error

	Initialize(m Model, id int) error
	PersistOldProperties(n Model, o Model) error

	Marshal(m Model) ([]byte, error)
	Unmarshal(buf []byte) (Model, error)
}

// Create persists the given Model.
func Create(m Model, ser Service) error {
	err := ser.Clean(m)
	if err != nil {
		return err
	}

	// Verify validity of model
	err = ser.Validate(m)
	if err != nil {
		return err
	}

	return ser.Database().Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		// Get next ID in sequence and assign to
		// model
		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		err = ser.Initialize(m, int(id))
		if err != nil {
			return err
		}

		// Save model in bucket
		buf, err := ser.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(itob(m.Iden()), buf)
	})
}

// Update replaces the value of the model with the given
// ID.
func Update(m Model, ser Service) error {
	err := ser.Clean(m)
	if err != nil {
		return err
	}

	// Verify validity of model
	err = ser.Validate(m)
	if err != nil {
		return err
	}

	return ser.Database().Update(func(tx *bolt.Tx) error {
		// Check if entity with ID exists
		v, err := GetRawByID(m.Iden(), ser.Bucket(), ser.Database())
		if err != nil {
			return err
		}

		// Unmarshall old
		o, err := ser.Unmarshal(v)
		if err != nil {
			return err
		}

		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		// Replace properties of updated with certain frozen
		// ones of old
		err = ser.PersistOldProperties(m, o)
		if err != nil {
			return err
		}

		// Save model
		buf, err := ser.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(itob(m.Iden()), buf)
	})
}

// Delete deletes the model with the given ID.
func Delete(id int, ser Service) error {
	return ser.Database().Update(func(tx *bolt.Tx) error {
		// Get bucket, exit if error
		b, err := Bucket(ser.Bucket(), tx)
		if err != nil {
			return err
		}

		// Store existing model to return
		return b.Delete(itob(id))
	})
}

// GetAll retrieves all persisted instances of a Model type.
func GetAll(ser Service) ([]Model, error) {
	return GetFilter(ser, func(m Model) bool { return true })
}

// GetFilter retrieves all persisted instances of a Model
// type that pass the filter.
func GetFilter(ser Service, keep func(m Model) bool) ([]Model, error) {
	var list []Model
	vlist, err := GetRawAll(ser.Bucket(), ser.Database())
	if err != nil {
		return nil, err
	}

	for _, v := range vlist {
		m, err := ser.Unmarshal(v)
		if err != nil {
			return nil, err
		}

		if keep(m) {
			list = append(list, m)
		}
	}
	return list, nil
}

// GetByID retrieves the persisted Model with the given
// ID.
func GetByID(id int, ser Service) (Model, error) {
	v, err := GetRawByID(id, ser.Bucket(), ser.Database())
	if err != nil {
		return nil, err
	}

	m, err := ser.Unmarshal(v)
	if err != nil {
		return nil, err
	}

	return m, nil
}