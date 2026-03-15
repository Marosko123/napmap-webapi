/*
 * NAPMap Stations API
 *
 * NAPMap - Alternative fuel stations registry API
 *
 * API version: 1.0.0
 * Contact: maros.bednar@stuba.sk
 */

package napmap

type Station struct {
	// Unique station identifier
	Id string `json:"id" bson:"id"`

	// Station name
	Name string `json:"name" bson:"name"`

	// Type of station - CHARGING or REFUELING
	StationType string `json:"stationType" bson:"stationType"`

	// Supported fuel types
	Fuels []string `json:"fuels" bson:"fuels"`

	// Operator company name
	OperatorName string `json:"operatorName" bson:"operatorName"`

	// Street address
	Address string `json:"address" bson:"address"`

	// City name
	City string `json:"city" bson:"city"`

	// Country code
	Country string `json:"country,omitempty" bson:"country,omitempty"`

	// GPS latitude
	Lat float64 `json:"lat" bson:"lat"`

	// GPS longitude
	Lng float64 `json:"lng" bson:"lng"`

	// Operating hours
	OpeningHours string `json:"openingHours,omitempty" bson:"openingHours,omitempty"`

	// Maximum charging power in kW
	MaxPowerKw *int32 `json:"maxPowerKw,omitempty" bson:"maxPowerKw,omitempty"`

	// Available connector types
	Connectors []string `json:"connectors,omitempty" bson:"connectors,omitempty"`

	// Additional services
	Services []string `json:"services,omitempty" bson:"services,omitempty"`

	// Station status - ACTIVE or INACTIVE
	Status string `json:"status,omitempty" bson:"status,omitempty"`
}
