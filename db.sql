create table http_mock
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

INSERT INTO public.http_mock (id, method, matching_url, matching_header, matching_body, response_status, response_body, response_header, description, is_deleted)
VALUES (1, 'GET', '1923', null, null, 200, '{"status": "DONE"}', '[{"header":"test_header1","value":"test_value1"},{"header":"test_header2","value":"test_value2"}]', 'test mock', false);