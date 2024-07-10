cqlsh localhost 9042


CREATE KEYSPACE auth WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 3};

use auth;
create table auth.user (username text, id uuid, createdat timestamp, password text, email text, PRIMARY KEY (id));

CREATE INDEX ON auth.user (email);
