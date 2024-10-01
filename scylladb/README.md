* Start Scylladb on Docker
```
  docker run --rm -it -p 9042:9042 scylladb/scylla --smp 2
```
* Create Keyspaces on Scylladb
```
    CREATE KEYSPACE Pets_Clinic WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor' : 1};
```
* Create Table
```
    CREATE TABLE IF NOT EXISTS pets_clinic.books (book_id int primary key, book_name text);
```