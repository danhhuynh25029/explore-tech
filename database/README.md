* Isolation level
  * Read Uncommitted
    * Transaction 1 thay đổi dữ liệu nhưng chưa commit
    * Transaction 2 có thể đọc dữ liệu thay đổi của transaction 1
  * Read Commited
    * Transaction 1 thay đổi dữ liệu nhưng chưa commit
    * Transaction 2 không thể đọc dữ liệu thay đổi của transaction 1
  * Repeatable Read
    * Transaction 1 thay đổi dữ liệu và commit
    * Transaction 2 không thể đọc dữ liệu thay đổi của transaction 1 (đọc lại dữ liệu cũ)
  * Serializable
    * Transaction 1 thay đổi dữ liệu nhưng chưa commit
    * Transaction 2 phải đợi khi transaction 1 thực hiện xong mới chạy được câu query
Mysql :

* Các loại index     
  * Normal Index:
  ```sql
    create index idx_username on users(username)
  ```
  * Unique Index:
  ```sql
    create unique index idx_username on users(username)
  ```
  * Primary Index:
  ```sql
    primary key(id)
  ```
  * Full text index:
  ```sql
    create fulltext index idx_fulltext_username on  users(username)
  ```
  * Composite Index  
    * Ngoài cùng bên trái
      * Các cột có ít dữ liệu trùng nhau thì đưa về phía trái khi tạo index composite index
      ```sql
        create index idx_username_gender on users(username,gender)
      ```
      * gender có nhiều dữ liệu trùng hên nên để phía cuối cùng. Không nên dánh normal index cho những dư liệu trùng nhau nhiều vd column : gender

* sử dụng explain
  type : all 