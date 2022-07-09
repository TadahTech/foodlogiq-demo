package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/mocks"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("SERVICE: GetEventRequest", func() {
	mockRepo := &mocks.EventsMongoDB{}

	rs := NewServer(mockRepo)

	w := httptest.NewRecorder()

	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	Describe("SERVICE: GetEventRequest", func() {
		Context("When GetEventRequest is successful", func() {
			It("HTTP response is OK, no error, and an event is returned", func() {
				event := &model.Event{
					CreatedAt: time.Now().Format(time.RFC3339),
					IsDeleted: false,
					CreatedBy: ajax.UserID,
					Contents: []*model.Content{{
						Lot:            "1234",
						Gtin:           "abcde",
						BestByDate:     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						ExpirationDate: time.Now().Add(48 * time.Hour).Format(time.RFC3339),
					}},
					ID:   "1235dfg231543j",
					Type: "shipping",
				}
				mockRepo.On("GetEvent", mock.Anything, mock.Anything).Return(event, nil).Once()
				r := httptest.NewRequest(http.MethodGet, "/event?event_id=1235dfg231543j", nil)
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				var mapper *model.Event
				decoder := json.NewDecoder(w.Body)
				decoder.UseNumber()
				err := decoder.Decode(&mapper)

				Expect(err).Should(Succeed())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(mapper.CreatedBy).To(Equal(ajax.UserID))
				Expect(mapper.ID).To(Equal(event.ID))
				Expect(mapper.Type).To(Equal(event.Type))
				Expect(mapper.CreatedAt).To(Equal(event.CreatedAt))
				Expect(mapper.Contents).To(Equal(event.Contents))
			})
		})
		Context("When DeleteEventRequest fails", func() {
			It("Authorization header is empty, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodGet, "/event", nil)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("no Authorization token"))
			})
			It("Authorization header is malformed, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodGet, "/event", nil)
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
				r := httptest.NewRequest(http.MethodGet, "/event", nil)
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
				r := httptest.NewRequest(http.MethodGet, "/event", nil)
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
				mockRepo.On("GetEvent", mock.Anything, mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				r := httptest.NewRequest(http.MethodGet, "/event?event_id=1235dfg231543j", nil)
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
				mockRepo.On("GetEvent", mock.Anything, mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				r := httptest.NewRequest(http.MethodGet, "/event?event_id=1235dfg231543j", nil)
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("cannot find an event id and created_by match"))
			})
			It("EventID exists, is owned by requester, but event is deleted and no event is returned", func() {
				mockRepo.On("GetEvent", mock.Anything, mock.Anything).Return(nil, mongo.ErrNoDocuments).Once()
				r := httptest.NewRequest(http.MethodGet, "/event?event_id=1235dfg231543j", nil)
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("cannot find an event id and created_by match"))
			})
			It("Request is OK, but DB error occurs, 500 code returned", func() {
				mockRepo.On("GetEvent", mock.Anything, mock.Anything).Return(nil, errors.New("random server error")).Once()
				r := httptest.NewRequest(http.MethodGet, "/event?event_id=1235dfg231543j", nil)
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()

				Expect(result.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
