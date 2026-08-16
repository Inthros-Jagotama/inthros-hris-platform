-- =============================================================================
-- 143_competency_rating_scales.sql
-- Competency 360 Module — Rating Scale (plan generik §7).
-- Rating scale reusable: satu scale (mis. "Skala 1-5") dengan item-item
-- (1 = Sangat Tidak Memenuhi ... 5 = Sangat Baik), dipakai oleh template
-- assessment, indicator, aggregation, dan report.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 143.1 Competency Rating Scales
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_rating_scales (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(50)  NOT NULL,
    description TEXT         NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_by  CHAR(36)     NULL,
    updated_by  CHAR(36)     NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_comp_scale_code ON competency_rating_scales (code);

CREATE INDEX IF NOT EXISTS idx_comp_scale_status ON competency_rating_scales (status);

-- ---------------------------------------------------------------------------
-- 143.2 Competency Rating Scale Items
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS competency_rating_scale_items (
    id          CHAR(36) PRIMARY KEY,
    scale_id    CHAR(36)     NOT NULL,
    value       SMALLINT     NOT NULL,
    label       VARCHAR(255) NOT NULL,
    description TEXT         NULL,
    weight      DECIMAL(6,2) NOT NULL DEFAULT 1,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_comp_scale_item_value UNIQUE (scale_id, value),

    CONSTRAINT fk_comp_scale_item_scale FOREIGN KEY (scale_id) REFERENCES competency_rating_scales(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comp_scale_item_scale ON competency_rating_scale_items (scale_id);
