FROM redis

VOLUME [ "/test" ]

# RUN apt update
# RUN apt-get install -y
# RUN apt-get upgrade -y

# RUN apt-get install build-essential -y

# WORKDIR /app
# COPY /redis-6.0.10 /app
# RUN make
# RUN make install

WORKDIR /test
COPY redis.conf /test/redis.conf
# RUN rm -rf ../app

EXPOSE 6379

CMD [ "redis-server", "/test/redis.conf", "--port", "6379" ]