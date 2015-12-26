heroku buildpacks:set https://github.com/ph3nx/heroku-binary-buildpack.git
heroku config:set PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/app/bin
mkdir bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/server.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/populate cmd/populate/populate.go

git checkout -b deploy
git add --all
git commit -m "Build"
git push heroku deploy:master -f
git checkout master
git branch -D deploy
