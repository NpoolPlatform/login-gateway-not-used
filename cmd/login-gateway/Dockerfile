FROM jrei/systemd-centos:7

COPY LoginGateway.viper.yaml /etc/LoginGateway/LoginGateway.viper.yaml
COPY login-gateway /usr/local/bin/login-gateway

CMD ["login-gateway", "run"]
