package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImages (t *testing.T) {
	newCategory := "cat"
	newUsername := "Iman Daneshi"
	newCaption := "test"
	newLink := "https://google.com/somepicture.jpg"
	SetupTestEngine()
	img := image{}
	_ = img.SetCategory(newCategory)
	_ = img.Create(newUsername, newCaption, newLink)
	assert.NotEqual(t, img.Slug, "")
	assert.NotEqual(t, img.ID, uint(0))
	var refreshedImage image
	engine.LoadByID(uint64(img.ID), &refreshedImage, "Category")
	assert.NotEqual(t, refreshedImage, nil)
	assert.NotEqual(t, refreshedImage.ID, uint(0))
	assert.Equal(t, refreshedImage.ID, img.ID)
	assert.Equal(t, refreshedImage.Category.Name, newCategory)
	assert.Equal(t, refreshedImage.UserName, newUsername)
	assert.Equal(t, refreshedImage.Link, newLink)
}