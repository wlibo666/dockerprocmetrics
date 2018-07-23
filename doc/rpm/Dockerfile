From wangchunyan/centos:7.5.1804
MAINTAINER wangchunyan <wlibo666@126.com>

LABEL version="dpm1.0.0"
LABEL describe="docker process metrics"

RUN mkdir -p /etc/dpm/
COPY metrics.json /etc/dpm/metrics.json
COPY dockerprocmetrics /usr/bin/dockerprocmetrics
RUN chmod +x /usr/bin/dockerprocmetrics
ENTRYPOINT /usr/bin/dockerprocmetrics