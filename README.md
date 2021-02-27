# interview_goldhub

client:
```
https://chrome.google.com/webstore/detail/websocket-test-client/fgponpodhbmadfljofbimhhlengambbn
```

To start the program
```
git clone https://github.com/miketsui3a/interview_goldhub.git
cd interview_goldhub
go run main.go
```


After program start up, use client to connect to:
```
ws://localhost:8089/ws
```

registration:
```
{"message":"29c9c30e0604515ced98b3d14fd88751a8f8e4b9bc69d483a67a257c14ab79fb","playerName":"your name","timestamp":1614399947}
```

guess:
```
{"message":"f1abe1b083d12d181ae136cfc75b8d18a8ecb43ac4e9d1a36d6a9c75b6016b61","guess":82,"timestamp":1614399947,"gameId":0}
```

