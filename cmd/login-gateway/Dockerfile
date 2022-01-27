FROM jrei/systemd-centos:7

ARG token_access_secret=
ENV LOGIN_TOKEN_ACCESS_SECRET=$token_access_secret

COPY LoginGateway.viper.yaml /etc/LoginGateway/LoginGateway.viper.yaml
COPY login-gateway /usr/local/bin/login-gateway

CMD ["login-gateway", "run"]
