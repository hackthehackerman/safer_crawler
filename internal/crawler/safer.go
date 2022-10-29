package crawler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"carrierleads.com/internal/dao"
	"carrierleads.com/internal/dao/carrierleads"
	"carrierleads.com/internal/lib/safer"
)

func CrawlSafer(DOTWatermark, bucketSize, numConnections int, dao dao.Dao) (err error) {
	ctx := context.Background()
	var wg sync.WaitGroup
	c := make(chan func())
	poll := func() {
		defer wg.Done()
		for {
			t1, more := <-c
			if more {
				t1()
			} else {
				return
			}
		}
	}

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go poll()
	}

	// map to store number of failed safer fetch in dot range segmented by bucket size.
	// this is used to detect the upperfound of available dotnumber
	notFound := sync.Map{}
	terminate := false
	dot := 1
	for {
		if dot%500 == 0 {
			log.Println("processing", dot)
		}
		if terminate {
			log.Println("upperbound reached, terminating")
			break
		}
		dotNumber := dot

		c <- func() {
			s, err := getSaferSnapshot(int(dotNumber))
			if err != nil {
				log.Print("failed to get snapshot ", dotNumber, err)

				bucket := dotNumber / bucketSize
				cnt, ok := notFound.Load(bucket)
				if ok {
					notFound.Store(bucket, cnt.(int)+1)
					if cnt.(int)+1 == bucketSize && dotNumber > DOTWatermark {
						terminate = true
					}
				} else {
					notFound.Store(bucket, 1)
				}
				return
			}

			err = writeToDB(s, dotNumber, ctx, dao)
			if err != nil {
				log.Printf("failed to write result for %d into db %v", dotNumber, err)
			}
		}
		dot += 1
	}
	close(c)
	wg.Wait()
	return
}

func getSaferSnapshot(dotNumber int) (ret *safer.CompanySnapshot, err error) {
	var backoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}

	for _, backoff := range backoffSchedule {
		client := safer.NewClient()
		ret, err = client.GetCompanyByDOTNumber(strconv.Itoa(dotNumber))
		if err == nil {
			return
		}
		if err.Error() == "company not found" {
			return
		}
		log.Printf("safer api error %d %+v \n", dotNumber, err)
		log.Printf("retrying %d in %v", dotNumber, backoff)
		time.Sleep(backoff)
	}
	return
}

