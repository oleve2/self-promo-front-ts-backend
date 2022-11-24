package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	dbaseServ "self_promo_back/pkg/dbase"
	emailServ "self_promo_back/pkg/email"

	"self_promo_back/pkg/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Server struct {
	mux         chi.Router
	emailSvc    *emailServ.Service
	dbaseSvc    *dbaseServ.Service
	accessToken string
}

// NewServer
func NewServer(
	mux chi.Router,
	emailSvc *emailServ.Service,
	dbaseSvc *dbaseServ.Service,
	accessToken string,
) *Server {
	return &Server{
		mux:         mux,
		emailSvc:    emailSvc,
		dbaseSvc:    dbaseSvc,
		accessToken: accessToken,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// https://stackoverflow.com/questions/40485248/cors-on-golang-server-javascript-fetch-frontend
// https://github.com/go-chi/cors

// Init
func (s *Server) Init() error {
	s.mux.Route("/api/v1", func(r chi.Router) {
		// middlewares
		checkTokenMd := Auth(func(ctx context.Context, token string) (bool, error) {
			// 01-3
			if token == s.accessToken {
				return true, nil
			}
			return false, errors.New("access denied")
		})

		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		r.Use(cors.Handler)

		r.Post("/blogs", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"Ack\": \"OK\"}"))
		})

		//
		r.Get("/echo", s.handleEcho)
		r.With(checkTokenMd).Get("/echo_auth", s.handleEcho)

		r.With(checkTokenMd).Post("/send_email", s.handleSendEmail)
		r.With(checkTokenMd).Get("/all_requests", s.handleAllRequests) //
		r.With(checkTokenMd).Post("/request_insert", s.handleCRInsert)
		r.With(checkTokenMd).Post("/request_update", s.handleCRUpdate)
		r.With(checkTokenMd).Post("/request_delete", s.handleCRDelete)
		r.With(checkTokenMd).Post("/check_name_emailphone", s.handleNameEmailPhoneCheck)
	})

	return nil
}

func WriteAnswer(dataJSON []byte, writer http.ResponseWriter) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write(dataJSON)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// -------------------------------
// echo
func (s *Server) handleEcho(writer http.ResponseWriter, request *http.Request) {
	//msg := fmt.Sprintf("this is echo page")

	msg2 := &models.EchoDTO{Message: "this is echo"}
	dataJSON, err := json.Marshal(msg2)
	if err != nil {
		log.Println(err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write(dataJSON)
}

// -------------------------------
func (s *Server) handleSendEmail(writer http.ResponseWriter, request *http.Request) {
	var ClientRequestDTO *models.ClientRequest
	err := json.NewDecoder(request.Body).Decode(&ClientRequestDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	//fmt.Println("ClientRequestDTO=", ClientRequestDTO)

	// construct email
	subjData := fmt.Sprintf("Письмо SelfPromo %s %s", ClientRequestDTO.Name, ClientRequestDTO.EmailPhone)
	bodyData := fmt.Sprintf("%s", ClientRequestDTO.Request)
	emailDTO := &models.EmailModel{Subject: subjData, Body: bodyData}

	// email send
	err = s.emailSvc.SendMail(emailDTO.Subject, emailDTO.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	// dbase insert
	s.dbaseSvc.InsertCRRow(ClientRequestDTO)

	writer.WriteHeader(http.StatusOK)
	return
}

// crud -------------------------------
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

func (s *Server) handleAllRequests(writer http.ResponseWriter, request *http.Request) {
	//enableCors(&writer)
	data, err := s.dbaseSvc.GetAllCRs()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	//fmt.Println(data)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	WriteAnswer(dataJSON, writer)

	/*
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		_, err = writer.Write(dataJSON)
	*/
}

func (s *Server) handleCRInsert(writer http.ResponseWriter, request *http.Request) {
	var ClientRequestDTO *models.ClientRequest
	err := json.NewDecoder(request.Body).Decode(&ClientRequestDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	fmt.Println("ClientRequestDTO=", ClientRequestDTO)

	s.dbaseSvc.InsertCRRow(ClientRequestDTO)

	writer.WriteHeader(http.StatusOK)
	return
}

func (s *Server) handleCRUpdate(writer http.ResponseWriter, request *http.Request) {
	var ClientRequestDTO *models.ClientRequest
	err := json.NewDecoder(request.Body).Decode(&ClientRequestDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	fmt.Println("ClientRequestDTO=", ClientRequestDTO)

	err = s.dbaseSvc.UpdateCR(ClientRequestDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	return
}

func (s *Server) handleCRDelete(writer http.ResponseWriter, request *http.Request) {
	var deleteDTO *models.CRDeleteDTO
	err := json.NewDecoder(request.Body).Decode(&deleteDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	fmt.Println("ClientRequestDTO=", deleteDTO)

	err = s.dbaseSvc.DeleteCR(deleteDTO.Id)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	return
}

// regexp checks
func (s *Server) handleNameEmailPhoneCheck(writer http.ResponseWriter, request *http.Request) {
	var cnpeDTO *models.CheckNameEmailPhoneDTO
	err := json.NewDecoder(request.Body).Decode(&cnpeDTO)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	checkresp := &models.CheckResponseDTO{
		NameCheck:  s.emailSvc.ValidateName(cnpeDTO.Name),
		EmailCheck: s.emailSvc.ValidateEmail(cnpeDTO.EmailPhone),
		PhoneCheck: s.emailSvc.ValidatePhone(cnpeDTO.EmailPhone),
	}
	dataJSON, err := json.Marshal(checkresp)
	if err != nil {
		log.Println(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write(dataJSON)
	return
}
