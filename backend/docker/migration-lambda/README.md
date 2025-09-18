build push方法

backendのディレクトリで作業した。コマンドはうろ覚えだが、buildしてタグ付けしてpushする流れでやった。

docker build -t migration-lambda:latest -f docker/migration-lambda/Dockerfile .

docker tag migration-lambda:latest 420302062688.dkr.ecr.ap-northeast-1.amazonaws.com/tasmap-migration-lambda:latest

awsのecrにログイン

docker push 420302062688.dkr.ecr.ap-northeast-1.amazonaws.com/tasmap-migration-lambda:latest

