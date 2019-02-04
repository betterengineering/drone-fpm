FROM debian:stretch

RUN apt-get update && apt-get install -y \
    ruby \
    ruby-dev \
    rubygems \
    build-essential

RUN gem install --no-ri --no-rdoc fpm

COPY out/entrypoint.sh /bin/entrypoint.sh

ENTRYPOINT ["/bin/entrypoint.sh"]