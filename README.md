# SAFER Crawler

A web crawler to scrape data from the U.S. Department of Transportation Safety and Fitness Electronic Records System ([SAFER](https://safer.fmcsa.dot.gov/)), and store the structured result into a mysql database.

![Screenshot 2022-10-30 at 11 35 54 AM](https://user-images.githubusercontent.com/1932380/198890241-691d379c-cd3f-4c6b-b249-861b10a5d859.png)

## Set up the database

1. Create a table in the mysql database with the following schema

```
CREATE TABLE `fmcsa_carrier_safer` (
    `entity_type` varchar(255) NOT NULL,
    `operating_status` varchar(255) NOT NULL,
    `oos_date` date DEFAULT NULL,
    `legal_name` varchar(255) NOT NULL,
    `dba_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `address` varchar(255) NOT NULL,
    `telephone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `mailing_address` varchar(255) NOT NULL,
    `dot_number` int NOT NULL,
    `state_carrier_id_number` varchar(255) NOT NULL,
    `docket_number` varchar(255) NOT NULL,
    `duns_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `power_units` int NOT NULL,
    `drivers` int NOT NULL,
    `mcs_150_form_date` date DEFAULT NULL,
    `mcs_150_mileage_year` varchar(255) NOT NULL,
    `carrier_operation` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `oc_authorized_for_hire` tinyint(1) NOT NULL,
    `oc_private_passenger_business` tinyint(1) NOT NULL,
    `oc_us_mail` tinyint(1) NOT NULL,
    `oc_local_government` tinyint(1) NOT NULL,
    `oc_exempt_for_hire` tinyint(1) NOT NULL,
    `oc_private_passenger_non_business` tinyint(1) NOT NULL,
    `oc_federal_government` tinyint(1) NOT NULL,
    `oc_indian_tribe` tinyint(1) NOT NULL,
    `oc_private_property` tinyint(1) NOT NULL,
    `oc_migrant` tinyint(1) NOT NULL,
    `oc_state_government` tinyint(1) NOT NULL,
    `oc_other` varchar(255) NOT NULL,
    `cc_general_freight` tinyint(1) NOT NULL,
    `cc_motor_vehicles` tinyint(1) NOT NULL,
    `cc_building_materials` tinyint(1) NOT NULL,
    `cc_fresh_product` tinyint(1) NOT NULL,
    `cc_passengers` tinyint(1) NOT NULL,
    `cc_grain_feed_hay` tinyint(1) NOT NULL,
    `cc_garbage_refuse_trash` tinyint(1) NOT NULL,
    `cc_commodities_dry_bulk` tinyint(1) NOT NULL,
    `cc_paper_products` tinyint(1) NOT NULL,
    `cc_construction` tinyint(1) NOT NULL,
    `cc_household_goods` tinyint(1) NOT NULL,
    `cc_drive_away_towaway` tinyint(1) NOT NULL,
    `cc_mobile_homes` tinyint(1) NOT NULL,
    `cc_liquids_gases` tinyint(1) NOT NULL,
    `cc_oilfield_equipment` tinyint(1) NOT NULL,
    `cc_coal_coke` tinyint(1) NOT NULL,
    `cc_us_mail` tinyint(1) NOT NULL,
    `cc_refrigerated_food` tinyint(1) NOT NULL,
    `cc_utility` tinyint(1) NOT NULL,
    `cc_waterwell` tinyint(1) NOT NULL,
    `cc_metal_sheets_coils_rolls` tinyint(1) NOT NULL,
    `cc_logs_poles_beams_lumber` tinyint(1) NOT NULL,
    `cc_machinery_large_objects` tinyint(1) NOT NULL,
    `cc_intermodal_containers` tinyint(1) NOT NULL,
    `cc_livestock` tinyint(1) NOT NULL,
    `cc_meat` tinyint(1) NOT NULL,
    `cc_chemicals` tinyint(1) NOT NULL,
    `cc_beverages` tinyint(1) NOT NULL,
    `cc_farm_supplies` tinyint(1) NOT NULL,
    `cc_other` varchar(255) NOT NULL,
    `us_inspection_vehicle` json NOT NULL,
    `us_inspection_driver` json NOT NULL,
    `us_inspection_hazmat` json NOT NULL,
    `us_inspection_iep` json NOT NULL,
    `us_crash_summary` json NOT NULL,
    `can_inspection_vehicle` json NOT NULL,
    `can_inspection_driver` json NOT NULL,
    `can_crash_summary` json NOT NULL,
    `safety_rating_date` date DEFAULT NULL,
    `safety_rating_review_date` date DEFAULT NULL,
    `safety_rating` varchar(255) NOT NULL,
    `safety_rating_type` varchar(255) NOT NULL,
    `latest_update_time` date DEFAULT NULL,
    `created_at` int NOT NULL,
    PRIMARY KEY (`dot_number`)
);
```

2. Edit the `config.yaml` file and replace the default sql connection string to point to your mysql instance.

## Running the crawler

```
  go run main.go
```

## Performance

As on 10/29/2022, running the crawler on a digital ocean droplet with 2GM ram and 1 AMD vCPU, the cralwer finished in 18 hours and scraped 2,022,837 records. The result database size is at 929 MiB.

## Credit

The fetching and parsing of individual safer data depends on brandenc40's [safer](https://github.com/brandenc40/safer) project.
