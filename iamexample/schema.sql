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
    display_name STRING(63) NOT NULL,
    latitude FLOAT64,
    longitude FLOAT64,
    CONSTRAINT fk_sites_parent
        FOREIGN KEY (shipper_id)
        REFERENCES shippers (shipper_id),
) PRIMARY KEY(shipper_id, site_id);

CREATE TABLE shipments (
    shipper_id STRING(63) NOT NULL,
    shipment_id STRING(63) NOT NULL,
    create_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    update_time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    delete_time TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    origin_site_id STRING(63) NOT NULL,
    destination_site_id STRING(63) NOT NULL,
    pickup_earliest_time TIMESTAMP NOT NULL,
    pickup_latest_time TIMESTAMP NOT NULL,
    delivery_earliest_time TIMESTAMP NOT NULL,
    delivery_latest_time TIMESTAMP NOT NULL,
    annotations ARRAY<STRING(MAX)> NOT NULL,
    CONSTRAINT fk_shipments_parent
        FOREIGN KEY (shipper_id)
        REFERENCES shippers (shipper_id),
    CONSTRAINT fk_shipments_origin_site
        FOREIGN KEY (shipper_id, origin_site_id)
        REFERENCES sites (shipper_id, site_id),
    CONSTRAINT fk_shipments_destination_site
        FOREIGN KEY (shipper_id, destination_site_id)
        REFERENCES sites (shipper_id, site_id),
) PRIMARY KEY(shipper_id, shipment_id);

CREATE TABLE line_items (
    shipper_id STRING(63) NOT NULL,
    shipment_id STRING(63) NOT NULL,
    line_number INT64 NOT NULL,
    title STRING(63) NOT NULL,
    quantity FLOAT64 NOT NULL,
    weight_kg FLOAT64 NOT NULL,
    volume_m3 FLOAT64 NOT NULL,
) PRIMARY KEY(shipper_id, shipment_id, line_number ASC),
  INTERLEAVE IN PARENT shipments ON DELETE CASCADE;
