## olcode
An application based on webrtc aims at giving an excellent experience for both interviewer and programmer interviewee.
See our [demo](https://togetherers.com)

## Architecture
1. first step:
```
signal server   <-----> peer1
|
\-------------> peer2
```
peer1 & peer2 connect to signal server then exchange some basic information for p2p connecting.

2. second step
```
signal server   <-----> peer1
|                         |
\-------------> peer2 ---/
```
peer1 & peer2 now connected and they can talk with code/voice real-time.

## Contribution
Any Contributions including issues are appreciated. And some important issues have been created by collabarators already so you can choose some one for your start.

If you want to start developing now, just follows the steps to prepare a development environment first:

### develop dependencies
1. [golang](https://golang.org/)
2. [node](https://nodejs.org/)
3. [yarn](https://yarnpkg.com/)

### run signal server
```
cd olcode
cd devops
pip install -r requirements.txt
python dev.py -t run -s gw
```

### run webclient dev server
```
cd olcode/web
yarn install && yarn start
```
