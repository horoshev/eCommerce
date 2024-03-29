# Storage

Данный сервис отвечает за резервирование товаров для заказов.

## General

При старте приложения создаются товары. Они хранятся в таблице `products`.
Зарезервированные заказы хранятся в таблице `reservations`.

## Flow

Сервис имеет возможность зарезервировать товары для заказа.
Для этого сервис читает сообщения из топика `storage-reserve-order` и обрабатывает их.
Проверяется наличие каждого товара, после чего он резервируется и убирается и общего числа доступных товаров.
Также идет подсчет общей стоимости заказа.

При удачном выполнении создается бронь и ответ об удачном выполнении отравляется в топик `storage-reserve-order-response`.
При проверке наличия и резервировании товаров создается транзакция для сохранения консистентности данных.
Если какого-либо товара не хватает для заказа, весь заказ считается отмененным.

Также сервис позволяет отменить резервирование, и вернуть товары на склад.
Для этого сервис читает сообщения из топика `storage-cancel-order` и обрабатывает их.
При удачном выполнении отменяется бронь, товары из брони возвращаются в общий доступ и 
ответ об удачном выполнении отравляется в топик `storage-cancel-order-response`.
