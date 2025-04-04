CREATE TABLE qrcode_data (
    id BIGSERIAL PRIMARY KEY,
    qrcode VARCHAR(64) NOT NULL,
    to_url VARCHAR(255) NOT NULL,
    pv INTEGER NOT NULL DEFAULT 0, -- 总访问量
    created_at BIGINT NOT NULL, -- 创建时间
    updated_at BIGINT NOT NULL, -- 更新时间
    status SMALLINT NOT NULL DEFAULT 10,
    CONSTRAINT uk_qrcode UNIQUE (qrcode)
);

CREATE TABLE qrcode_data_query_log (
    id BIGSERIAL PRIMARY KEY,
    qrcode VARCHAR(64) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    request_headers TEXT NOT NULL,
    request_ip VARCHAR(30) NOT NULL,
    created_at BIGINT NOT NULL, -- 创建时间
    updated_at BIGINT NOT NULL, -- 更新时间
    CONSTRAINT idx_qrcode_query_log UNIQUE (qrcode)
);
