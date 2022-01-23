create if not exists table event_mock
(
    id bigserial not null
        constraint event_mock_pkey
            primary key,
    key varchar(255) not null,
    channel varchar(255) not null,
    description varchar(255) not null,
    is_deleted boolean
);

create if not exists unique index event_mock_unique_key
	on action (key);

create if not exists table http_mock
(
    id serial not null
        constraint http_mock_pk
            primary key,
    method varchar(12) not null,
    matching_url varchar(512),
    matching_header varchar(256),
    matching_body varchar(10000),
    response_status int not null,
    response_body varchar(5000),
    response_header jsonb,
    description varchar(256),
    is_deleted bool default false not null
);

create if not exists unique index http_mock_method_murl_mheader_mbody_unique_index
	on http_mock (method, matching_url, matching_header, matching_body);

--insert into http_mock (id, method, matching_url, matching_header, matching_body, response_status, response_body, response_header, description, is_deleted)
--values (1, 'GET', '1923', null, null, 200, '{"status": "DONE"}', '[{"header":"test_header1","value":"test_value1"},{"header":"test_header2","value":"test_value2"}]', 'test mock', false);