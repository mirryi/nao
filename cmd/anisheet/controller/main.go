package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/anisheet/pkg/api"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
	bolt "go.etcd.io/bbolt"
)

// SubController represents a group of endpoints
// with the same prefix in the controller layer
type SubController struct {
	SubRouter *mux.Router
	Service   *data.Service
}

// Controller represents the API controller layer
type Controller struct {
	Router                *mux.Router
	MediaService          *data.MediaService
	EpisodeService        *data.EpisodeService
	CharacterService      *data.CharacterService
	GenreService          *data.GenreService
	ProducerService       *data.ProducerService
	PersonService         *data.PersonService
	MediaRelationService  *data.MediaRelationService
	MediaCharacterService *data.MediaCharacterService
	MediaGenreService     *data.MediaGenreService
	MediaProducerService  *data.MediaProducerService
}

// New returns a new instance of Controller
func New(db *bolt.DB) Controller {
	// Instantiate controller
	router := mux.NewRouter().StrictSlash(true)
	c := Controller{
		Router:                router,
		MediaService:          &data.MediaService{DB: db},
		EpisodeService:        &data.EpisodeService{DB: db},
		CharacterService:      &data.CharacterService{DB: db},
		GenreService:          &data.GenreService{DB: db},
		ProducerService:       &data.ProducerService{DB: db},
		PersonService:         &data.PersonService{DB: db},
		MediaRelationService:  &data.MediaRelationService{DB: db},
		MediaCharacterService: &data.MediaCharacterService{DB: db},
		MediaGenreService:     &data.MediaGenreService{DB: db},
		MediaProducerService:  &data.MediaProducerService{DB: db},
	}

	// Map routing handlers
	c.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := api.StatusGet()
		json.NewEncoder(w).Encode(status)
	})

	MediaSubrouter(&c)
	EpisodeSubrouter(&c)
	CharacterSubrouter(&c)
	GenreSubrouter(&c)
	ProducerSubrouter(&c)
	PersonSubrouter(&c)
	MediaRelationSubrouter(&c)
	MediaCharacterSubrouter(&c)
	MediaCharacterSubrouter(&c)
	MediaGenreSubrouter(&c)
	MediaProducerSubrouter(&c)

	return c
}

func parseID(r *http.Request) (id int, err error) {
	var vars map[string]string = mux.Vars(r)
	idVal, found := vars["id"]
	if !found {
		err = errors.New("no id specified")
		return
	}

	id, err = strconv.Atoi(idVal)
	if err != nil {
		return
	}

	return
}

func encodeResponseBody(body interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(body)
}

func encodeError(err string, debug error, w http.ResponseWriter) {
	errorResponse := api.ErrorResponseNew(err, debug)
	json.NewEncoder(w).Encode(errorResponse)
	return
}

func withDefaultResponseHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Add("Content-Type", "application/json")
	return w
}
