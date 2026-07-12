CREATE TABLE IF NOT EXISTS newsletter_subscribe.users
(
    id           uuid primary key,
    created_at   timestamptz  not null,
    updated_at   timestamptz  not null,
    deleted_at   timestamptz,
    first_name   varchar(100) not null,
    last_name    varchar(100) not null,
    patronymic   varchar(100),
    phone_number varchar(15),
    status       integer      not null,

    constraint users_first_name_capitalized check (first_name ~ '^[А-ЯЁA-Z]'),
    constraint users_last_name_capitalized check (last_name ~ '^[А-ЯЁA-Z]'),
    constraint users_patronymic_capitalized check (patronymic is null or patronymic ~ '^[А-ЯЁA-Z]'),
    constraint users_phone_number_format check (phone_number is null or phone_number ~ '^\+[0-9]{10,15}$')
);
