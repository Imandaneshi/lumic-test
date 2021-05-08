package main

import (
	"github.com/google/uuid"
	"github.com/latolukasz/orm"
	"time"
)

type imageCategory struct {
	orm.ORM `orm:"redisSearch=default"`
	ID   uint
	Name string `orm:"length=70;required;searchable"`
}

type image struct {
	orm.ORM    `orm:"redisSearch=default"`
	ID         uint
	UserName   string `orm:"length=100;searchable"`
	Caption    string `orm:"length=max;searchable"`
	Link    string `orm:"length=max;searchable"`
	Slug       string `orm:"length=40;required;unique=SlugIndex;searchable"`
	Meta       interface{}
	Created    time.Time `orm:"time=true"`
	FakeDelete bool
	Category   *imageCategory
	Status     string `orm:"enum=status;searchable"`
}

func (image *image) SetEssentialFields() {
	image.Slug = uuid.New().String()
	image.Created = time.Now()
	image.Status = "visible"
}

func (image *image) SetCategory(name string) error {
	category := &imageCategory{}
	sq := &orm.RedisSearchQuery{}
	engine.RedisSearchOne(category, sq.FilterString("Name", name))
	if !category.IsLoaded() {
		category.Name = name
		engine.Flush(category)
	}
	image.Category = category
	return nil
}

func (image *image) Create(userName string, caption string, link string) error {
	image.Link = link
	image.UserName = userName
	image.Caption = caption
	image.SetEssentialFields()
	engine.Flush(image)
	return nil
}