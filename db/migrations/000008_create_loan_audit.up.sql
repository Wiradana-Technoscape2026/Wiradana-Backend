CREATE TABLE loan_audit_log (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  loan_id uuid NOT NULL REFERENCES loan(id),
  action text NOT NULL,
  performed_by uuid NOT NULL REFERENCES app_user(id),
  performed_at timestamptz NOT NULL DEFAULT now(),
  before_data jsonb,
  after_data jsonb,
  note text,
  is_flagged boolean NOT NULL DEFAULT false,
  flagged_by_name text,
  flagged_at timestamptz,
  flagged_reason text
);

CREATE TABLE loan_audit_token (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id uuid NOT NULL REFERENCES cooperative(id),
  loan_id uuid NOT NULL REFERENCES loan(id),
  token_hash text NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  created_by uuid NOT NULL REFERENCES app_user(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  revoked boolean NOT NULL DEFAULT false
);

CREATE INDEX idx_loan_audit_log_loan ON loan_audit_log(loan_id);
CREATE INDEX idx_loan_audit_token_hash ON loan_audit_token(token_hash);
