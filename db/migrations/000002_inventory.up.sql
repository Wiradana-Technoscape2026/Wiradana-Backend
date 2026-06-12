-- Inventory module tables (Tier 3)
CREATE TABLE inventory_field_def (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  field_key text NOT NULL,
  label text NOT NULL,
  data_type text NOT NULL DEFAULT 'text' CHECK (data_type IN ('text','number','select','date')),
  options jsonb NOT NULL DEFAULT '[]',
  required boolean NOT NULL DEFAULT false,
  sort_order int NOT NULL DEFAULT 0,
  UNIQUE (cooperative_id, field_key)
);

CREATE TABLE inventory_product (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  sku text NOT NULL,
  name text NOT NULL,
  unit text NOT NULL DEFAULT 'pcs',
  custom_attributes jsonb NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cooperative_id, sku)
);

CREATE TABLE inventory_movement (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  product_id uuid NOT NULL REFERENCES inventory_product(id),
  direction text NOT NULL CHECK (direction IN ('masuk','keluar')),
  quantity bigint NOT NULL CHECK (quantity > 0),
  note text,
  custom_attributes jsonb NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_movement_product ON inventory_movement(product_id);
CREATE INDEX idx_product_coop ON inventory_product(cooperative_id);
