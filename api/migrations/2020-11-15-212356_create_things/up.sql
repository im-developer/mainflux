create table things
(
	name varchar(255) constraint things_pk primary key,
	secret varchar(255),
	user_secret varchar(255) not null constraint things_users_secret_fk references users (secret) on delete cascade
);

create unique index things_secret_uindex on things (secret);
create unique index things_name_user_secret_uindex
	on things (name, user_secret);