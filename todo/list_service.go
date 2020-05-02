package todo

import (
	"net/http"

	"github.com/google/uuid"
)

type ListService struct {
	repo ListRepository
}

func NewListService(repo ListRepository) (service ListService) {
	return ListService{
		repo: repo,
	}
}

func (service ListService) AddItem(w http.ResponseWriter, r *http.Request) {
	param1 := r.URL.Query().Get("item")
	if param1 != "" {
		ID := uuid.New().String()
		err := service.repo.SaveItem(ID, param1)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}

		w.Write([]byte(ID))
	} else {
		w.WriteHeader(500)
	}
}

func (service ListService) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("ID")
	if ID != "" {
		err := service.repo.DeleteItem(ID)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}

		w.Write([]byte("delete success"))
	} else {
		w.WriteHeader(500)
	}
}

func (service ListService) FindItem(w http.ResponseWriter, r *http.Request) {
	param1 := r.URL.Query().Get("ID")
	if param1 != "" {
		res, err := service.repo.FindItem(param1)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}

		w.Write([]byte(res))
	} else {
		w.WriteHeader(500)
	}
}
