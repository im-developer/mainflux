create table channels_things
(
	channel varchar(255) not null constraint channels_things_channels_secret_fk references channels (name) on delete cascade,
	thing varchar(255) not null constraint channels_things_things_secret_fk references things (name) on delete cascade
);