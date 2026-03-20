package napmap

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Marosko123/napmap-webapi/internal/db_service"
)

type implStationsAPI struct {
}

func NewStationsApi() StationsAPI {
	return &implStationsAPI{}
}

func (o *implStationsAPI) GetStations(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	filter := bson.M{}

	if city := c.Query("city"); city != "" {
		filter["city"] = bson.M{"$regex": city, "$options": "i"}
	}

	if fuel := c.Query("fuel"); fuel != "" {
		filter["fuels"] = fuel
	}

	if stationType := c.Query("stationType"); stationType != "" {
		filter["stationType"] = stationType
	}

	if operator := c.Query("operator"); operator != "" {
		filter["operatorName"] = bson.M{"$regex": operator, "$options": "i"}
	}

	status := c.DefaultQuery("status", "ACTIVE")
	if status != "" {
		filter["status"] = status
	}

	if minPower := c.Query("minPowerKw"); minPower != "" {
		if power, err := strconv.Atoi(minPower); err == nil {
			filter["maxPowerKw"] = bson.M{"$gte": power}
		}
	}

	stations, err := db.FindDocuments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to load stations",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stations)
}

func (o *implStationsAPI) CreateStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	var station Station

	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if station.Name == "" || station.StationType == "" || len(station.Fuels) == 0 ||
		station.OperatorName == "" || station.Address == "" || station.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Missing required fields",
		})
		return
	}

	if station.Id == "" || station.Id == "@new" {
		station.Id = uuid.NewString()
	}

	if station.Status == "" {
		station.Status = "ACTIVE"
	}

	if station.Country == "" {
		station.Country = "SK"
	}

	err := db.CreateDocument(c.Request.Context(), station.Id, &station)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, station)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Station already exists",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create station",
			"error":   err.Error(),
		})
	}
}

func (o *implStationsAPI) GetStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Station ID is required",
		})
		return
	}

	station, err := db.FindDocument(c.Request.Context(), stationId)

	switch err {
	case nil:
		c.JSON(http.StatusOK, station)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Station not found",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to load station",
			"error":   err.Error(),
		})
	}
}

func (o *implStationsAPI) UpdateStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Station ID is required",
		})
		return
	}

	// check if station exists
	_, err := db.FindDocument(c.Request.Context(), stationId)
	switch err {
	case nil:
		// ok
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Station not found",
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to load station",
			"error":   err.Error(),
		})
		return
	}

	var station Station
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	station.Id = stationId

	err = db.UpdateDocument(c.Request.Context(), stationId, &station)
	switch err {
	case nil:
		c.JSON(http.StatusOK, station)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update station",
			"error":   err.Error(),
		})
	}
}

func (o *implStationsAPI) DeleteStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Station ID is required",
		})
		return
	}

	// soft delete - set status to INACTIVE
	station, err := db.FindDocument(c.Request.Context(), stationId)
	switch err {
	case nil:
		// ok
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Station not found",
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to load station",
			"error":   err.Error(),
		})
		return
	}

	station.Status = "INACTIVE"

	err = db.UpdateDocument(c.Request.Context(), stationId, station)
	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete station",
			"error":   err.Error(),
		})
	}
}

func getDbService(c *gin.Context) (db_service.DbService[Station], bool) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return nil, false
	}

	db, ok := value.(db_service.DbService[Station])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "db_service context is not of type db_service.DbService",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return nil, false
	}

	return db, true
}
