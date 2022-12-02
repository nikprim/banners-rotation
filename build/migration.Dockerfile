FROM alpine:3.9

ENV WAITFORIT_VERSION="v2.4.1"
ENV WAIT_FOR_IT_PATH "/usr/local/bin/waitforit"
RUN wget -q -O $WAIT_FOR_IT_PATH https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 \
    && chmod +x $WAIT_FOR_IT_PATH

ENV GOOSE_PATH "/usr/local/bin/goose"
RUN wget -q -O $GOOSE_PATH https://github.com/pressly/goose/releases/latest/download/goose_linux_x86_64 \
    && chmod +x $GOOSE_PATH

CMD goose