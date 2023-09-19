# HTTP proxy server, MITM

1. Сгенерировать корневой сертификат и его приватный ключ. Важно, чтобы файлы с именами ca.crt и ca.key находились в
   корневой директории проекта.

```./scripts/gen_ca.sh```

2. Скопировать сертификат.

```sudo cp ca.crt /usr/local/share/ca-certificates/ca.crt```

```sudo cp ca.crt /usr/share/ca-certificates/ca.crt```

3. Выполнить команду.

```sudo update-ca-certificates```

4. Запустить контейнеры.

```docker compose up -d --build storage api proxy```

5. Тестовые запросы.

```curl -v -x http://127.0.0.1:8080 https://example.org```

```curl -v http://127.0.0.1:8000/requests ```

6. Для работы в браузерах. Но я задавал свой центр сертификации через интерфейс настроек хрома.

```mkcert -install```

7. Просмотр логов контейнеров.

```docker logs -f mitm_proxy```

```docker logs -f mitm_api```