func writeToDB(s *safer.CompanySnapshot, dotNumber int, ctx context.Context, dao dao.Dao) (err error) {
	buildNullTime := func(in *time.Time) sql.NullTime {
		if in == nil {
			return sql.NullTime{Valid: false}
		} else {
			return sql.NullTime{Time: *in, Valid: true}
		}
	}

	buildInspectionSummary := func(in safer.InspectionSummary) json.RawMessage {
		ret, _ := json.Marshal([]any{in.Inspections, in.OutOfService, in.OutOfServicePct, in.NationalAverage})
		return json.RawMessage(ret)
	}

	buildCrashSummary := func(in safer.CrashSummary) json.RawMessage {
		ret, _ := json.Marshal([]any{in.Fatal, in.Injury, in.Tow, in.Total})
		return json.RawMessage(ret)
	}

	newDotNumber, err := strconv.Atoi(s.DOTNumber)
	if err != nil {
		log.Print("failed to convert dotnumber", s, dotNumber)
		return
	}

	params := carrierleads.CreateSaferSnapshotParams{
		EntityType:             s.EntityType,
		OperatingStatus:        s.OperatingStatus,
		OosDate:                buildNullTime(s.OutOfServiceDate),
		LegalName:              s.LegalName,
		DbaName:                s.DBAName,
		Address:                s.PhysicalAddress,
		Telephone:              s.Phone,
		MailingAddress:         s.MailingAddress,
		DotNumber:              int32(newDotNumber),
		StateCarrierIDNumber:   s.StateCarrierID,
		DocketNumber:           strings.Join(s.MCMXFFNumbers, ","),
		DunsNumber:             s.DUNSNumber,
		PowerUnits:             int32(s.PowerUnits),
		Drivers:                int32(s.Drivers),
		Mcs150FormDate:         buildNullTime(s.MCS150FormDate),
		Mcs150MileageYear:      s.MCS150Year,
		CarrierOperation:       strings.Join(s.CarrierOperation, ","),
		UsInspectionVehicle:    buildInspectionSummary(s.USVehicleInspections),
		UsInspectionDriver:     buildInspectionSummary(s.USDriverInspections),
		UsInspectionHazmat:     buildInspectionSummary(s.USHazmatInspections),
		UsInspectionIep:        buildInspectionSummary(s.USIEPInspections),
		UsCrashSummary:         buildCrashSummary(s.USCrashes),
		CanInspectionVehicle:   buildInspectionSummary(s.CanadaVehicleInspections),
		CanInspectionDriver:    buildInspectionSummary(s.CanadaDriverInspections),
		CanCrashSummary:        buildCrashSummary(s.CanadaCrashes),
		SafetyRatingDate:       buildNullTime(s.Safety.RatingDate),
		SafetyRatingReviewDate: buildNullTime(s.Safety.ReviewDate),
		SafetyRating:           s.Safety.Rating,
		SafetyRatingType:       s.Safety.Type,
		LatestUpdateTime:       buildNullTime(s.LatestUpdateDate),
		CreatedAt:              int32(time.Now().Unix()),
	}

	for _, oc := range s.OperationClassification {
		switch oc {
		case "Auth. For Hire":
			params.OcAuthorizedForHire = true
		case "Priv. Pass.(Non-business)":
			params.OcPrivatePassengerNonBusiness = true
		case "State Gov't":
			params.OcStateGovernment = true
		case "Exempt For Hire":
			params.OcExemptForHire = true
		case "Migrant":
			params.OcMigrant = true
		case "Local Gov't":
			params.OcLocalGovernment = true
		case "Private(Property)":
			params.OcPrivateProperty = true
		case "U.S. Mail":
			params.OcUsMail = true
		case "Indian Nation":
			params.OcIndianTribe = true
		case "Priv. Pass. (Business)":
			params.OcPrivatePassengerBusiness = true
		case "Fed. Gov't":
			params.OcFederalGovernment = true
		default:
			params.OcOther = oc
		}
	}
	for _, cc := range s.CargoCarried {
		switch cc {
		case "General Freight":
			params.CcGeneralFreight = true
		case "Liquids/Gases":
			params.CcLiquidsGases = true
		case "Chemicals":
			params.CcChemicals = true

		case "Household Goods":
			params.CcHouseholdGoods = true
		case "Intermodal Cont.":
			params.CcIntermodalContainers = true
		case "Commodities Dry Bulk":
			params.CcCommoditiesDryBulk = true

		case "Metal: sheets, coils, rolls":
			params.CcMetalSheetsCoilsRolls = true
		case "Passengers":
			params.CcPassengers = true
		case "Refrigerated Food":
			params.CcRefrigeratedFood = true

		case "Motor Vehicles":
			params.CcMotorVehicles = true
		case "Oilfield Equipment":
			params.CcOilfieldEquipment = true
		case "Beverages":
			params.CcBeverages = true

		case "Drive/Tow away":
			params.CcDriveAwayTowaway = true
		case "Livestock":
			params.CcLivestock = true
		case "Paper Products":
			params.CcPaperProducts = true

		case "Logs, Poles, Beams, Lumber":
			params.CcLogsPolesBeamsLumber = true
		case "Grain, Feed, Hay":
			params.CcGrainFeedHay = true
		case "Utilities":
			params.CcUtility = true

		case "Building Materials":
			params.CcBuildingMaterials = true
		case "Coal/Coke":
			params.CcCoalCoke = true
		case "Agricultural/Farm Supplies":
			params.CcFarmSupplies = true

		case "Mobile Homes":
			params.CcMobileHomes = true
		case "Meat":
			params.CcMeat = true
		case "Construction":
			params.CcConstruction = true

		case "Machinery, Large Objects":
			params.CcMachineryLargeObjects = true
		case "Garbage/Refuse":
			params.CcGarbageRefuseTrash = true
		case "Water Well":
			params.CcWaterwell = true

		case "Fresh Produce":
			params.CcFreshProduct = true
		case "US Mail":
			params.CcUsMail = true
		default:
			params.CcOther = cc
		}
	}
	_, err = dao.Queries.CreateSaferSnapshot(ctx, params)
	return
}
