# ТЗ на демон, осуществляющий "Системный мониторинг"

## Общее описание
Демон - программа, собирающая информацию о системе, на которой запущена,
и отправляющая её своим клиентам по GRPC.

## Архитектура
- GRPC сервер;
- допускается использование временных (`/tmp`) файлов;
- статистика хранится в памяти, долговременное хранение не предусмотрено.

## Требования
Необходимо каждые **N** секунд выдавать информацию, усредненную за последние **M** секунд.

Например, N = 5с, а M = 15с, тогда демон "молчит" первые 15 секунд,
затем выдает снапшот за 0-15с; через 5с (в 20с) выдает снапшот за 5-20с;
через 5с (в 25с) выдает снапшот за 10-25с и т.д.

**N** и **M** указывает клиент в запросе на получение статистики.

Что необходимо собирать:
- Средняя загрузка системы (load average).

- Средняя загрузка CPU (%user_mode, %system_mode, %idle).

- Загрузка дисков:
    - tps (transfers per second);
    - KB/s (kilobytes (read+write) per second);
    - CPU (%user_mode, %system_mode, %idle).

- Информация о дисках по каждой файловой системе:
    - использовано мегабайт, % от доступного количества;
    - использовано inode, % от доступного количества.

- Top talkers по сети:
    - по протоколам: protocol (TCP, UDP, ICMP, etc), bytes, % от sum(bytes) за последние **M**), сортируем по убыванию процента;
    - по трафику: source ip:port, destination ip:port, protocol, bytes per second (bps), сортируем по убыванию bps.

- Статистика по сетевым соединениям:
    - слушающие TCP & UDP сокеты: command, pid, user, protocol, port;
    - количество TCP соединений, находящихся в разных состояниях (ESTAB, FIN_WAIT, SYN_RCV и пр.).

#### Разрешено использовать только стандартную библиотеку языка Go!

Команды, которые могут пригодиться:
```
$ top -b -n1
$ df -k
$ df -i
$ iostat -d -k
$ cat /proc/net/dev
$ sudo netstat -lntup
$ ss -ta
$ tcpdump -ntq -i any -P inout -l
$ tcpdump -nt -i any -P inout -ttt -l
```

Статистика представляет собой объекты, описанные в формате Protobuf.

Информацию необходимо выдавать всем подключенным по GRPC клиентам
с использованием [однонаправленного потока](https://grpc.io/docs/tutorials/basic/go/#server-side-streaming-rpc).

Выдавать "снапшот" системы можно как отдельными сообщениями, так и одним жирным объектом.

Сбор информации, её парсинг и пр. должен осуществлятся как можно более конкуррентно.

## Поддерживаемая ОС
Минимум - Linux (Ubuntu 18.04).

Максимум - несколько сборок под набор из популярных ОС/процессоров:
- darwin, linux, windows
- 386, amd64

[Список возможных вариантов](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).

Но тогда придется постараться с реализацией использования различных команд для сбора данных.

Пригодятся [build тэги](https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures).

## Конфигурация
- Через аргументы командной строки можно указать, на каком порту стартует сервер.
- Через файл можно указать, какие из подсистем сбора включены/выключены.

## Тестирование
#### Юнит-тесты
- по возможности мок интерфейсов и проверка вызовов конкретных методов;
- тесты вспомогательных функций и пр.

#### Интеграционные тесты
- потестировать факт потока статистики, можно без конкретных цифр;
- можно посоздавать файлы, пооткрывать сокеты и посмотреть на изменение снапшота.

#### Клиент
Необходимо реализовать простой клиент, который в реальном времени получает
и выводит в STDOUT статистику по одному из пунктов (например, сетевую информацию)
в читаемом формате (например, в виде таблицы).

## Разбалловка
Максимум - **20 баллов**
(при условии выполнения [обязательных требований](./README.md)):

* Реализован сбор:
    - load average - 1 балл;
    - загрузка CPU - 1 балл;
    - загрузка дисков - 1 балл;
    - top talkers по сети - 1 балла;
    - статистика по сети - 1 балл.
* Через конфигурацию можно отключать отдельную статистику - 2 балла.
* Написаны юнит-тесты - 1 балл.
* Написаны интеграционные тесты - 2 балла.
* Реализован простой клиент к демону - 2 балла.
* Сбор хотя бы одного типа статистики работает на разных ОС - 5 баллов.
* Понятность и чистота кода - до 3 баллов.

#### Зачёт от 10 баллов