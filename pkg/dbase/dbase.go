package dbase

import (
	"fmt"
	"log"

	"self_promo_back/pkg/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Init() {
	s.db.AutoMigrate(&models.ClientRequest{})
}

func (s *Service) InsertCRRow(paramRow *models.ClientRequest) {
	s.db.Create(&models.ClientRequest{
		Name:       paramRow.Name,
		EmailPhone: paramRow.EmailPhone,
		Request:    paramRow.Request,
	})
	fmt.Println("client request inserted")
}

func (s *Service) GetAllCRs() ([]*models.ClientRequest, error) {
	var rowsCR []*models.ClientRequest
	//s.db.Raw("select id, name, email_phone, request, created_at, updated_at from client_requests").Scan(rowsCR)
	s.db.Find(&models.ClientRequest{}).Find(&rowsCR)

	return rowsCR, nil
}

func (s *Service) UpdateCR(record *models.ClientRequest) error {
	//result := s.db.Save(&record)
	result := s.db.Where(&models.ClientRequest{ID: record.ID}).Updates(models.ClientRequest{
		Name:       record.Name,
		EmailPhone: record.EmailPhone,
		Request:    record.Request,
	})
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func (s *Service) DeleteCR(id int) error {
	result := s.db.Delete(&models.ClientRequest{}, id)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}
