create table packages
(
	id serial not null
		constraint packages_pk
			primary key,
	name varchar not null
);

alter table packages owner to postgres;

create unique index packages_id_uindex
	on packages (id);

create unique index packages_name_uindex
	on packages (name);

create table versions
(
	id serial not null
		constraint versions_pk
			primary key,
	name varchar not null,
	size integer not null,
	checksum varchar not null,
	file_name varchar not null,
	package_id integer not null
		constraint versions_packages_id_fk
			references packages
				on delete cascade
);

alter table versions owner to postgres;

create unique index versions_id_uindex
	on versions (id);

create unique index versions_package_id_name_uindex
	on versions (package_id, name);