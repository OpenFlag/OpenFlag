package config

//nolint:lll
const defaultConfig = `
logger:
  access:
    enabled: true
    path: "./logs/access.log"
    format: "${remote_ip} - - [${time_rfc3339}] \"${method} ${uri} HTTP/1.1\" ${status} \
    ${bytes_out} ${bytes_in} ${latency} \"${referer}\" \"${user_agent}\"\n"
    max-size: 1024
    max-backups: 7
    max-age: 7
  app:
    level: debug
    path: "./logs/app.log"
    max-size: 1024
    max-backups: 7
    max-age: 7
    stdout: true

server:
  address: :7677
  read-timeout: 20s
  write-timeout: 20s
  graceful-timeout: 5s

postgres:
  host: 127.0.0.1
  port: 5432
  user: openflag
  pass: secret
  dbname: openflag
  connect-timeout: 30s
  connection-lifetime: 30m
  max-open-connections: 10
  max-idle-connections: 5

monitoring:
  prometheus:
    enabled: true
    address: ":9001"
`
