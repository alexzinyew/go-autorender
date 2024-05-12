DROP TABLE IF EXISTS videos;

CREATE TABLE videos (
    id varchar(8) PRIMARY KEY,
	status int,
	
	title varchar(100)
);

select * from videos