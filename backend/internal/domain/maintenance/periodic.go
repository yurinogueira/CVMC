package maintenance

import (
	"math"
	"strings"
	"time"

	domaincar "cvmc/internal/domain/car"
)

type UpcomingMaintenanceStatus string

const (
	StatusOK       UpcomingMaintenanceStatus = "ok"
	StatusUpcoming UpcomingMaintenanceStatus = "upcoming"
	StatusOverdue  UpcomingMaintenanceStatus = "overdue"

	TypeEngineOil   = "Óleo de Motor"
	TypeOilFilter   = "Filtro do Óleo de Motor"
	TypeAirFilter   = "Filtro de Ar do Motor"
	TypeCabinFilter = "Filtro do Ar-Condicionado"
	TypeFuelFilter  = "Filtro de Combustível"

	UpcomingKmThreshold   = 1000
	UpcomingDaysThreshold = 30

	AlertKmThreshold   = 500
	AlertDaysThreshold = 15
)

type PeriodicRule struct {
	Type           string `json:"type"`
	Label          string `json:"label"`
	IntervalKm     int    `json:"intervalKm"`
	IntervalMonths int    `json:"intervalMonths"`
}

type UpcomingMaintenanceItem struct {
	Type            string                    `json:"type"`
	Label           string                    `json:"label"`
	Status          UpcomingMaintenanceStatus `json:"status"`
	IntervalKm      int                       `json:"intervalKm"`
	IntervalMonths  int                       `json:"intervalMonths"`
	HasHistory      bool                      `json:"hasHistory"`
	LastPerformedAt *time.Time                `json:"lastPerformedAt,omitempty"`
	LastMileage     int                       `json:"lastMileage"`
	NextDueMileage  int                       `json:"nextDueMileage"`
	NextDueDate     time.Time                 `json:"nextDueDate"`
	RemainingKm     int                       `json:"remainingKm"`
	RemainingDays   int                       `json:"remainingDays"`
	ProgressPercent int                       `json:"progressPercent"`
	AlertEligible   bool                      `json:"alertEligible"`
}

type CarUpcomingMaintenances struct {
	CarID         string                    `json:"carId"`
	CarName       string                    `json:"carName"`
	Manufacturer  string                    `json:"manufacturer"`
	Model         string                    `json:"model"`
	VehicleType   string                    `json:"vehicleType"`
	CurrentKm     int                       `json:"currentKm"`
	Status        UpcomingMaintenanceStatus `json:"status"`
	OKCount       int                       `json:"okCount"`
	UpcomingCount int                       `json:"upcomingCount"`
	OverdueCount  int                       `json:"overdueCount"`
	Items         []UpcomingMaintenanceItem `json:"items"`
}

type FleetUpcomingSummary struct {
	TotalVehicles int                       `json:"totalVehicles"`
	Status        UpcomingMaintenanceStatus `json:"status"`
	OKCount       int                       `json:"okCount"`
	UpcomingCount int                       `json:"upcomingCount"`
	OverdueCount  int                       `json:"overdueCount"`
	Vehicles      []CarUpcomingMaintenances `json:"vehicles"`
}

// PeriodicRulesForVehicleType returns the preventive maintenance schedule for the given vehicle category.
func PeriodicRulesForVehicleType(vehicleType string) []PeriodicRule {
	vt := strings.ToLower(strings.TrimSpace(vehicleType))
	if vt == "motorcycles" || vt == "motorcycle" || vt == "moto" || vt == "motos" {
		return []PeriodicRule{
			{
				Type:           TypeEngineOil,
				Label:          "Troca de óleo",
				IntervalKm:     2500,
				IntervalMonths: 12,
			},
			{
				Type:           TypeOilFilter,
				Label:          "Troca de filtro de óleo",
				IntervalKm:     2500,
				IntervalMonths: 12,
			},
			{
				Type:           TypeAirFilter,
				Label:          "Troca de filtro de ar",
				IntervalKm:     5000,
				IntervalMonths: 12,
			},
		}
	}

	// Default: cars and trucks ("cars", "trucks")
	return []PeriodicRule{
		{
			Type:           TypeEngineOil,
			Label:          "Troca de óleo",
			IntervalKm:     10000,
			IntervalMonths: 12,
		},
		{
			Type:           TypeOilFilter,
			Label:          "Troca de filtro de óleo",
			IntervalKm:     10000,
			IntervalMonths: 12,
		},
		{
			Type:           TypeAirFilter,
			Label:          "Troca de filtro de ar",
			IntervalKm:     10000,
			IntervalMonths: 12,
		},
		{
			Type:           TypeCabinFilter,
			Label:          "Troca do filtro de ar-condicionado",
			IntervalKm:     20000,
			IntervalMonths: 12,
		},
		{
			Type:           TypeFuelFilter,
			Label:          "Troca do filtro de combustível",
			IntervalKm:     20000,
			IntervalMonths: 12,
		},
	}
}

