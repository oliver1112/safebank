package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InstituteDAO struct {
	Db *gorm.DB
}

func NewInstituteDao(db *gorm.DB) *InstituteDAO {
	return &InstituteDAO{
		Db: db,
	}
}

type Institute struct {
	InstituteID   int64  `gorm:"primaryKey" json:"institute_id"`
	InstituteName string `gorm:"unique" json:"edu_institute"`
}

func (Institute) TableName() string {
	return "wsj_institute"
}

func (i *InstituteDAO) GetByID(ctx *gin.Context, ID int64) (Institute, error) {
	var institute Institute
	err := i.Db.WithContext(ctx).Where(&Institute{InstituteID: ID}).First(&institute).Error
	return institute, err
}

func (i *InstituteDAO) GetByName(ctx *gin.Context, instituteName string) (Institute, error) {
	var institute Institute
	err := i.Db.WithContext(ctx).Where(&Institute{InstituteName: instituteName}).First(&institute).Error
	return institute, err
}

func (i *InstituteDAO) CreateOrUpdate(ctx *gin.Context, data Institute) (Institute, error) {
	var institute Institute
	err := i.Db.Where(data).Assign(data).FirstOrCreate(&institute).Error
	return institute, err
}
