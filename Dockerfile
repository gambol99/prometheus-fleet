#
#   Author: Rohith
#   Date: 2015-08-22 13:38:12 +0100 (Sat, 22 Aug 2015)
#
#  vim:ts=2:sw=2:et
#
FROM progrium/busybox
MAINTAINER Rohith <gambol99@gmail.com>

ADD bin/prometheus-fleet /prometheus-fleet
RUN chmod +x prometheus-fleet

ENTRYPOINT [ "/prometheus-fleet" ]
