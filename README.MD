# Разбираемся с работой в RabbitMQ в golang

Представим себе приложение, которое состоит из клиента и worker'ов.

Клиент отправляет сообщения в различные exchange'ы.

Worker принимает сообщения и выводит их в stdout.
Предполагается, что один worker обрабатывает одну очередь.  

Оно вроде работает. В дальнейшем будем измываться над этим каркасом.