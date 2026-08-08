package attendance

import (
	"math"

	"github.com/google/uuid"
)

const earthRadiusMeters = 6371000.0

// haversineDistanceMeters returns the great-circle distance between two
// lat/lng points in meters.
func haversineDistanceMeters(lat1, lng1, lat2, lng2 float64) int {
	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return int(earthRadiusMeters * c)
}

// validateEventLocation checks event.Latitude/Longitude against registered
// attendance_locations (nearest one wins), falling back to the company
// setting's single lat/lng + max_distance_meter when no locations are
// registered. It mutates event's DistanceM/IsInGeofence/ValidatedLocationID
// and returns whether the check passed.
//
// Per §18, this is only meaningful when IsLocationRequired is true - the
// caller decides whether a failed/unresolvable check should fail the event.
func validateEventLocation(event *AttendanceEvent, locations []AttendanceLocation, setting *AttendanceCompanySetting) bool {
	if len(locations) > 0 {
		bestDistance := -1
		var bestLocationID uuid.UUID
		withinAny := false
		for i := range locations {
			loc := locations[i]
			d := haversineDistanceMeters(event.Latitude, event.Longitude, loc.Latitude, loc.Longitude)
			if bestDistance == -1 || d < bestDistance {
				bestDistance = d
				bestLocationID = loc.ID
			}
			if d <= loc.RadiusM {
				withinAny = true
			}
		}
		event.DistanceM = &bestDistance
		event.ValidatedLocationID = &bestLocationID
		event.IsInGeofence = withinAny
		return withinAny
	}

	if setting != nil && setting.Latitude != nil && setting.Longitude != nil && setting.MaxDistanceMeter != nil {
		d := haversineDistanceMeters(event.Latitude, event.Longitude, *setting.Latitude, *setting.Longitude)
		event.DistanceM = &d
		within := d <= *setting.MaxDistanceMeter
		event.IsInGeofence = within
		return within
	}

	// No geofence configured anywhere - can't validate a required check.
	return false
}
