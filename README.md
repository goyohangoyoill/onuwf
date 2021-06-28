# ONUWF
One Night Ultimate Werewolf ver 2

# 실행시 설치해야하는 go package

- mongoDB
```
go get go.mongodb.org/mongo-driver/mongo
```

- discordgo
```
go get github.com/bwmarrin/discordgo
```

- discordgo embed
```
go get github.com/clinet/discordgo-embed
```

# 개발 진행 상황

branch `file-seperate` : 강제 종료 구현을 용이하게 하기 위해 게임별로 다른 프로세스를 실행하고, 해당 게임이 종료되었을 때 (강제종료, 정상종료) 해당 프로세스를 끝마친다. 이를 위한 외부 프로세스와 내부 프로세스를 각각 `ONUWF`, `GameHandler` 가 맡는다.

# 빌드 방법

```
git clone https://github.com/splkm97/ONUWF.git
cd ONUWF/GameHandler
/* add env.json and Log dir */
go build
cd ..
go build
```

# 설계 진행 상황

https://draw.io 접속 후 struct_diagram.drawio 선택하여 편집/확인 가능

https://gitmind.com/app/doc/1362405404 -> 이전 설계
