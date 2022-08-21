create database mockidoki
	with owner postgres;

\c mockidoki

create table http_mock
(
    id              serial                not null
        constraint http_mock_pk
            primary key,
    method          varchar(12)           not null,
    matching_url    varchar(512),
    matching_header varchar(256),
    matching_body   varchar(10000),
    response_status integer               not null,
    response_body   varchar(5000),
    response_header jsonb,
    description     varchar(256),
    is_deleted      boolean default false not null
);

alter table http_mock
    owner to postgres;

create unique index http_mock_method_murl_mheader_mbody_unique_index
    on http_mock (method, matching_url, matching_header, matching_body);

INSERT INTO public.http_mock (method, matching_url, matching_header, matching_body, response_status, response_body, response_header, description, is_deleted) VALUES ('POST', '(.*)/warehouses/3/stock-transactions(.*)', '(.*)mockidoki-user(.*)', '(.*)(?=.*12345)(?=.*45.37)(.*)', 201, null, '[{"value": "test_perfect1", "header": "test_header1"}, {"value": "test_perfect2", "header": "test_header2"}]', 'Mock for adding new stock transaction', false);
INSERT INTO public.http_mock (method, matching_url, matching_header, matching_body, response_status, response_body, response_header, description, is_deleted) VALUES ('GET', '(.*)/warehouses/3/stock-transactions\?type=SALES(.*)', null, null, 200, '{"transactions":[{"Id":"6259166211f5d6001565a5dc","ettn":"462B3870-8F69-4FC6-B5E3-7D1E6D47BF25","iNumber":"4272019563da4174b25e4ccc059a75ec","status":"OK","createdDate":"2022-08-15T09:53:22","amount":377.97},{"Id":"6239566211f5d6892465a5df","ettn":"566B3870-8F69-4AC6-B5E3-7B1E6D47BF68","iNumber":"5679019563da4174b25e4ccc059a75fd","status":"OK","createdDate":"2022-08-16T10:57:26","amount":67.15},{"Id":"1439566211f5d6892465a5ab","ettn":"986B3870-8F69-4AF6-A5D3-7B1E6D47BF11","iNumber":"1879019563da4074c25e4ccc059a75ed","status":"OK","createdDate":"2022-08-17T12:21:12","amount":90.25}]}', null, 'Mock stock transaction list', false);
INSERT INTO public.http_mock (method, matching_url, matching_header, matching_body, response_status, response_body, response_header, description, is_deleted) VALUES ('POST', '(.*)/product(.*)', null, '(.*)991188|17789|828929(.*)', 201, null, null, 'Mock for adding new product (if productId is 991188 or 17789 or 828929)', false);