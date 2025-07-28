#!/bin/bash

# CONFIG
DB_NAME=""
DB_USER=""
DB_PASSWORD=""
TABLE_NAME="aws_ec2_pricing"

# Step 1: Create DB and User
psql -U postgres <<EOF
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

# Step 2: Create Table
psql -U $DB_USER -d $DB_NAME <<EOF
CREATE TABLE IF NOT EXISTS $TABLE_NAME (
    sku TEXT PRIMARY KEY,
    instance_type TEXT,
    region TEXT,
    operating_system TEXT,
    tenancy TEXT,
    price_per_hour FLOAT,
    currency TEXT DEFAULT 'USD',
    raw_attributes JSONB,
    last_updated TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_instance_type ON $TABLE_NAME(instance_type);
CREATE INDEX IF NOT EXISTS idx_region ON $TABLE_NAME(region);
EOF

echo "✅ Database and table created!"


