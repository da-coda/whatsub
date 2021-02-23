FROM postgres:13.2-alpine

EXPOSE 5432
ENV MUSL_LOCALE_DEPS cmake make musl-dev gcc gettext-dev libintl
ENV MUSL_LOCPATH /usr/share/i18n/locales/musl
ENV POSTGRES_DB whatsub
ENV POSTGRES_INITDB_ARGS "--encoding=UTF8"
ENV POSTGRES_USER whatsub
ENV POSTGRES_PASSWORD whatsub

RUN apk add --no-cache \
    $MUSL_LOCALE_DEPS \
    && wget https://gitlab.com/rilian-la-te/musl-locales/-/archive/master/musl-locales-master.zip \
        && unzip musl-locales-master.zip \
          && cd musl-locales-master \
          && cmake -DLOCALE_PROFILE=OFF -D CMAKE_INSTALL_PREFIX:PATH=/usr . && make && make install \
          && cd .. && rm -r musl-locales-master


USER postgres
COPY Docker/postgresql/subreddits_2021-02-19.csv /subreddits_2021-02-19.csv
COPY Docker/postgresql/init.sql /docker-entrypoint-initdb.d/-init.sql

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
