// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package carrierleads

import (
	"database/sql"
	"encoding/json"
)

type FmcsaCarrierSafer struct {
	EntityType                    string
	OperatingStatus               string
	OosDate                       sql.NullTime
	LegalName                     string
	DbaName                       string
	Address                       string
	Telephone                     string
	MailingAddress                string
	DotNumber                     int32
	StateCarrierIDNumber          string
	DocketNumber                  string
	DunsNumber                    string
	PowerUnits                    int32
	Drivers                       int32
	Mcs150FormDate                sql.NullTime
	Mcs150MileageYear             string
	CarrierOperation              string
	OcAuthorizedForHire           bool
	OcPrivatePassengerBusiness    bool
	OcUsMail                      bool
	OcLocalGovernment             bool
	OcExemptForHire               bool
	OcPrivatePassengerNonBusiness bool
	OcFederalGovernment           bool
	OcIndianTribe                 bool
	OcPrivateProperty             bool
	OcMigrant                     bool
	OcStateGovernment             bool
	OcOther                       string
	CcGeneralFreight              bool
	CcMotorVehicles               bool
	CcBuildingMaterials           bool
	CcFreshProduct                bool
	CcPassengers                  bool
	CcGrainFeedHay                bool
	CcGarbageRefuseTrash          bool
	CcCommoditiesDryBulk          bool
	CcPaperProducts               bool
	CcConstruction                bool
	CcHouseholdGoods              bool
	CcDriveAwayTowaway            bool
	CcMobileHomes                 bool
	CcLiquidsGases                bool
	CcOilfieldEquipment           bool
	CcCoalCoke                    bool
	CcUsMail                      bool
	CcRefrigeratedFood            bool
	CcUtility                     bool
	CcWaterwell                   bool
	CcMetalSheetsCoilsRolls       bool
	CcLogsPolesBeamsLumber        bool
	CcMachineryLargeObjects       bool
	CcIntermodalContainers        bool
	CcLivestock                   bool
	CcMeat                        bool
	CcChemicals                   bool
	CcBeverages                   bool
	CcFarmSupplies                bool
	CcOther                       string
	UsInspectionVehicle           json.RawMessage
	UsInspectionDriver            json.RawMessage
	UsInspectionHazmat            json.RawMessage
	UsInspectionIep               json.RawMessage
	UsCrashSummary                json.RawMessage
	CanInspectionVehicle          json.RawMessage
	CanInspectionDriver           json.RawMessage
	CanCrashSummary               json.RawMessage
	SafetyRatingDate              sql.NullTime
	SafetyRatingReviewDate        sql.NullTime
	SafetyRating                  string
	SafetyRatingType              string
	LatestUpdateTime              sql.NullTime
	CreatedAt                     int32
}
