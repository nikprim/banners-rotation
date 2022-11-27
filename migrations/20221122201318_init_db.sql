-- +goose Up
create
extension if not exists "uuid-ossp";

CREATE TABLE slots
(
    guid uuid NOT NULL,
    name text NOT NULL,
    CONSTRAINT unq_slots_guid UNIQUE (guid),
    CONSTRAINT pk_slots_guid PRIMARY KEY (guid)
);

CREATE TABLE banners
(
    guid uuid NOT NULL,
    name text NOT NULL,
    CONSTRAINT unq_banners_guid UNIQUE (guid),
    CONSTRAINT pk_banners_guid PRIMARY KEY (guid)
);

CREATE TABLE social_groups
(
    guid uuid NOT NULL,
    name text NOT NULL,
    CONSTRAINT unq_social_groups_guid UNIQUE (guid),
    CONSTRAINT pk_social_groups_guid PRIMARY KEY (guid)
);

CREATE TABLE banners_link_slots
(
    guid        uuid NOT NULL,
    banner_guid uuid NOT NULL,
    slot_guid   uuid NOT NULL,
    CONSTRAINT unq_banners_link_slots_guid UNIQUE (guid),
    CONSTRAINT unq_banners_link_slots_row UNIQUE (banner_guid, slot_guid),
    CONSTRAINT pk_banners_link_slots_guid PRIMARY KEY (guid),
    CONSTRAINT fk_banners_link_slots_banner_guid FOREIGN KEY (banner_guid) REFERENCES banners (guid),
    CONSTRAINT fk_banners_link_slots_slot_guid FOREIGN KEY (slot_guid) REFERENCES slots (guid)
);

CREATE TABLE stats
(
    guid              uuid NOT NULL,
    banner_guid       uuid NOT NULL,
    slot_guid         uuid NOT NULL,
    social_group_guid uuid NOT NULL,
    shows             int  NOT NULL,
    clicks            int  NOT NULL,
    CONSTRAINT unq_stats_guid UNIQUE (guid),
    CONSTRAINT unq_stats_row UNIQUE (banner_guid, slot_guid, social_group_guid),
    CONSTRAINT pk_stats_guid PRIMARY KEY (guid),
    CONSTRAINT fk_stats_banner_guid FOREIGN KEY (banner_guid) REFERENCES banners (guid),
    CONSTRAINT fk_stats_slot_guid FOREIGN KEY (slot_guid) REFERENCES slots (guid),
    CONSTRAINT fk_stats_social_group_guid FOREIGN KEY (social_group_guid) REFERENCES social_groups (guid)
);

INSERT INTO slots(guid, name)
VALUES ('00000000-0000-0000-0000-000000000001', 'slot 1'),
       ('00000000-0000-0000-0000-000000000002', 'slot 2'),
       ('00000000-0000-0000-0000-000000000003', 'slot 3'),
       ('00000000-0000-0000-0000-000000000004', 'slot 4'),
       ('00000000-0000-0000-0000-000000000005', 'slot 5'),
       ('00000000-0000-0000-0000-000000000006', 'slot 6');

INSERT INTO banners(guid, name)
VALUES ('00000000-0000-0000-1111-000000000001', 'banner 1'),
       ('00000000-0000-0000-1111-000000000002', 'banner 2'),
       ('00000000-0000-0000-1111-000000000003', 'banner 3'),
       ('00000000-0000-0000-1111-000000000004', 'banner 4'),
       ('00000000-0000-0000-1111-000000000005', 'banner 5');

INSERT INTO social_groups(guid, name)
VALUES ('00000000-0000-0000-2222-000000000001', 'social group 1'),
       ('00000000-0000-0000-2222-000000000002', 'social group 2'),
       ('00000000-0000-0000-2222-000000000003', 'social group 3'),
       ('00000000-0000-0000-2222-000000000004', 'social group 4'),
       ('00000000-0000-0000-2222-000000000005', 'social group 5');

INSERT INTO banners_link_slots(guid, banner_guid, slot_guid)
VALUES (uuid_generate_v4(), '00000000-0000-0000-1111-000000000001', '00000000-0000-0000-0000-000000000001'),

       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000001', '00000000-0000-0000-0000-000000000002'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000002', '00000000-0000-0000-0000-000000000002'),

       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000001', '00000000-0000-0000-0000-000000000003'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000002', '00000000-0000-0000-0000-000000000003'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000003', '00000000-0000-0000-0000-000000000003'),

       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000001', '00000000-0000-0000-0000-000000000004'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000002', '00000000-0000-0000-0000-000000000004'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000003', '00000000-0000-0000-0000-000000000004'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000004', '00000000-0000-0000-0000-000000000004'),

       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000001', '00000000-0000-0000-0000-000000000005'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000002', '00000000-0000-0000-0000-000000000005'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000003', '00000000-0000-0000-0000-000000000005'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000004', '00000000-0000-0000-0000-000000000005'),
       (uuid_generate_v4(), '00000000-0000-0000-1111-000000000005', '00000000-0000-0000-0000-000000000005');

-- +goose Down
drop
extension if exists "uuid-ossp";

DROP TABLE stats;
DROP TABLE banners_link_slots;
DROP TABLE social_groups;
DROP TABLE banners;
DROP TABLE slots;