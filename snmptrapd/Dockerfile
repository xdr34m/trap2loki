FROM alpine:latest

RUN apk --no-cache add net-snmp
#RUN echo 'disableAuthorization yes' > /etc/snmp/snmptrapd.conf
# accept all traps without v3 - only v2

#EXPOSE 162

CMD ["snmptrapd","-L","o","-f"]
# setup correct snmptrapd options to log to /spool? for snmptt to read from it?