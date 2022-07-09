package service

import (
	"bytes"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("SERVICE: CreateEventRequest", func() {
	mockRepo := &mocks.EventsMongoDB{}

	rs := NewServer(mockRepo)

	w := httptest.NewRecorder()
	BeforeEach(func() {
		w = httptest.NewRecorder()
	})

	Describe("SERVICE: CreateEventRequest", func() {
		Context("When CreateEventRequest is successful", func() {
			It("HTTP response is added, no error, and an event ID is returned", func() {
				eventId := primitive.NewObjectID().Hex()
				mockRepo.On("CreateEvent", mock.Anything).Return(eventId, nil).Once()
				request := &model.Event{
					Type: "shipping",
					Contents: []*model.Content{
						{
							Gtin: "1234",
							Lot:  "abcdef",
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				var mapper map[string]interface{}
				decoder := json.NewDecoder(w.Body)
				decoder.UseNumber()
				err := decoder.Decode(&mapper)

				Expect(err).Should(Succeed())
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(mapper["event_id"].(string)).To(Equal(eventId))
			})
		})
		Context("When CreateEvent fails", func() {
			It("Authorization header is empty, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodPost, "/event", nil)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("no Authorization token"))
			})
			It("Authorization header is malformed, 401 code returned", func() {
				r := httptest.NewRequest(http.MethodPost, "/event", nil)
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
				r := httptest.NewRequest(http.MethodPost, "/event", nil)
				r.Header.Add("Authorization", "Bearer "+"badToken1234")
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(string(b)).To(ContainSubstring("no user found for that token"))
			})
			It("Type is invalid, 400 code returned", func() {
				request := &model.Event{
					Type: "somethingInvalid",
					Contents: []*model.Content{
						{
							Gtin: "1234",
							Lot:  "abcdef",
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("invalid type"))
			})
			It("Content is invalid, 400 code returned", func() {
				request := &model.Event{
					Type: "receiving",
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.createEvent(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("content is invalid"))
			})
			It("Content is populated, but GTIN is not given, 400 code returned", func() {
				request := &model.Event{
					Type: "receiving",
					Contents: []*model.Content{
						{
							Lot: "abcdef",
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("gtin is empty"))
			})
			It("Content is populated, but LOT is not given, 400 code returned", func() {
				request := &model.Event{
					Type: "receiving",
					Contents: []*model.Content{
						{
							Gtin: "abcdef",
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("lot is empty"))
			})
			It("Content is populated, and BestBuyDate is given, but not RFC3339, 400 code returned", func() {
				request := &model.Event{
					Type: "receiving",
					Contents: []*model.Content{
						{
							Gtin:       "abcdef",
							Lot:        "1234",
							BestByDate: time.Now().Format(time.RFC850),
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("bestByDate is not RFC3339"))
			})
			It("Content is populated, and ExpirationDate is given, but not RFC3339, 400 code returned", func() {
				request := &model.Event{
					Type: "receiving",
					Contents: []*model.Content{
						{
							Gtin:           "abcdef",
							Lot:            "1234",
							ExpirationDate: time.Now().Format(time.RFC850),
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()
				defer result.Body.Close()

				b, err := io.ReadAll(result.Body)

				Expect(err).To(Succeed())
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(string(b)).To(ContainSubstring("expirationDate is not RFC3339"))
			})
			It("Request is OK, but DB error occurs, 500 code returned", func() {
				mockRepo.On("CreateEvent", mock.Anything).Return("", errors.New("random server error")).Once()
				request := &model.Event{
					Type: "receiving",
					Contents: []*model.Content{
						{
							Gtin: "abcdef",
							Lot:  "1234",
						},
					},
				}
				jsonBlock, _ := json.Marshal(&request)
				r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBlock))
				r.Header.Add("Authorization", "Bearer "+ajax.Token)
				rs.Router.ServeHTTP(w, r)

				result := w.Result()

				Expect(result.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