// AllPeriodicTypes returns all monitored periodic maintenance type strings.
func AllPeriodicTypes() []string {
	return []string{
		TypeEngineOil,
		TypeOilFilter,
		TypeAirFilter,
		TypeCabinFilter,
		TypeFuelFilter,
	}
}

// MatchesPeriodicType checks whether a Maintenance record includes the given canonical periodic type.
func MatchesPeriodicType(m Maintenance, periodicType string) bool {
	target := strings.TrimSpace(periodicType)
	for _, t := range m.Types {
		if strings.EqualFold(strings.TrimSpace(t), target) {
			return true
		}
	}

	// Fallback for records without explicit Types slice: match against Title
	if len(m.Types) == 0 && strings.TrimSpace(m.Title) != "" {
		titleLower := strings.ToLower(m.Title)
		switch target {
		case TypeEngineOil:
			if strings.Contains(titleLower, "óleo") || strings.Contains(titleLower, "oleo") {
				// Match "troca de óleo", "óleo de motor", "óleo e filtro"
				if !strings.Contains(titleLower, "filtro") || strings.Contains(titleLower, "óleo e filtro") || strings.Contains(titleLower, "oleo e filtro") || strings.Contains(titleLower, "troca de óleo") || strings.Contains(titleLower, "troca de oleo") {
					return true
				}
			}
		case TypeOilFilter:
			if strings.Contains(titleLower, "filtro de óleo") || strings.Contains(titleLower, "filtro de oleo") || strings.Contains(titleLower, "filtro do óleo") || strings.Contains(titleLower, "filtro do oleo") || strings.Contains(titleLower, "óleo e filtro") || strings.Contains(titleLower, "oleo e filtro") {
				return true
			}
		case TypeAirFilter:
			if (strings.Contains(titleLower, "filtro de ar") || strings.Contains(titleLower, "filtro do ar")) && !strings.Contains(titleLower, "condicionado") {
				return true
			}
		case TypeCabinFilter:
			if strings.Contains(titleLower, "ar-condicionado") || strings.Contains(titleLower, "ar condicionado") || strings.Contains(titleLower, "cabine") {
				return true
			}
		case TypeFuelFilter:
			if strings.Contains(titleLower, "filtro de combustível") || strings.Contains(titleLower, "filtro de combustivel") {
				return true
			}
		}
	}

	return false
}

// FindLatestByTypesMap extracts the most recent maintenance record for each requested periodic type.
func FindLatestByTypesMap(maintenances []Maintenance, types []string) map[string]Maintenance {
	result := make(map[string]Maintenance, len(types))
	for _, t := range types {
		var latest *Maintenance
		for i := range maintenances {
			m := &maintenances[i]
			if !MatchesPeriodicType(*m, t) {
				continue
			}
			if latest == nil || m.Date.After(latest.Date) || (m.Date.Equal(latest.Date) && m.Mileage > latest.Mileage) {
				latest = m
			}
		}
		if latest != nil {
			result[t] = *latest
		}
	}
	return result
}

