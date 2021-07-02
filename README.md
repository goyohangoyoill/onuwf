# ONUWF
One Night Ultimate Werewolf ver 2

# 실행시 설치해야하는 go package

- mongoDB
```
go get go.mongodb.org/mongo-driver/mongo
// 아직 사용하지 않았다.
```

- discordgo
```
go get github.com/bwmarrin/discordgo
```

- discordgo embed
```
go get github.com/clinet/discordgo-embed
```

# 디렉터리 구조

`asset` : 게임을 위한 에셋
`util` : 데이터베이스 연동, 파일 입출력 등의 유틸리티 코드
`game` : 게임 내부 데이터를 위한 스트럭처 모음

# 추가해야하는 파일들

`config/env.json` -> 슬랙 고요일 채널 참고

# 설계 진행 상황

https://draw.io 접속 후 struct_diagram.drawio 선택하여 편집/확인 가능

https://gitmind.com/app/doc/1362405404 -> 이전 설계
