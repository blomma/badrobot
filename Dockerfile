FROM arm32v6/alpine

COPY bin/linux-arm-7-badrobot /linux-arm-7-badrobot
COPY badfriends.html /badfriends.html

EXPOSE 8001

CMD ["/linux-arm-7-badrobot"]