// CalculateCarUpcoming evaluates all periodic maintenance items for a vehicle based on its history or registration baseline.
func CalculateCarUpcoming(car domaincar.Car, maintenances []Maintenance, now time.Time) CarUpcomingMaintenances {
	if now.IsZero() {
		now = time.Now().UTC()
	} else {
		now = now.UTC()
	}

	rules := PeriodicRulesForVehicleType(car.VehicleType)
	ruleTypes := make([]string, 0, len(rules))
	for _, r := range rules {
		ruleTypes = append(ruleTypes, r.Type)
	}

	latestByType := FindLatestByTypesMap(maintenances, ruleTypes)

	baseDate := car.CreatedAt.UTC()
	if baseDate.IsZero() {
		baseDate = now
	}

	baseMileage := resolveBaselineMileage(car, maintenances)

	items := make([]UpcomingMaintenanceItem, 0, len(rules))
	okCount := 0
	upcomingCount := 0
	overdueCount := 0

	for _, rule := range rules {
		intervalKm := rule.IntervalKm
		intervalMonths := rule.IntervalMonths
		if car.MaintenanceIntervals != nil {
			if custom, exists := car.MaintenanceIntervals[rule.Type]; exists {
				if custom.IntervalKm > 0 {
					intervalKm = custom.IntervalKm
				}
				if custom.IntervalMonths > 0 {
					intervalMonths = custom.IntervalMonths
				}
			}
		}

		var (
			lastMileage     int
			lastDate        time.Time
			lastPerformedAt *time.Time
			hasHistory      bool
		)

		if latest, exists := latestByType[rule.Type]; exists {
			hasHistory = true
			lastMileage = latest.Mileage
			lastDate = latest.Date.UTC()
			tCopy := lastDate
			lastPerformedAt = &tCopy
		} else {
			hasHistory = false
			lastMileage = baseMileage
			lastDate = baseDate
		}

		nextDueMileage := lastMileage + intervalKm
		nextDueDate := lastDate.AddDate(0, intervalMonths, 0)

		remainingKm := nextDueMileage - car.LastMileage
		diffHours := nextDueDate.Sub(now).Hours()
		remainingDays := int(math.Ceil(diffHours / 24.0))

		// Calculate wear progress percentage (0 - 100)
		kmDriven := car.LastMileage - lastMileage
		if kmDriven < 0 {
			kmDriven = 0
		}
		kmPct := 0
		if intervalKm > 0 {
			kmPct = (kmDriven * 100) / intervalKm
		}

		totalPeriodDays := int(nextDueDate.Sub(lastDate).Hours() / 24.0)
		if totalPeriodDays <= 0 {
			totalPeriodDays = 365
		}
		daysElapsed := int(now.Sub(lastDate).Hours() / 24.0)
		if daysElapsed < 0 {
			daysElapsed = 0
		}
		daysPct := (daysElapsed * 100) / totalPeriodDays

		progressPercent := kmPct
		if daysPct > progressPercent {
			progressPercent = daysPct
		}
		if progressPercent < 0 {
			progressPercent = 0
		}
		if progressPercent > 100 {
			progressPercent = 100
		}

		var status UpcomingMaintenanceStatus
		switch {
		case remainingKm <= 0 || remainingDays <= 0:
			status = StatusOverdue
			overdueCount++
		case remainingKm <= UpcomingKmThreshold || remainingDays <= UpcomingDaysThreshold:
			status = StatusUpcoming
			upcomingCount++
		default:
			status = StatusOK
			okCount++
		}

		alertEligible := remainingKm <= AlertKmThreshold || remainingDays <= AlertDaysThreshold

		items = append(items, UpcomingMaintenanceItem{
			Type:            rule.Type,
			Label:           rule.Label,
			Status:          status,
			IntervalKm:      intervalKm,
			IntervalMonths:  intervalMonths,
			HasHistory:      hasHistory,
			LastPerformedAt: lastPerformedAt,
			LastMileage:     lastMileage,
			NextDueMileage:  nextDueMileage,
			NextDueDate:     nextDueDate,
			RemainingKm:     remainingKm,
			RemainingDays:   remainingDays,
			ProgressPercent: progressPercent,
			AlertEligible:   alertEligible,
		})
	}

	overallStatus := StatusOK
	if overdueCount > 0 {
		overallStatus = StatusOverdue
	} else if upcomingCount > 0 {
		overallStatus = StatusUpcoming
	}

	vehicleType := strings.TrimSpace(car.VehicleType)
	if vehicleType == "" {
		vehicleType = "cars"
	}

	return CarUpcomingMaintenances{
		CarID:         car.ID,
		CarName:       car.Name,
		Manufacturer:  car.Manufacturer,
		Model:         car.Model,
		VehicleType:   vehicleType,
		CurrentKm:     car.LastMileage,
		Status:        overallStatus,
		OKCount:       okCount,
		UpcomingCount: upcomingCount,
		OverdueCount:  overdueCount,
		Items:         items,
	}
}

func resolveBaselineMileage(car domaincar.Car, maintenances []Maintenance) int {
	if car.InitialMileage > 0 {
		return car.InitialMileage
	}
	if len(maintenances) > 0 {
		minMileage := car.LastMileage
		for _, m := range maintenances {
			if m.Mileage >= 0 && (minMileage == 0 || m.Mileage < minMileage) {
				minMileage = m.Mileage
			}
		}
		return minMileage
	}
	return car.LastMileage
}
