create table channels
(
	name varchar(255) constraint channels_pk primary key,
	secret varchar(255),
	user_secret varchar(255) not null constraint channels_users_secret_fk references users (secret) on delete cascade
);

create unique index channels_secret_uindex on channels (secret);
create unique index channels_name_user_secret_uindex
	on channels (name, user_secret);