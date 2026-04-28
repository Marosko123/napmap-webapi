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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var missing []string
	if station.Name == "" {
		missing = append(missing, "name")
	}
	if station.StationType == "" {
		missing = append(missing, "stationType")
	}
	if len(station.Fuels) == 0 {
		missing = append(missing, "fuels")
	}
	if station.OperatorName == "" {
		missing = append(missing, "operatorName")
	}
	if station.Address == "" {
		missing = append(missing, "address")
	}
	if station.City == "" {
		missing = append(missing, "city")
	}
	if len(missing) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":  "missing required fields",
			"fields": missing,
		})
		return
	}

	if station.Lat < -90 || station.Lat > 90 || station.Lng < -180 || station.Lng > 180 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid GPS coordinates"})
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
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "station already exists"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (o *implStationsAPI) GetStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "stationId is required"})
		return
	}

	station, err := db.FindDocument(c.Request.Context(), stationId)

	switch err {
	case nil:
		c.JSON(http.StatusOK, station)
	case db_service.ErrNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "station not found"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (o *implStationsAPI) UpdateStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "stationId is required"})
		return
	}

	// over že stanica existuje pred update-om
	_, err := db.FindDocument(c.Request.Context(), stationId)
	switch err {
	case nil:
		// ok
	case db_service.ErrNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "station not found"})
		return
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var station Station
	if err := c.ShouldBindJSON(&station); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station.Id = stationId

	err = db.UpdateDocument(c.Request.Context(), stationId, &station)
	switch err {
	case nil:
		c.JSON(http.StatusOK, station)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (o *implStationsAPI) DeleteStation(c *gin.Context) {
	db, ok := getDbService(c)
	if !ok {
		return
	}

	stationId := c.Param("stationId")
	if stationId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "stationId is required"})
		return
	}

	// soft delete - status na INACTIVE
	station, err := db.FindDocument(c.Request.Context(), stationId)
	switch err {
	case nil:
		// ok
	case db_service.ErrNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "station not found"})
		return
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	station.Status = "INACTIVE"

	err = db.UpdateDocument(c.Request.Context(), stationId, station)
	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getDbService(c *gin.Context) (db_service.DbService[Station], bool) {
	value, exists := c.Get("db_service")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service not found in context"})
		return nil, false
	}

	db, ok := value.(db_service.DbService[Station])
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service has unexpected type"})
		return nil, false
	}

	return db, true
}
