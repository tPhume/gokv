package kv

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tPhume/gokv/btree"
	"github.com/tPhume/gokv/store"
	"net/http"
	"strings"
)

var (
	errorWhiteSpaces = "bad format, key cannot contain white spaces"
	errorValueEmpty  = "bad format, value cannot be empty"
	errorBadJSON     = "bad format, json"
	errorInternal    = "an error occurred"
	errorKeyNotFound = "key not found"
)

// Returns gin's Engine that has KeyValue store handlers
// Will return standalone Rest server
func DefaultRestServer() *gin.Engine {
	kvHandlers := NewKeyValueHandlers(btree.NewBtree(3))
	router := gin.Default()
	setHandlers(kvHandlers, router)

	return router
}

// Set default handlers given a gin Engine
func DefaultRestWithEngine(router *gin.Engine) {
	kvHandlers := NewKeyValueHandlers(btree.NewBtree(3))
	setHandlers(kvHandlers, router)
}

// Create new Rest server with store as parameter
func RestWithStore(store store.Store) *gin.Engine {
	kvHandlers := NewKeyValueHandlers(store)
	router := gin.Default()
	setHandlers(kvHandlers, router)

	return router
}

// Utility function to set insert,update,search and delete routes
func setHandlers(kvHandlers *KeyValueHandlers, r *gin.Engine) {
	storeGroup := r.Group("/store")

	storeGroupV1 := storeGroup.Group("/v1")
	storeGroupV1.POST("/:key", kvHandlers.insert)
	storeGroupV1.PATCH("/:key", kvHandlers.update)
	storeGroupV1.GET("/:key", kvHandlers.search)
	storeGroupV1.DELETE("/:key", kvHandlers.remove)
}

// Handles request to the store
type KeyValueHandlers struct {
	store store.Store
}

func NewKeyValueHandlers(store store.Store) *KeyValueHandlers {
	return &KeyValueHandlers{store: store}
}

func (kv *KeyValueHandlers) insert(c *gin.Context) {
	key := c.Param("key")
	if strings.Contains(key, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorWhiteSpaces})
		return
	}

	body := c.Request.Body
	if body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorValueEmpty})
		return
	}

	var value store.Value
	err := json.NewDecoder(body).Decode(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorBadJSON})
		return
	}

	if err := kv.store.Insert(key, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorInternal})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("%v inserted", key)})
}

func (kv *KeyValueHandlers) update(c *gin.Context) {
	key := c.Param("key")
	if strings.Contains(key, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorWhiteSpaces})
		return
	}

	body := c.Request.Body
	if body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorValueEmpty})
		return
	}

	var value store.Value
	err := json.NewDecoder(body).Decode(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorBadJSON})
		return
	}

	if err := kv.store.Update(key, value); err == btree.KeyDoesNotExist {
		c.JSON(http.StatusNotFound, gin.H{"message": errorKeyNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%v updated", key)})
}

func (kv *KeyValueHandlers) search(c *gin.Context) {
	key := c.Param("key")
	if strings.Contains(key, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorWhiteSpaces})
		return
	}

	value := kv.store.Search(key)
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": errorKeyNotFound})
		return
	}

	c.JSON(http.StatusOK, value)
}

func (kv *KeyValueHandlers) remove(c *gin.Context) {
	key := c.Param("key")
	if strings.Contains(key, " ") {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorWhiteSpaces})
		return
	}

	if err := kv.store.Remove(key); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": errorKeyNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key/value deleted"})
}
