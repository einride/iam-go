CREATE TABLE shippers (
    shipper_id STRING(63) NOT NULL,
    create_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    update_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    delete_time TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    display_name STRING(63) NOT NULL,
) PRIMARY KEY(shipper_id);

CREATE TABLE sites (
    shipper_id STRING(63) NOT NULL,
    site_id STRING(63) NOT NULL,
    create_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    update_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    delete_time TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    display_name STRING(63),
    latitude FLOAT64,
    longitude FLOAT64,
) PRIMARY KEY(shipper_id, site_id);

CREATE TABLE shipments (
    shipper_id STRING(63) NOT NULL,
    shipment_id STRING(63) NOT NULL,
    create_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    update_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    delete_time TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    origin_site_id STRING(63),
    destination_site_id STRING(63),
    pickup_earliest_time TIMESTAMP,
    pickup_latest_time TIMESTAMP,
    delivery_earliest_time TIMESTAMP,
    delivery_latest_time TIMESTAMP,
) PRIMARY KEY(shipper_id, shipment_id);

CREATE TABLE line_items (
    shipper_id STRING(63) NOT NULL,
    shipment_id STRING(63) NOT NULL,
    line_number INT64 NOT NULL,
    title STRING(63),
    quantity FLOAT64,
    weight_kg FLOAT64,
    volume_m3 FLOAT64,
) PRIMARY KEY(shipper_id, shipment_id, line_number ASC),
  INTERLEAVE IN PARENT shipments ON DELETE CASCADE;
