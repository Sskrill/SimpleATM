package transport

import (
	"encoding/json"
	"github.com/Sskrill/SimpleATM/internal/domain"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Service interface {
	CreateAccount()
	AddBalance(id int, amount float64) error
	WithdrawBalance(id int, amount float64) error
	ShowBalance(id int) (float64, error)
}
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", h.createAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{id}/deposit", h.depositAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{id}/withdraw", h.withdrawAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{id}/balance", h.showAccount).Methods(http.MethodGet)

	return router
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {

	go h.service.CreateAccount()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account has been created"))

	return
}
func (h *Handler) depositAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("nil id"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var sum domain.Sum
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte(err.Error())))
		return
	}
	err = json.Unmarshal(data, &sum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}
	errCh := make(chan error)

	go func() {
		err := h.service.AddBalance(id, sum.Amount)
		errCh <- err
	}()
	err = <-errCh
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("deposited"))
	return
}

func (h *Handler) withdrawAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("nil id"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var sum domain.Sum
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte(err.Error())))
		return
	}
	err = json.Unmarshal(data, &sum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}
	errCh := make(chan error)
	go func() {
		err := h.service.WithdrawBalance(id, sum.Amount)
		errCh <- err
	}()
	err = <-errCh
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("withdraw"))
	return

}
func (h *Handler) showAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("nil id"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	errCh := make(chan error)
	amountCh := make(chan float64)
	go func() {
		amount, err := h.service.ShowBalance(id)
		errCh <- err
		amountCh <- amount
	}()
	err = <-errCh
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}
	var sum domain.Sum

	sum.Amount = <-amountCh

	data, err := json.Marshal(sum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
