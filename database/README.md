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