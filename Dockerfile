FROM alpine
ADD sso-srv /sso-srv
ENTRYPOINT [ "/sso-srv" ]
