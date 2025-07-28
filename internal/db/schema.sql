CREATE TABLE IF NOT EXISTS aws_ec2_pricing (
    sku TEXT PRIMARY KEY,
    instance_type TEXT,
    region TEXT,
    operating_system TEXT,
    tenancy TEXT,
    price_per_hour NUMERIC,
    raw_attributes JSONB
);
