* Cache problem:
  * Cache avalanche
    * Set thời gian expired của dữ liệu khác nhau
    * Sử dụng singleflight nghĩa nhiều request sử dụng chung key cache nếu expired thì chỉ có 1 request gọi xuống database
    * Cache cùng 1 dữ liệu với 2 key khác nhau nếu key chính ko tìm thấy trả về dữ liệu key phụ 
  * Cache Breakdown
    * Set up redis cluster ( master - slave +  sential)
  * Cache Penetration 
    * Bloom Filter
    * set giá trị không tồn tại trong redis
* Data type:
  * String
    * store string
  * LIST
    * Linked list string
    * Sử dụng: 
    * Sử dụng làm stack and queue
    * build queue
  * SET
    * Dữ liệu duy nhất và không sắp xếp
    * Sử dụng : 
    * Kiểm tra dữ liệu unique
    * Kiểm tra value nào đó có nằm trong set không
    * Kiểm tra 2 set có chung một value nào không  
  * HASHES
    * Dữ liệu lưu vào dạng key value
    * Sử dụng:
      * Lưu dữ liệu của một object
  * Sorted SETS
    * Chứa dữ liệu duy nhất và được sắp xếp bới điểm số
    * Sử dụng 
    * Sử dụng để sắp xếp dữ liệu với score