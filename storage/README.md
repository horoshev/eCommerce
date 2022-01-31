# Storage

Заказ формируется. 

При резервировании заказа создается транзакция для сохранения консистентности данных.

Проверяется наличие каждого товара, после чего он резервируется и убирается и общего числа доступных товаров.

Если какого-либо товара не хватает для заказа, весь заказ считается отмененным.

Топик для создания резервирования - `storage_order_create`
Топик для отмены резервирования - `storage_order_cancel`


## TODO

Получает заказ и резервирует товары из него.

При проверке наличия и резервировании товаров используется транзакция.

Внутри транзакции проверяется, что товары в наличии, если так, то мы резервирует товары и вычитает из общего количества.

Заказы хранятся в таблице reservations