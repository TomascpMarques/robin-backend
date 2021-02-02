FROM redis

# define um volume para perssistência
VOLUME [ "/test" ]

# RUN apt update
# RUN apt-get install -y
# RUN apt-get upgrade -y

# RUN apt-get install build-essential -y

# WORKDIR /app
# COPY /redis-6.0.10 /app
# RUN make
# RUN make install

# muda para a dir que tem perssistência
WORKDIR /test
# Copia a configuração redis prefeita para o servidor
COPY redis.conf /test/redis.conf
# RUN rm -rf ../app

# expõe o container na porta 6379 (onde o serviço redis estará localizado)
EXPOSE 6379

# inicia o redis-server com as configurações pré-carregadas, e "liga-se" á porta 6379
CMD [ "redis-server", "/test/redis.conf", "--port", "6379" ]