package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/TadahTech/foodlogiq-demo/pkg/mocks"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("SERVICE: DeleteEventRequest", func() {
	mockRepo := &mocks.EventsMongoDB{}

	rs := NewServer(mockRepo)

	w := httptest.NewRecorder()

	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	Describe("SERVICE: DeleteEventRequest", func() {
		Context("When DeleteEventRequest is successful", func() {
			It("HTTP response is OK and no error", func() {
				mockRepo.On("DeleteEvent", mock.Anything, mock.Anything).Return(nil).Once()
				request := &model.Event{
					ID: "1235dfg231543j",
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodDelete, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
		Context("When DeleteEventRequest fails", func() {
			It("Authorization header is empty, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodDelete, "/event", nil)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("no Authorization token"))
			})
			It("Authorization header is malformed, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodDelete, "/event", nil)
				r.Header.Add("Authorization", "Bearer:"+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("bearer token malformed"))
			})
			It("Authorization code is unknown, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodDelete, "/event", nil)
				r.Header.Add("Authorization", "Bearer "+"badToken1234")
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("no user found for that token"))
			})
			It("EventID is empty", func() {
				request := &model.Event{}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodDelete, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("event id is empty"))
			})
			It("EventID does not exist", func() {
				mockRepo.On("DeleteEvent", mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments).Once()
				request := &model.Event{
					ID: "totallyReal",
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodDelete, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("cannot find an event id and created_by match"))
			})
			It("EventID exists, but not owned by requester", func() {
				mockRepo.On("DeleteEvent", mock.Anything, mock.Anything).Return(mongo.ErrNoDocuments).Once()
				request := &model.Event{
					ID: "realId",
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodDelete, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("cannot find an event id and created_by match"))
			})
			It("Request is OK, but DB error occurs, 500 code returned", func() {
				mockRepo.On("DeleteEvent", mock.Anything, mock.Anything).Return(errors.New("random server error")).Once()
				request := &model.Event{
					ID: "realId",
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodDelete, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()

				Expect(result.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
