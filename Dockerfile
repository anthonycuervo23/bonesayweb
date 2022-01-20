FROM scratch
MAINTAINER Rid <rid@cylo.io>
ADD dist/bonettpsay_linux_amd64/bonettpsay bonettpsay
ADD static static
ADD templates templates
CMD ["/bonettpsay"]
EXPOSE 8000